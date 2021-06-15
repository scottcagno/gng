package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"strconv"
)

const (
	screenWidth  = 600
	screenHeight = 800
)

func textureFromBMP(renderer *sdl.Renderer, filename string) *sdl.Texture {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("loading %v: %v", filename, err))
	}
	defer img.Free()
	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture from %v: %v", filename, err))
	}
	return tex
}

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

	// create new player instance
	plr := newPlayer(renderer)

	// create matrix of enemies
	var enemies []*basicEnemy
	for i := 0; i < 5; i++ {
		for j := 0; j < 3; j++ {
			offset := +(basicEnemySize / 2.0) + 1.0
			x := (float64(i)/5)*screenWidth + offset
			y := float64(j)*basicEnemySize + offset
			enemy := newBasicEnemy(renderer, x, y)
			handleErr("creating basic enemy at "+strconv.Itoa(i)+", "+strconv.Itoa(j), err)
			enemies = append(enemies, enemy)
		}
	}

	// initialize bullet pool
	initBulletPoll(renderer, 30)

	// start event loop
	running := true
	for running {
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

		plr.draw(renderer)
		plr.update() // update player location

		for _, enemy := range enemies {
			enemy.draw(renderer)
		}

		for _, bul := range bulletPool {
			bul.draw(renderer)
			bul.update()
		}

		renderer.Present()
	}
}
