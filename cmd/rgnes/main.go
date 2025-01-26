package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"runtime"
	"sync"

	"github.com/gordonklaus/portaudio"
	"github.com/ichirin2501/rgnes/nes"

	rl "github.com/gen2brain/raylib-go/raylib"
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
	evenImg  *image.RGBA
	oddImg   *image.RGBA
	oddFrame bool
	mu       *sync.Mutex
}

func newRenderer(evenImg, oddImg *image.RGBA) *renderer {
	return &renderer{
		evenImg: evenImg,
		oddImg:  oddImg,
		mu:      &sync.Mutex{},
	}
}

func (r *renderer) Render(x, y int, c color.Color) {
	if r.oddFrame {
		r.oddImg.Set(x, y, c)
	} else {
		r.evenImg.Set(x, y, c)
	}
}
func (r *renderer) Refresh() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.oddFrame {
		r.oddFrame = false
	} else {
		r.oddFrame = true
	}
}
func (r *renderer) CurerntImage() *rl.Image {
	r.mu.Lock()
	defer r.mu.Unlock()
	if !r.oddFrame {
		return rl.NewImageFromImage(r.oddImg)
	} else {
		return rl.NewImageFromImage(r.evenImg)
	}
}

// ref: https://github.com/fogleman/nes/blob/3880f3400500b1ff2e89af4e12e90be46c73ae07/ui/audio.go#L5
type Player struct {
	stream         *portaudio.Stream
	sampleRate     float64
	volume         float32
	outputChannels int
	channel        chan float32
}

func newPlayer(volume float32) (*Player, error) {
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		return nil, err
	}
	parameters := portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice)

	p := Player{
		sampleRate:     parameters.SampleRate,
		outputChannels: parameters.Output.Channels,
		volume:         volume,
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
	p.channel <- v * p.volume
}

func (p *Player) SampleRate() float64 {
	return p.sampleRate
}

func realMain() error {
	var (
		rom    string
		scale  int
		volume float64
	)
	flag.StringVar(&rom, "rom", "", "rome filepath")
	flag.IntVar(&scale, "scale", 2, "window scale size")
	flag.Float64Var(&volume, "volume", 0.5, "volume scale size")
	flag.Parse()

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
	player, nil := newPlayer(float32(volume))
	if err != nil {
		return err
	}
	if err := player.Start(); err != nil {
		return err
	}
	defer player.Stop()

	evenImg := image.NewRGBA(image.Rect(0, 0, nes.ScreenWidth, nes.ScreenHeight))
	oddImg := image.NewRGBA(image.Rect(0, 0, nes.ScreenWidth, nes.ScreenHeight))
	renderer := newRenderer(evenImg, oddImg)
	n := nes.New(mapper, renderer, player)

	rl.SetTraceLogLevel(rl.LogWarning)
	rl.InitWindow(nes.ScreenWidth*int32(scale), nes.ScreenHeight*int32(scale), "rgnes")
	defer rl.CloseWindow()

	currImg := renderer.CurerntImage()
	texture := rl.LoadTextureFromImage(currImg)
	defer rl.UnloadTexture(texture)
	rl.UnloadImage(currImg)

	rl.SetTargetFPS(60)

	go func() {
		n.PowerUp()
		n.Run()
	}()

	for !rl.WindowShouldClose() {
		rl.UnloadTexture(texture)
		currImg = renderer.CurerntImage()
		newTexture := rl.LoadTextureFromImage(currImg)
		rl.UnloadImage(currImg)
		texture = newTexture

		if rl.IsKeyDown(rl.KeyR) {
			n.Reset()
		}

		if rl.IsKeyDown(rl.KeyUp) {
			n.SetButtonStatus(nes.ButtonUP, true)
		}
		if rl.IsKeyDown(rl.KeyDown) {
			n.SetButtonStatus(nes.ButtonDown, true)
		}
		if rl.IsKeyDown(rl.KeyLeft) {
			n.SetButtonStatus(nes.ButtonLeft, true)
		}
		if rl.IsKeyDown(rl.KeyRight) {
			n.SetButtonStatus(nes.ButtonRight, true)
		}
		if rl.IsKeyDown(rl.KeySpace) {
			n.SetButtonStatus(nes.ButtonSelect, true)
		}
		if rl.IsKeyDown(rl.KeyEnter) {
			n.SetButtonStatus(nes.ButtonStart, true)
		}
		if rl.IsKeyDown(rl.KeyZ) {
			n.SetButtonStatus(nes.ButtonA, true)
		}
		if rl.IsKeyDown(rl.KeyX) {
			n.SetButtonStatus(nes.ButtonB, true)
		}
		// release
		if rl.IsKeyReleased(rl.KeyUp) {
			n.SetButtonStatus(nes.ButtonUP, false)
		}
		if rl.IsKeyReleased(rl.KeyDown) {
			n.SetButtonStatus(nes.ButtonDown, false)
		}
		if rl.IsKeyReleased(rl.KeyLeft) {
			n.SetButtonStatus(nes.ButtonLeft, false)
		}
		if rl.IsKeyReleased(rl.KeyRight) {
			n.SetButtonStatus(nes.ButtonRight, false)
		}
		if rl.IsKeyReleased(rl.KeySpace) {
			n.SetButtonStatus(nes.ButtonSelect, false)
		}
		if rl.IsKeyReleased(rl.KeyEnter) {
			n.SetButtonStatus(nes.ButtonStart, false)
		}
		if rl.IsKeyReleased(rl.KeyZ) {
			n.SetButtonStatus(nes.ButtonA, false)
		}
		if rl.IsKeyReleased(rl.KeyX) {
			n.SetButtonStatus(nes.ButtonB, false)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)
		rl.DrawTextureEx(texture, rl.NewVector2(0, 0), 0, float32(scale), rl.White)

		rl.EndDrawing()
	}

	return nil
}
