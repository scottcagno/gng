package main

import "github.com/veandco/go-sdl2/sdl"

const (
	playerSpeed = 0.2
	playerSize  = 100
)

type player struct {
	tex  *sdl.Texture
	x, y float64 // use float for smoother movement maths
}

func newPlayer(renderer *sdl.Renderer) (*player, error) {
	img, err := sdl.LoadBMP("cmd/space/sprites/triangle.bmp")
	if err != nil {
		return nil, err
	}
	defer img.Free()
	ptex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		return nil, err
	}
	return &player{
		tex: ptex,
		x:   screenWidth / 2.0,
		y:   screenHeight - playerSize/2.0,
	}, nil
}

func (p *player) draw(renderer *sdl.Renderer) {
	// converting coordinates from top left corner, to center of sprite
	x := p.x - playerSize/2.0
	y := p.y - playerSize/2.0
	renderer.Copy(p.tex,
		&sdl.Rect{X: 0, Y: 0, W: 100, H: 100},
		&sdl.Rect{X: int32(x), Y: int32(y), W: 100, H: 100},
	)
}

func keyIsPressed(ok uint8) bool {
	return ok == 1
}

func (p *player) update() {
	keys := sdl.GetKeyboardState()
	if keyIsPressed(keys[sdl.SCANCODE_LEFT]) {
		p.x -= playerSpeed
		return
	}
	if keyIsPressed(keys[sdl.SCANCODE_RIGHT]) {
		p.x += playerSpeed
		return
	}
}
