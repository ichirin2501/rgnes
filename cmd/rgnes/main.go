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
	"github.com/gordonklaus/portaudio"
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
	img *canvas.Image
}

func newRenderer(win fyne.Window, img *canvas.Image) *renderer {
	return &renderer{
		Window: win,
		img:    img,
	}
}

func (r *renderer) Render(x, y int, c color.Color) {
	r.img.Image.(*image.RGBA).Set(x, y, c)
}
func (r *renderer) Refresh() {
	r.img.Refresh()
}

// ref: https://github.com/fogleman/nes/blob/3880f3400500b1ff2e89af4e12e90be46c73ae07/ui/audio.go#L5
type Player struct {
	stream         *portaudio.Stream
	sampleRate     float64
	outputChannels int
	channel        chan float32
}

func newPlayer() (*Player, error) {
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)

	p := Player{
		sampleRate:     parameters.SampleRate,
		outputChannels: parameters.Output.Channels,
		// If this channel size is too large (e.g. 44100), the BGM will be delayed. Make the size not too big
		channel: make(chan float32, 3000),
	}

	cbFunc := func(out []float32) {
		var output float32
		for i := range out {
			if i%p.outputChannels == 0 {
				select {
				case sample := <-p.channel:
					output = sample
				default:
					output = 0
				}
			}
			out[i] = output
		}
	}
	stream, err := portaudio.OpenStream(parameters, cbFunc)
	if err != nil {
		return nil, err
	}

	p.stream = stream
	return &p, nil
}

func (p *Player) Start() error {
	return p.stream.Start()
}

func (p *Player) Stop() error {
	return p.stream.Close()
}

func (p *Player) Sample(v float32) {
	p.channel <- v
}

func (p *Player) SampleRate() float64 {
	return p.sampleRate
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
	canvasImg.SetMinSize(fyne.NewSize(256*2, 240*2))
	canvasImg.ScaleMode = canvas.ImageScalePixels
	win.SetContent(container.NewStack(canvasImg))
	renderer := newRenderer(win, canvasImg)

	f, err := os.Open(rom)
	if err != nil {
		return err
	}
	defer f.Close()

	mapper, err := nes.NewMapper(f)
	if err != nil {
		return err
	}

	portaudio.Initialize()
	defer portaudio.Terminate()
	player, nil := newPlayer()
	if err != nil {
		return err
	}
	if err := player.Start(); err != nil {
		return err
	}
	defer player.Stop()

	n := nes.New(mapper, renderer, player)

	if deskCanvas, ok := win.Canvas().(desktop.Canvas); ok {
		deskCanvas.SetOnKeyDown(func(k *fyne.KeyEvent) {
			updateKey(n, k.Name, true)
		})
		deskCanvas.SetOnKeyUp(func(k *fyne.KeyEvent) {
			updateKey(n, k.Name, false)
		})
	}

	go func() {
		defer win.Close()
		// win.ShowAndRun() is required to run in the main thread, but a data race occurs during initialization within ShowAndRun().
		// So so wait a little while.
		time.Sleep(2 * time.Second)
		n.PowerUp()
		n.Run()
	}()

	win.ShowAndRun()

	return nil
}

func updateKey(n *nes.NES, k fyne.KeyName, pressed bool) {
	switch k {
	case fyne.KeyEscape:
		n.Close()
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
