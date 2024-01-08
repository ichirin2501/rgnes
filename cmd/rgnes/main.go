package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/hajimehoshi/oto"
	"github.com/ichirin2501/rgnes/nes"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type renderer struct {
	fyne.Window
	currImg *canvas.Image
	nextImg *canvas.Image
	ticker  *time.Ticker
}

// todo: fix data race
func newRenderer(win fyne.Window, curr, next *canvas.Image) *renderer {
	return &renderer{
		Window:  win,
		currImg: curr,
		nextImg: next,
		ticker:  time.NewTicker(16 * time.Millisecond),
	}
}

func (r *renderer) Render(x, y int, c color.Color) {
	r.nextImg.Image.(*image.RGBA).Set(x, y, c)
}
func (r *renderer) Refresh() {
	<-r.ticker.C
	r.currImg, r.nextImg = r.nextImg, r.currImg // swap
	r.SetContent(container.NewMax(r.currImg))
	r.currImg.Refresh()
}

type player struct {
	p   *oto.Player
	buf []byte
}

func newPlayer() (*player, error) {
	c, err := oto.NewContext(44100, 1, 1, 1)
	if err != nil {
		return nil, err
	}
	p := c.NewPlayer()
	return &player{
		p:   p,
		buf: make([]byte, 1),
	}, nil
}

func (p *player) Sample(v float32) {
	p.buf[0] = byte(v * 0xFF)
	if _, err := p.p.Write(p.buf); err != nil {
		fmt.Println("why: ", err)
	}
}

func realMain() error {
	var (
		rom string
	)
	flag.StringVar(&rom, "rom", "", "rome filepath")
	flag.Parse()

	myapp := app.New()
	win := myapp.NewWindow("rgnes")
	img1 := image.NewRGBA(image.Rect(0, 0, 256, 240))
	img2 := image.NewRGBA(image.Rect(0, 0, 256, 240))

	canvasImg1 := canvas.NewImageFromImage(img1)
	canvasImg2 := canvas.NewImageFromImage(img2)
	canvasImg1.SetMinSize(fyne.NewSize(256, 240))
	canvasImg2.SetMinSize(fyne.NewSize(256, 240))
	canvasImg1.ScaleMode = canvas.ImageScalePixels
	canvasImg2.ScaleMode = canvas.ImageScalePixels

	win.SetContent(container.NewMax(
		canvasImg1,
	))
	renderer := newRenderer(win, canvasImg1, canvasImg2)

	f, err := os.Open(rom)
	if err != nil {
		return err
	}
	defer f.Close()

	mapper, err := nes.NewMapper(f)
	if err != nil {
		return err
	}

	player, err := newPlayer()
	if err != nil {
		return err
	}

	trace := &nes.Trace{}
	irp := &nes.Interrupter{}

	m := mapper.MirroingType()
	ppu := nes.NewPPU(renderer, mapper, m, irp)
	joypad := nes.NewJoypad()
	apu := nes.NewAPU(irp, player)
	cpuBus := nes.NewBus(ppu, apu, mapper, joypad)

	cpu := nes.NewCPU(cpuBus, irp, nes.WithTracer(trace))
	apu.PowerUp()
	cpu.PowerUp()

	if deskCanvas, ok := win.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(k *fyne.KeyEvent) {
			updateKey(win, cpu, joypad, k.Name, true)
		})
		deskCanvas.SetOnKeyUp(func(k *fyne.KeyEvent) {
			updateKey(win, cpu, joypad, k.Name, false)
		})
	}

	go func() {
		for {
			trace.Reset()

			// ここでppuの状態を記録しておく
			trace.SetPPUX(uint16(ppu.Cycle))
			trace.SetPPUY(uint16(ppu.Scanline))
			// v := ppu.FetchV()
			// mp0 := mapper.Read(0)
			// ppuBuf := ppu.FetchBuffer()
			//beforeScanline := ppu.Scanline
			cpu.Step()

			// fmt.Printf("%s apuSteps:%d\tapuFrameMode:%d\tapuFrameSeqStep:%d\tapuPulse1LC:%d\tframeIRQFlag:%v\tnewfval:%v\twriteDelayFC:%v\n", trace.NESTestString(),
			// 	apu.FetchFrameStep(),
			// 	apu.FetchFrameMode(),
			// 	apu.FetchFrameSeqStep(),
			// 	apu.FetchPulse1LC(),
			// 	apu.FetchFrameIRQFlag(),
			// 	apu.FetchNewFrameCounterVal(),
			// 	apu.FetchWriteDelayFC(),
			// )
			//fmt.Printf("%s ppu.v:0x%04X ppu.buf:0x%02X mapper[0]:0x%02X\n", trace.NESTestString(), v, ppuBuf, mp0)
			//fmt.Printf("0x6000 = 0x%02X\n", cpuBus.ReadForTest(0x6000))

			// if cpu.FetchCycles()*3 != ppu.Clock {
			// 	panic("eeeeeeeeeeeeeeeeee")
			// }
		}
	}()

	win.ShowAndRun()

	return nil
}

func updateKey(win fyne.Window, cpu *nes.CPU, j *nes.Joypad, k fyne.KeyName, pressed bool) {
	switch k {
	case fyne.KeyEscape:
		win.Close()
	case fyne.KeyR:
		cpu.Reset()
	case fyne.KeySpace:
		j.SetButtonStatus(nes.ButtonSelect, pressed)
	case fyne.KeyReturn:
		j.SetButtonStatus(nes.ButtonStart, pressed)
	case fyne.KeyUp:
		j.SetButtonStatus(nes.ButtonUP, pressed)
	case fyne.KeyDown:
		j.SetButtonStatus(nes.ButtonDown, pressed)
	case fyne.KeyLeft:
		j.SetButtonStatus(nes.ButtonLeft, pressed)
	case fyne.KeyRight:
		j.SetButtonStatus(nes.ButtonRight, pressed)
	case fyne.KeyZ:
		j.SetButtonStatus(nes.ButtonA, pressed)
	case fyne.KeyX:
		j.SetButtonStatus(nes.ButtonB, pressed)
	}
}
