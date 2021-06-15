package main

import "github.com/veandco/go-sdl2/sdl"

const basicEnemySize = 100

type enemy = basicEnemy

type basicEnemy struct {
	tex  *sdl.Texture
	x, y float64 // use float for smoother movement maths
}

func newBasicEnemy(renderer *sdl.Renderer, x, y float64) *basicEnemy {
	return &basicEnemy{
		tex: textureFromBMP(renderer, "cmd/space/sprites/square.bmp"),
		x:   x, //screenWidth / 2.0,
		y:   y, //screenHeight - playerSize/2.0,
	}
}

func (e *basicEnemy) draw(renderer *sdl.Renderer) {
	// converting coordinates from top left corner, to center of sprite
	x := e.x - basicEnemySize/2.0
	y := e.y - basicEnemySize/2.0
	renderer.CopyEx(e.tex,
		&sdl.Rect{X: 0, Y: 0, W: basicEnemySize, H: basicEnemySize},
		&sdl.Rect{X: int32(x), Y: int32(y), W: basicEnemySize, H: basicEnemySize},
		180, // angle of rotation
		&sdl.Point{X: basicEnemySize / 2.0, Y: basicEnemySize / 2.0},
		sdl.FLIP_NONE,
	)
}
