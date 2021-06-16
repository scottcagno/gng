package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

func keyIsPressed(ok uint8) bool {
	return ok == 1
}

type keyboardMover struct {
	container *element
	speed     float64
	sr        *spriteRenderer
}

func newKeyboardMover(container *element, speed float64) *keyboardMover {
	return &keyboardMover{
		container: container,
		speed:     speed,
		sr:        container.getComponent(&spriteRenderer{}).(*spriteRenderer),
	}
}

func (mover *keyboardMover) onUpdate() error {
	keys := sdl.GetKeyboardState()
	cont := mover.container
	if keyIsPressed(keys[sdl.SCANCODE_LEFT]) {
		if cont.position.x-(mover.sr.width/2.0) > 0 {
			cont.position.x -= mover.speed
		}

	} else if keyIsPressed(keys[sdl.SCANCODE_RIGHT]) {
		if cont.position.x+(mover.sr.height/2.0) < screenWidth {
			cont.position.x += mover.speed
		}
	}
	return nil
}

func (mover *keyboardMover) onDraw(renderer *sdl.Renderer) error {
	return nil
}

func (mover *keyboardMover) onCollision(other *element) error {
	return nil
}

type keyboardShooter struct {
	container *element
	cooldown  time.Duration
	lastShot  time.Time
}

func newKeyboardShooter(container *element, cooldown time.Duration) *keyboardShooter {
	return &keyboardShooter{
		container: container,
		cooldown:  cooldown,
	}
}

func (shooter *keyboardShooter) shoot(x, y float64) {
	if bul, ok := bulletFromPool(); ok {
		bul.active = true
		bul.position.x = x
		bul.position.y = y
		bul.rotation = 270 * (math.Pi / 180)
	}
}

func (shooter *keyboardShooter) onUpdate() error {
	keys := sdl.GetKeyboardState()
	if keyIsPressed(keys[sdl.SCANCODE_SPACE]) {
		if time.Since(shooter.lastShot) >= shooter.cooldown {
			pos := shooter.container.position
			shooter.shoot(pos.x+25, pos.y)
			shooter.shoot(pos.x-25, pos.y)
			shooter.lastShot = time.Now()
		}
	}
	return nil
}

func (shooter *keyboardShooter) onDraw(renderer *sdl.Renderer) error {
	return nil
}

func (shooter *keyboardShooter) onCollision(other *element) error {
	return nil
}
