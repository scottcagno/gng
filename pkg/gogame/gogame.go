package gogame

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

var delta float64

type Game struct {
	win      *sdl.Window
	ren      *sdl.Renderer
	ticks    float64
	elements []interface{}
}

func NewGame(title string, width int32, height int32) (*Game, error) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return nil, fmt.Errorf("could not initialize sdl: %v", err)
	}

	win, err := sdl.CreateWindow(
		"My Game",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		width,
		height,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return nil, fmt.Errorf("could not initialize window: %v", err)
	}

	ren, err := sdl.CreateRenderer(win, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return nil, fmt.Errorf("could not initialize renderer: %v", err)
	}

	return &Game{
		win:   win,
		ren:   ren,
		ticks: 60,
	}, nil
}

func (g *Game) Run() error {
	for {
		frameStart := time.Now()
		g.ren.Clear()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return nil
			}
		}

		for _, elem := range g.elements {
			if elem.active {
				err = elem.update()
				handleErr("updating element", err)

				err = elem.draw(renderer)
				handleErr("drawing element", err)
			}
		}

		g.ren.Present()
		delta = time.Since(frameStart).Seconds() * g.ticks
	}
}

func (g *Game) Destroy() {
	g.win.Destroy()
	g.ren.Destroy()
}
