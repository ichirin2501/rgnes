package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/ichirin2501/rgnes/nes"
	"github.com/ichirin2501/rgnes/nes/apu"
	"github.com/ichirin2501/rgnes/nes/bus"
	"github.com/ichirin2501/rgnes/nes/cassette"
	"github.com/ichirin2501/rgnes/nes/cpu"
	"github.com/ichirin2501/rgnes/nes/memory"
	"github.com/ichirin2501/rgnes/nes/ppu"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type renderer struct {
	img *image.RGBA
}

func newRenderer(img *image.RGBA) *renderer {
	return &renderer{
		img: img,
	}
}

func (r *renderer) Render(x, y int, c color.Color) {
	r.img.Set(x, y, c)
}

func realMain() error {
	var (
		rom string
	)
	flag.StringVar(&rom, "rom", "", "rome filepath")
	flag.Parse()

	myapp := app.New()
	win := myapp.NewWindow("rgnes")
	img := image.NewRGBA(image.Rect(0, 0, 256, 240))

	canvasImg := canvas.NewImageFromImage(img)
	// TODO: windowを調節したときに比を維持してほしい
	canvasImg.FillMode = canvas.ImageFillOriginal
	win.SetContent(container.NewVBox(
		canvasImg,
	))
	keyEvents := make(chan fyne.KeyName, 5)
	win.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		fmt.Println("press ", k.Name)
		keyEvents <- k.Name
	})

	renderer := newRenderer(img)

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

	cycle := cpu.NewCPUCycle()
	ram := memory.NewMemory(0x8100)

	// debug
	// for i := 0; i < len(c.CHR); i++ {
	// 	if c.CHR[i] != 0 {
	// 		fmt.Printf("%04x: %02x\n", i, c.CHR[i])
	// 	}
	// }
	// return nil

	trace := &cpu.Trace{}
	irp := &cpu.Interrupter{}

	ppu := ppu.NewPPU(renderer, mapper, c.Mirror, irp, trace)

	joypad := nes.NewJoypad()
	apu := apu.NewAPU()
	cpuBus := bus.NewCPUBus(ram, ppu, apu, mapper, joypad)

	cpu := cpu.NewCPU(cpuBus, cycle, irp, trace)
	cpu.Reset()
	trace.AddCPUCycle(7)
	for i := 0; i < 7*3; i++ {
		ppu.Step()
	}

	go func() {
		for {
			trace.Reset()
			cycle := cpu.Step()

			//fmt.Println(trace.NESTestString())

			trace.AddCPUCycle(cycle)

			//beforeI := irp.I
			for i := 0; i < cycle*3; i++ {
				ppu.Step()
			}
			//afterI := irp.I
			//fmt.Printf("interrupt type: before:%v after:%v\n", beforeI, afterI)
			select {
			case k := <-keyEvents:
				switch k {
				case fyne.KeyEscape:
					win.Close()
				case fyne.KeySpace:
					joypad.SetButtonPressedStatus(nes.ButtonSelect, true)
				case fyne.KeyReturn:
					joypad.SetButtonPressedStatus(nes.ButtonStart, true)
				case fyne.KeyUp:
					joypad.SetButtonPressedStatus(nes.ButtonUP, true)
				case fyne.KeyDown:
					joypad.SetButtonPressedStatus(nes.ButtonDown, true)
				case fyne.KeyLeft:
					joypad.SetButtonPressedStatus(nes.ButtonLeft, true)
				case fyne.KeyRight:
					joypad.SetButtonPressedStatus(nes.ButtonRight, true)
				}
			default:
			}
			canvasImg.Refresh()
			//time.Sleep(10 * time.Microsecond)
			//time.Sleep(1 * time.Millisecond)
		}
	}()

	win.ShowAndRun()

	return nil
}
