package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	playerSpeed        = 5
	playerSize         = 100
	playerShotCooldown = time.Millisecond * 250
)

func newPlayer(renderer *sdl.Renderer) *element {
	plyr := &element{
		position: vector{
			x: screenWidth / 2.0,
			y: screenHeight - playerSize/2.0,
		},
		active: true,
	}
	sr := newSpriteRenderer(plyr, renderer, "cmd/space/sprites/player.bmp")
	plyr.addComponent(sr)

	mover := newKeyboardMover(plyr, playerSpeed)
	plyr.addComponent(mover)

	shooter := newKeyboardShooter(plyr, playerShotCooldown)
	plyr.addComponent(shooter)

	return plyr
}
