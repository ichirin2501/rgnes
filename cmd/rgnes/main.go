package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"runtime"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/hajimehoshi/oto"
	"github.com/ichirin2501/rgnes/nes"
)

func init() {
	runtime.LockOSThread()
}

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

	n := nes.New(mapper, renderer, player)
	n.PowerUp()

	if deskCanvas, ok := win.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(k *fyne.KeyEvent) {
			updateKey(win, n, k.Name, true)
		})
		deskCanvas.SetOnKeyUp(func(k *fyne.KeyEvent) {
			updateKey(win, n, k.Name, false)
		})
	}

	go func() {
		for {
			n.Step()
		}
	}()

	win.ShowAndRun()

	return nil
}

func updateKey(win fyne.Window, n *nes.NES, k fyne.KeyName, pressed bool) {
	switch k {
	case fyne.KeyEscape:
		win.Close()
	case fyne.KeyR:
		n.Reset()
	case fyne.KeySpace:
		n.SetButtonStatus(nes.ButtonSelect, pressed)
	case fyne.KeyReturn:
		n.SetButtonStatus(nes.ButtonStart, pressed)
	case fyne.KeyUp:
		n.SetButtonStatus(nes.ButtonUP, pressed)
	case fyne.KeyDown:
		n.SetButtonStatus(nes.ButtonDown, pressed)
	case fyne.KeyLeft:
		n.SetButtonStatus(nes.ButtonLeft, pressed)
	case fyne.KeyRight:
		n.SetButtonStatus(nes.ButtonRight, pressed)
	case fyne.KeyZ:
		n.SetButtonStatus(nes.ButtonA, pressed)
	case fyne.KeyX:
		n.SetButtonStatus(nes.ButtonB, pressed)
	}
}
