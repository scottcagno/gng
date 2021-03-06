package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	bulletSpeed = 10
	bulletSize  = 25
)

func newBullet(renderer *sdl.Renderer) *element {
	bullet := &element{
		active: false,
	}
	sr := newSpriteRenderer(bullet, renderer, "cmd/space/sprites/bullet.bmp")
	bullet.addComponent(sr)

	mover := newBulletMover(bullet, bulletSpeed)
	bullet.addComponent(mover)

	bullet.collisions = append(bullet.collisions, newCollisionCircle(bullet.position, bulletSize))
	bullet.tag = "bullet"

	return bullet
}

var bulletPool []*element

func initBulletPool(renderer *sdl.Renderer, max int) {
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
