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
	"github.com/ichirin2501/rgnes/nes"
	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/cpu"
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
}

func newRenderer(win fyne.Window, curr, next *canvas.Image) *renderer {
	return &renderer{
		Window:  win,
		currImg: curr,
		nextImg: next,
	}
}

// todo: fix data race
func (r *renderer) Render(x, y int, c color.Color) {
	r.nextImg.Image.(*image.RGBA).Set(x, y, c)
}
func (r *renderer) Refresh() {
	r.currImg, r.nextImg = r.nextImg, r.currImg // swap
	r.SetContent(container.NewVBox(r.currImg))
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
	// TODO: windowを調節したときに比を維持してほしい
	canvasImg1.FillMode = canvas.ImageFillOriginal
	canvasImg2.FillMode = canvas.ImageFillOriginal
	win.SetContent(container.NewVBox(
		canvasImg1,
	))
	keyEvents := make(chan fyne.KeyName, 5)
	win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		keyEvents <- k.Name
	})

	renderer := newRenderer(win, canvasImg1, canvasImg2)

	f, err := os.Open(rom)
	if err != nil {
		return err
	}
	defer f.Close()

	c, err := cassette.NewCassette(f)
	if err != nil {
		return err
	}
	mapper := cassette.NewMapper(c)

	trace := &cpu.Trace{}
	irp := &cpu.Interrupter{}

	ppu := ppu.NewPPU(renderer, mapper, c.Mirror, irp, trace)
	joypad := nes.NewJoypad()
	apu := apu.NewAPU()
	cpuBus := cpu.NewBus(ppu, apu, mapper, joypad)

	cpu := cpu.NewCPU(cpuBus, irp, trace)
	cpu.Reset()
	trace.AddCPUCycle(7)
	for i := 0; i < 15; i++ {
		ppu.Step()
	}

	go func() {
		ticker := time.NewTicker(16 * time.Millisecond)
		defer ticker.Stop()
		beforeppuy := uint16(0)

		for {
			trace.Reset()

			// ここでppuの状態を記録しておく
			trace.SetPPUX(uint16(ppu.Cycle))
			trace.SetPPUY(uint16(ppu.FetchScanline()))
			beforeScanline := ppu.FetchScanline()
			cycle := cpu.Step()

			// fmt.Println(trace.NESTestString())

			if cpu.FetchCycles()*3 != ppu.Clock {
				panic("eeeeeeeeeeeeeeeeee")
			}

			trace.AddCPUCycle(cycle)
			if beforeScanline != 240 && ppu.FetchScanline() == 240 {
				updateKey(win, keyEvents, joypad)
			}

			if beforeppuy > trace.PPUY {
				<-ticker.C
			}
			beforeppuy = trace.PPUY
		}
	}()

	win.ShowAndRun()

	return nil
}

func updateKey(win fyne.Window, keyEvents <-chan fyne.KeyName, j *nes.Joypad) {
	keySt := byte(0)
	loop := true
	for loop {
		select {
		case k := <-keyEvents:
			switch k {
			case fyne.KeyEscape:
				win.Close()
			case fyne.KeySpace:
				keySt |= nes.ButtonSelect
			case fyne.KeyReturn:
				keySt |= nes.ButtonStart
			case fyne.KeyUp:
				keySt |= nes.ButtonUP
			case fyne.KeyDown:
				keySt |= nes.ButtonDown
			case fyne.KeyLeft:
				keySt |= nes.ButtonLeft
			case fyne.KeyRight:
				keySt |= nes.ButtonRight
			case fyne.KeyA:
				keySt |= nes.ButtonA
			case fyne.KeyS:
				keySt |= nes.ButtonB
			}
		default:
			loop = false
		}
	}
	j.SetButtonStatus(keySt)
}
