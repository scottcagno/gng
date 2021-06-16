package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time  int64
	bg    *sdl.Texture
	birds []*sdl.Texture
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "cmd/flappy/res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	var birds []*sdl.Texture
	for i := 1; i < 4; i++ {
		path := fmt.Sprintf("cmd/flappy/res/imgs/bird_frame_%d.png", i)
		bird, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load bird (frame 1) image: %v", err)
		}
		birds = append(birds, bird)
	}

	return &scene{
		bg:    bg,
		birds: birds,
	}, nil
}

func (s *scene) paint(r *sdl.Renderer) error {
	s.time++
	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	rect := &sdl.Rect{X: 10, Y: screenHeight - 43/2, W: 50, H: 43}
	i := s.time % int64(len(s.birds))
	err = r.Copy(s.birds[i], nil, rect)
	if err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
}
