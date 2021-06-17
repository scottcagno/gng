package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.25
	jumpSpeed = 5
)

type bird struct {
	time     int
	textures []*sdl.Texture
	y        float64
	speed    float64
}

func newBird(r *sdl.Renderer) (*bird, error) {
	var textures []*sdl.Texture
	for i := 1; i < 4; i++ {
		path := fmt.Sprintf("cmd/flappy/res/imgs/bird_frame_%d.png", i)
		texture, err := img.LoadTexture(r, path)
		if err != nil {
			return nil, fmt.Errorf("could not load bird (frame 1) image: %v", err)
		}
		textures = append(textures, texture)
	}
	return &bird{
		time:     0,
		textures: textures,
		y:        300,
		speed:    0,
	}, nil
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.time++
	b.y -= b.speed
	if b.y < 0 {
		b.speed = -b.speed
		b.y = 0
	}
	b.speed += gravity

	rect := &sdl.Rect{X: 10, Y: (600 - int32(b.y)) - 43/2, W: 50, H: 43}

	i := b.time / 10 % len(b.textures)
	err := r.Copy(b.textures[i], nil, rect)
	if err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}
	return nil
}

func (b *bird) destroy() {
	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *bird) jump() {
	b.speed = -jumpSpeed
}
