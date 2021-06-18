package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math"
	"time"
)

const (
	winWidth  = 800
	winHeight = 600
)

type gameState int

const (
	start gameState = iota
	play
)

var state = start

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xvel   float32
	yvel   float32
	color  color
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	if num > 9 {
		num = 0
	}

	for i, v := range numbers[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

func (b *ball) draw(pixels []byte) {
	// YAGNI - ya aint gonna need it
	for y := -b.radius; y < b.radius; y++ {
		for x := -b.radius; x < b.radius; x++ {
			// avoid sqrt because it's expensive
			if x*x+y*y < b.radius*b.radius {
				setPixel(int(b.x+x), int(b.y+y), b.color, pixels)
			}
		}
	}
}

func getCenter() pos {
	return pos{winWidth / 2, winHeight / 2}
}

func (b *ball) update(p1, p2 *paddle, gameSpeed float32) {
	b.x += b.xvel * gameSpeed
	b.y += b.yvel * gameSpeed

	// ball colliding (bounces) off top OR bottom of screen
	if b.y-b.radius < 0 || b.y+b.radius > winHeight {
		b.yvel = -b.yvel
	}

	// ball colliding (scores) off left of screen
	if b.x-b.radius < 0 {
		p2.score++
		b.pos = getCenter()
		state = start
	}

	// ball colliding (scores) off right of screen
	if b.x+b.radius > winWidth {
		p1.score++
		b.pos = getCenter()
		state = start
	}

	// collision detection with player 1, paddle 1 (left)
	if b.x-b.radius < p1.x+p1.w/2 {
		if b.y > p1.y-p1.h/2 && b.y < p1.y+p1.h/2 {
			b.xvel = -b.xvel // reverse ball.x velocity
			b.x = p1.x + p1.w/2 + b.radius
		}
	}

	// collision detection with player 2, paddle 2 (right)
	if b.x+b.radius > p2.x-p2.w/2 {
		if b.y > p2.y-p2.h/2 && b.y < p2.y+p2.h/2 {
			b.xvel = -b.xvel // reverse ball.x velocity
			b.x = p2.x - p2.w/2 - b.radius
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	yvel  float32
	score int
	color color
}

// linear interpolation
func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func (p *paddle) draw(pixels []byte) {
	startX := p.x - p.w/2
	startY := p.y - p.h/2

	for y := 0; y < int(p.h); y++ {
		for x := 0; x < int(p.w); x++ {
			setPixel(int(startX)+x, int(startY)+y, p.color, pixels)
		}
	}

	numX := lerp(p.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, p.color, 10, p.score, pixels)
}

func (p *paddle) update(keys []uint8, controllerAxis int16, elapsedTime float32) {
	keyUP := sdl.GetScancodeFromName("UP")
	keyDN := sdl.GetScancodeFromName("DOWN")
	if keys[keyUP] != 0 {
		p.y -= p.yvel * elapsedTime
	}

	if keys[keyDN] != 0 {
		p.y += p.yvel * elapsedTime
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
		p.y += p.yvel * pct * elapsedTime
	}
}

func (p *paddle) aiUpdate(ball *ball, gameSpeed float32) {
	p.y = ball.y
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func main() {

	// init sdl
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Panicf("init sdl: %v", err)
	}
	defer sdl.Quit()

	// create window
	window, err := sdl.CreateWindow(
		"Testing SDL2",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth,
		winHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		log.Panicf("window: %v", err)
	}
	defer window.Destroy()

	// create renderer
	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Panicf("renderer: %v", err)
	}
	defer renderer.Destroy()

	// create texture
	tex, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		winWidth,
		winHeight,
	)
	if err != nil {
		log.Panicf("texture: %v", err)
	}
	defer tex.Destroy()

	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	// create pixel matrix
	pixels := make([]byte, winWidth*winHeight*4)

	// create new player
	p1 := &paddle{
		pos:   pos{50, 100},
		w:     20,
		h:     100,
		yvel:  400,
		color: color{255, 255, 255},
	}

	// create new player
	p2 := &paddle{
		pos:   pos{winWidth - 50, 100},
		w:     20,
		h:     100,
		yvel:  400,
		color: color{255, 255, 255},
	}

	// create new ball
	b1 := &ball{
		pos:    pos{300, 300},
		radius: 20,
		xvel:   400,
		yvel:   400,
		color:  color{255, 255, 255},
	}

	// grab ptr to sdl's internally managed keyboard events
	keys := sdl.GetKeyboardState()

	var frameStart time.Time
	var elapsedTime float32
	var controllerAxis int16

	// main game loop
	//
	for {
		frameStart = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		for _, controller := range controllerHandlers {
			if controller != nil {
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}

		if state == play {
			p1.update(keys, controllerAxis, elapsedTime)
			p2.aiUpdate(b1, elapsedTime)
			b1.update(p1, p2, elapsedTime)
		}

		if state == start {
			if keys[sdl.SCANCODE_SPACE] != 0 {
				if p1.score == 3 || p2.score == 3 {
					p1.score = 0
					p2.score = 0
				}
				state = play
			}
		}

		clear(pixels) // clear screen
		p1.draw(pixels)
		p2.draw(pixels)
		b1.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		// ensure we get 200fps??
		if elapsedTime < 0.005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}

	}

}
