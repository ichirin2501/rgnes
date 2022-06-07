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
	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/cpu"
	"github.com/ichirin2501/rgnes/nes/joypad"
	"github.com/ichirin2501/rgnes/nes/ppu"
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

	mapper, err := cassette.NewMapper(f)
	if err != nil {
		return err
	}

	trace := &cpu.Trace{}
	irp := &cpu.Interrupter{}

	m := mapper.MirroingType()
	ppu := ppu.New(renderer, mapper, &m, irp)
	joypad := joypad.New()
	apu := apu.New(irp)
	cpuBus := cpu.NewBus(ppu, apu, mapper, joypad)

	cpu := cpu.New(cpuBus, irp, cpu.WithTracer(trace))
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

			// fmt.Printf("%s apuSteps:%d\tapuFrameMode:%d\tapuFrameSeqStep:%d\tapuPulse1LC:%d xxx\n", trace.NESTestString(),
			// 	apu.FetchFrameStep(),
			// 	apu.FetchFrameMode(),
			// 	apu.FetchFrameSeqStep(),
			// 	apu.FetchPulse1LC(),
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

func updateKey(win fyne.Window, cpu *cpu.CPU, j *joypad.Joypad, k fyne.KeyName, pressed bool) {
	switch k {
	case fyne.KeyEscape:
		win.Close()
	case fyne.KeyR:
		cpu.Reset()
	case fyne.KeySpace:
		j.SetButtonStatus(joypad.ButtonSelect, pressed)
	case fyne.KeyReturn:
		j.SetButtonStatus(joypad.ButtonStart, pressed)
	case fyne.KeyUp:
		j.SetButtonStatus(joypad.ButtonUP, pressed)
	case fyne.KeyDown:
		j.SetButtonStatus(joypad.ButtonDown, pressed)
	case fyne.KeyLeft:
		j.SetButtonStatus(joypad.ButtonLeft, pressed)
	case fyne.KeyRight:
		j.SetButtonStatus(joypad.ButtonRight, pressed)
	case fyne.KeyZ:
		j.SetButtonStatus(joypad.ButtonA, pressed)
	case fyne.KeyX:
		j.SetButtonStatus(joypad.ButtonB, pressed)
	}
}
