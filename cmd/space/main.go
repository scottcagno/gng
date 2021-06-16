package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
	"time"
)

const (
	screenWidth          = 600
	screenHeight         = 800
	targetTicksPerSecond = 60
)

type vector struct {
	x, y float64
}

var delta float64

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		handleErr("initializing sdl", err)
	}
	window, err := sdl.CreateWindow(
		"My Game",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		screenWidth,
		screenHeight,
		sdl.WINDOW_SHOWN,
	)
	handleErr("initializing window", err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	handleErr("initializing renderer", err)
	defer renderer.Destroy()

	elements = append(elements, newPlayer(renderer)) // GLOBAL

	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			offset := +(basicEnemySize / 2.0) + 1.0
			x := (float64(i)/5)*screenWidth + offset
			y := float64(j)*basicEnemySize + offset
			enemy := newBasicEnemy(renderer, vector{x: x, y: y})
			handleErr("creating basic enemy at "+strconv.Itoa(i)+", "+strconv.Itoa(j), err)

			elements = append(elements, enemy) // GLOBAL
		}
	}

	// initialize bullet pool
	initBulletPool(renderer, 30)

	// start event loop
	running := true
	for running {
		frameStart := time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()

		for _, elem := range elements {
			if elem.active {
				err = elem.update()
				handleErr("updating element", err)

				err = elem.draw(renderer)
				handleErr("drawing element", err)
			}
		}

		err = checkCollisions()
		handleErr("checking collisions", err)

		renderer.Present()

		delta = time.Since(frameStart).Seconds() * targetTicksPerSecond

	}
}
