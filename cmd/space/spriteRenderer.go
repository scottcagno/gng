package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

type spriteRenderer struct {
	container     *element
	tex           *sdl.Texture
	width, height float64
}

func newSpriteRenderer(contianer *element, renderer *sdl.Renderer, filename string) *spriteRenderer {
	tex, err := loadTextureFromBMP(renderer, filename)
	if err != nil {
		panic(fmt.Errorf("query texture: %v\n", err))
	}
	_, _, width, height, err := tex.Query()
	return &spriteRenderer{
		container: contianer,
		tex:       tex,
		width:     float64(width),
		height:    float64(height),
	}
}

func (sr *spriteRenderer) onUpdate() error {
	return nil
}

func (sr *spriteRenderer) onDraw(renderer *sdl.Renderer) error {
	return drawTexture(sr.tex, sr.container.position, sr.container.rotation, renderer)
}

func (sr *spriteRenderer) onCollision(other *element) error {
	return nil
}
