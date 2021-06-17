package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"time"
)

type scene struct {
	bg   *sdl.Texture
	bird *bird
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "cmd/flappy/res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}
	bird, err := newBird(r)
	if err != nil {
		return nil, err
	}
	return &scene{
		bg:   bg,
		bird: bird,
	}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)

		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
			case <-tick:
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()
	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent, *sdl.TextInputEvent:
		s.bird.jump()
		return false
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent:
		return false
	default:
		log.Printf("unknown event %T", event)
		return false
	}
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	err := r.Copy(s.bg, nil, nil)
	if err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	err = s.bird.paint(r)
	if err != nil {
		return err
	}

	r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
}
