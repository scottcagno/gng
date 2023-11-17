package main

import "github.com/veandco/go-sdl2/sdl"

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
}
