package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"
	"runtime"
	"time"

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
	img  *image.RGBA
	rimg *rl.Image
}

func newRenderer(img *image.RGBA) *renderer {
	return &renderer{
		img:  img,
		rimg: rl.NewImageFromImage(img),
	}
}

func (r *renderer) Render(x, y int, c color.Color) {
	r.img.Set(x, y, c)
}
func (r *renderer) Refresh() {
	r.rimg = rl.NewImageFromImage(r.img)
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
		rom   string
		scale int
	)
	flag.StringVar(&rom, "rom", "", "rome filepath")
	flag.IntVar(&scale, "scale", 2, "window scale size")
	flag.Parse()

	img := image.NewRGBA(image.Rect(0, 0, 256, 240))
	renderer := newRenderer(img)

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

	go func() {
		time.Sleep(2 * time.Second)
		n.PowerUp()
		n.Run()
	}()

	rl.SetTraceLogLevel(rl.LogWarning)
	rl.InitWindow(256*int32(scale), 240*int32(scale), "rgnes")
	defer rl.CloseWindow()

	texture := rl.LoadTextureFromImage(renderer.rimg)
	defer rl.UnloadTexture(texture)

	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		newTexture := rl.LoadTextureFromImage(renderer.rimg)
		rl.UnloadTexture(texture)
		texture = newTexture

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
