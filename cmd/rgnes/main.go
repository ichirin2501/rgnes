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

	// i := 0
	// vblankOnClock := 0
	// for {
	// 	beforeN := ppu.FetchNMIDelay()
	// 	beforeV := ppu.FetchVBlankStarted()
	// 	ppu.Step()
	// 	i++
	// 	afterV := ppu.FetchVBlankStarted()
	// 	afterN := ppu.FetchNMIDelay()
	// 	if !beforeV && afterV {
	// 		fmt.Printf("vblank:1\tclock:%d\n", i)
	// 		vblankOnClock = i
	// 	} else if beforeV && !afterV {
	// 		fmt.Printf("vblank:0\tclock:%d\t1->0:%d\t=cpuClock:%d(=%d)\n", i, i-vblankOnClock, (i-vblankOnClock)/3, (i-vblankOnClock)%3)
	// 		fmt.Println("")
	// 		return nil
	// 	} else if beforeN > 0 && afterN == 0 {
	// 		fmt.Printf("trigger nmi clock: %d\n", i)
	// 	}
	// }

	// return nil

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
			cycle := cpu.Step()

			// fmt.Println(trace.NESTestString())
			//fmt.Printf("cpu.clock:%d\tppu.clock:%d\tdiff:%d\n", cpu.FetchCycles(), ppu.Clock, cpu.FetchCycles()*3-ppu.Clock)
			if cpu.FetchCycles()*3 != ppu.Clock {
				panic("eeeeeeeeeeeeeeeeee")
			}

			trace.AddCPUCycle(cycle)

			//beforeI := irp.I
			// for i := 0; i < cycle*3; i++ {
			// 	ppu.Step()
			// }

			if beforeppuy > trace.PPUY {
				joypad.SetButtonPressedStatus(0)
				<-ticker.C
				//fmt.Println("ticker on: ", time.Now())
			}
			beforeppuy = trace.PPUY
			//afterI := irp.I
			//fmt.Printf("interrupt type: before:%v after:%v\n", beforeI, afterI)
			select {
			case k := <-keyEvents:
				switch k {
				case fyne.KeyEscape:
					win.Close()
				case fyne.KeySpace:
					joypad.SetButtonPressedStatus(nes.ButtonSelect)
				case fyne.KeyReturn:
					joypad.SetButtonPressedStatus(nes.ButtonStart)
				case fyne.KeyUp:
					joypad.SetButtonPressedStatus(nes.ButtonUP)
				case fyne.KeyDown:
					joypad.SetButtonPressedStatus(nes.ButtonDown)
				case fyne.KeyLeft:
					joypad.SetButtonPressedStatus(nes.ButtonLeft)
				case fyne.KeyRight:
					joypad.SetButtonPressedStatus(nes.ButtonRight)
				case fyne.KeyA:
					joypad.SetButtonPressedStatus(nes.ButtonA)
				case fyne.KeyS:
					joypad.SetButtonPressedStatus(nes.ButtonB)
				}
			default:
				//joypad.SetButtonPressedStatus(0)
			}
			canvasImg.Refresh()
			//time.Sleep(10 * time.Microsecond)
			//time.Sleep(1 * time.Millisecond)
		}
	}()

	win.ShowAndRun()

	return nil
}
