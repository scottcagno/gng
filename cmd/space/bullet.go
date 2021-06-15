package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

const (
	bulletSize  = 25
	bulletSpeed = 0.7
)

type bullet struct {
	tex   *sdl.Texture
	x, y  float64
	angle float64

	active bool
}

func newBullet(renderer *sdl.Renderer) *bullet {
	return &bullet{
		tex: textureFromBMP(renderer, "cmd/space/sprites/bullet.bmp"),
		x:   screenWidth / 2.0,
		y:   screenHeight - playerSize/2.0,
	}
}

func (b *bullet) draw(renderer *sdl.Renderer) {
	if !b.active {
		return
	}
	// converting coordinates from top left corner, to center of sprite
	x := b.x - bulletSize/2.0
	y := b.y - bulletSize/2.0
	renderer.Copy(b.tex,
		&sdl.Rect{X: 0, Y: 0, W: bulletSize, H: bulletSize},
		&sdl.Rect{X: int32(x), Y: int32(y), W: bulletSize, H: bulletSize},
	)
}

func (b *bullet) update() {
	b.x += bulletSpeed * math.Cos(b.angle)
	b.y += bulletSpeed * math.Sin(b.angle)
	if b.x > screenWidth || b.x < 0 || b.y > screenHeight || b.y < 0 {
		b.active = false
	}
}

var bulletPool []*bullet

func initBulletPoll(renderer *sdl.Renderer, max int) {
	for i := 0; i < max; i++ {
		bulletPool = append(bulletPool, newBullet(renderer))
	}
}

func bulletFromPool() (*bullet, bool) {
	for _, b := range bulletPool {
		if !b.active {
			return b, true
		}
	}
	return nil, false
}
