package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	bulletSize  = 25
	bulletSpeed = 0.7
)

func newBullet(renderer *sdl.Renderer) *element {
	bullet := &element{
		active: false,
	}
	sr := newSpriteRenderer(bullet, renderer, "cmd/space/sprites/bullet.bmp")
	bullet.addComponent(sr)

	mover := newBulletMover(bullet, bulletSpeed)
	bullet.addComponent(mover)

	return bullet
}

var bulletPool []*element

func initBulletPoll(renderer *sdl.Renderer, max int) {
	for i := 0; i < max; i++ {
		bul := newBullet(renderer)
		elements = append(elements, bul) // GLOBAL
		bulletPool = append(bulletPool, bul)
	}
}

func bulletFromPool() (*element, bool) {
	for _, b := range bulletPool {
		if !b.active {
			return b, true
		}
	}
	return nil, false
}
