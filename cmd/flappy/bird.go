package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"sync"
)

const (
	gravity   = 0.15
	jumpSpeed = 5
)

type bird struct {
	mu sync.RWMutex

	time     int
	textures []*sdl.Texture

	x, y int32
	w, h int32

	speed float64
	dead  bool
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
		x:        10,
		y:        300,
		w:        50,
		h:        43,
		speed:    0,
	}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	rect := &sdl.Rect{X: b.x, Y: 600 - b.y - b.h/2, W: b.w, H: b.h}

	i := b.time / 10 % len(b.textures)
	err := r.Copy(b.textures[i], nil, rect)
	if err != nil {
		return fmt.Errorf("could not copy bird: %v", err)
	}
	return nil
}

func (b *bird) destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, t := range b.textures {
		t.Destroy()
	}
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.dead
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.speed = -jumpSpeed
}

// collision detection
func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if p.x > b.x+b.w {
		return // too far right, not touching
	}
	if p.x+p.w < b.x {
		return // too far left, not touching
	}
	if !p.inverted && p.h < b.y-b.h/2 {
		return // pipe is too low, not touching
	}
	if p.inverted && 600-p.h > b.y+b.h/2 {
		return // inverted pipe is too high, not touching
	}
	// otherwise, it's touching
	b.dead = true
}
