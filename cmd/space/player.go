package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	playerSpeed        = 0.4
	playerSize         = 100
	playerShotCooldown = time.Millisecond * 250
)

func newPlayer(renderer *sdl.Renderer) *element {
	p := &element{
		position: vector{
			x: screenWidth / 2.0,
			y: screenHeight - playerSize/2.0,
		},
		active: true,
	}
	sr := newSpriteRenderer(p, renderer, "cmd/space/sprites/triangle.bmp")
	p.addComponent(sr)

	mover := newKeyboardMover(p, playerSpeed)
	p.addComponent(mover)

	shooter := newKeyboardShooter(p, playerShotCooldown)
	p.addComponent(shooter)

	return p
}
