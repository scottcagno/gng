package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
	"runtime"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

const (
	screenWidth  = 800
	screenHeight = 600
)

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return fmt.Errorf("could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	err = ttf.Init()
	if err != nil {
		return fmt.Errorf("could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	w, r, err := sdl.CreateWindowAndRenderer(screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		return fmt.Errorf("could not create window: %v", err)
	}
	defer w.Destroy()

	err = drawTitle(r, "Flappy Gopher")
	if err != nil {
		return fmt.Errorf("could not draw title: %v", err)
	}
	time.Sleep(1 * time.Second)

	s, err := newScene(r)
	if err != nil {
		return fmt.Errorf("could not create scene: %v", err)
	}
	defer s.destroy()

	events := make(chan sdl.Event)
	errc := s.run(events, r)

	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

func drawTitle(r *sdl.Renderer, text string) error {
	r.Clear() // clear buffer

	font, err := ttf.OpenFont("cmd/flappy/res/fonts/flappy.ttf", 10)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer font.Close()

	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	sur, err := font.RenderUTF8Solid(text, c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer sur.Free()

	tex, err := r.CreateTextureFromSurface(sur)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer tex.Destroy()

	_, _, w, h, err := tex.Query()
	if err != nil {
		return fmt.Errorf("could not query texture: %v", err)
	}

	w, h = w*4, h*4

	err = r.Copy(tex, nil, &sdl.Rect{
		X: screenWidth/2 - w/2,
		Y: screenHeight/2 - h/2,
		W: w,
		H: h,
	})
	if err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present() // paint buffer
	return nil
}

func drawTexture(tex *sdl.Texture, x, y int32, rotation float64, renderer *sdl.Renderer) error {
	_, _, width, height, err := tex.Query()
	if err != nil {
		return fmt.Errorf("query texture: %v\n", err)
	}
	// converting coordinates from top left corner, to center of sprite
	x -= width / 2.0
	y -= height / 2.0
	return renderer.CopyEx(tex,
		&sdl.Rect{X: 0, Y: 0, W: width, H: height},
		&sdl.Rect{X: x, Y: y, W: width, H: height},
		rotation, // angle of rotation
		&sdl.Point{X: width / 2, Y: height / 2},
		sdl.FLIP_NONE,
	)
}
