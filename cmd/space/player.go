package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"time"
)

const (
	playerSpeed        = 0.4
	playerSize         = 100
	playerShotCooldown = time.Millisecond * 250
)

type player struct {
	tex      *sdl.Texture
	x, y     float64 // use float for smoother movement maths
	lastShot time.Time
}

func newPlayer(renderer *sdl.Renderer) *player {
	return &player{
		tex: textureFromBMP(renderer, "cmd/space/sprites/triangle.bmp"),
		x:   screenWidth / 2.0,
		y:   screenHeight - playerSize/2.0,
	}
}

func (p *player) draw(renderer *sdl.Renderer) {
	// converting coordinates from top left corner, to center of sprite
	x := p.x - playerSize/2.0
	y := p.y - playerSize/2.0
	renderer.Copy(p.tex,
		&sdl.Rect{X: 0, Y: 0, W: playerSize, H: playerSize},
		&sdl.Rect{X: int32(x), Y: int32(y), W: playerSize, H: playerSize},
	)
}

func keyIsPressed(ok uint8) bool {
	return ok == 1
}

func (p *player) update() {
	keys := sdl.GetKeyboardState()
	if keyIsPressed(keys[sdl.SCANCODE_LEFT]) {
		if p.x-(playerSize/2.0) > 0 {
			p.x -= playerSpeed
		}

	} else if keyIsPressed(keys[sdl.SCANCODE_RIGHT]) {
		if p.x+(playerSize/2.0) < screenWidth {
			p.x += playerSpeed
		}
	}
	if keyIsPressed(keys[sdl.SCANCODE_SPACE]) {
		if time.Since(p.lastShot) >= playerShotCooldown {
			p.shoot(p.x+25, p.y)
			p.shoot(p.x-25, p.y)
			p.lastShot = time.Now()
		}
	}

}

func (p *player) shoot(x, y float64) {
	if bul, ok := bulletFromPool(); ok {
		bul.active = true
		bul.x = x
		bul.y = y
		bul.angle = 270 * (math.Pi / 180)
	}
}
