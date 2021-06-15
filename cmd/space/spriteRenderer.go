package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type spriteRenderer struct {
	container *element
	tex       *sdl.Texture
}

func newSpriteRenderer(contianer *element, renderer *sdl.Renderer, filename string) *spriteRenderer {
	return &spriteRenderer{
		container: contianer,
		tex:       textureFromBMP(renderer, filename),
	}
}

func (s *spriteRenderer) onUpdate() error {
	return nil
}

func (s *spriteRenderer) onDraw(renderer *sdl.Renderer) error {
	_, _, width, height, err := s.tex.Query()
	if err != nil {
		return fmt.Errorf("query texture: %v\n", err)
	}
	// converting coordinates from top left corner, to center of sprite
	x := s.container.position.x - float64(width)/2.0
	y := s.container.position.y - float64(height)/2.0

	renderer.CopyEx(s.tex,
		&sdl.Rect{X: 0, Y: 0, W: width, H: height},
		&sdl.Rect{X: int32(x), Y: int32(y), W: width, H: height},
		s.container.rotation, // angle of rotation
		&sdl.Point{X: width / 2, Y: height / 2},
		sdl.FLIP_NONE,
	)
	return nil
}

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
