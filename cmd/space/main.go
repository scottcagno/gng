package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

const (
	screenWidth  = 600
	screenHeight = 800
)

func handleErr(msg string, err error) {
	log.Panicf("%s: %v\n", msg, err)
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
		sdl.WINDOW_OPENGL,
	)
	handleErr("initializing window", err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	handleErr("initializing renderer", err)
	defer renderer.Destroy()

	for {
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		renderer.Present()
	}
}
