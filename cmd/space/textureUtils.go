package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

func drawTexture(tex *sdl.Texture, position vector, rotation float64, renderer *sdl.Renderer) error {
	_, _, width, height, err := tex.Query()
	if err != nil {
		panic(fmt.Errorf("query texture: %v\n", err))
	}
	// converting coordinates from top left corner, to center of sprite
	position.x -= float64(width) / 2.0
	position.y -= float64(height) / 2.0
	return renderer.CopyEx(tex,
		&sdl.Rect{X: 0, Y: 0, W: width, H: height},
		&sdl.Rect{X: int32(position.x), Y: int32(position.y), W: width, H: height},
		rotation, // angle of rotation
		&sdl.Point{X: width / 2, Y: height / 2},
		sdl.FLIP_NONE,
	)
}

func loadTextureFromBMP(renderer *sdl.Renderer, filename string) (*sdl.Texture, error) {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		return nil, fmt.Errorf("loading %v: %v", filename, err)
	}
	defer img.Free()
	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		return nil, fmt.Errorf("creating texture from %v: %v", filename, err)
	}
	return tex, nil
}
