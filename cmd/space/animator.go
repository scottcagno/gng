package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"path"
	"time"
)

type animator struct {
	container       *element
	sequences       map[string]*sequence
	current         string // sequence name
	lastFrameChange time.Time
	finished        bool
}

func newAnimator(container *element, sequences map[string]*sequence, defaultSequence string) *animator {
	return &animator{
		container:       container,
		sequences:       sequences,
		current:         defaultSequence,
		lastFrameChange: time.Now(),
	}
}

func (a *animator) onUpdate() error {
	seq := a.sequences[a.current]
	frameInterval := float64(time.Second) / seq.sampleRate
	if time.Since(a.lastFrameChange) >= time.Duration(frameInterval) {
		a.finished = seq.nextFrame()
		a.lastFrameChange = time.Now()
	}
	return nil
}

func (a *animator) onDraw(renderer *sdl.Renderer) error {
	tex := a.sequences[a.current].texture()
	return drawTexture(tex, a.container.position, a.container.rotation, renderer)
}

func (a *animator) onCollision(other *element) error {
	return nil
}

func (a *animator) setSequence(name string) {
	a.current = name
	a.lastFrameChange = time.Now()
}

type sequence struct {
	textures   []*sdl.Texture
	frame      int
	sampleRate float64 // num of times a frame should play per second
	loop       bool
}

func newSequence(filepath string, sampleRate float64, loop bool, renderer *sdl.Renderer) (*sequence, error) {
	files, err := os.ReadDir(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading directory %v: %v\n", filepath, err)
	}
	var seq sequence
	for _, file := range files {
		filename := path.Join(filepath, file.Name())
		tex, err := loadTextureFromBMP(renderer, filename)
		if err != nil {
			return nil, fmt.Errorf("loading sequence frame: %v\n", err)
		}
		seq.textures = append(seq.textures, tex)
	}
	seq.sampleRate = sampleRate
	seq.loop = loop

	return &seq, nil
}

func (s *sequence) texture() *sdl.Texture {
	return s.textures[s.frame]
}

func (s *sequence) nextFrame() bool {
	if s.frame == len(s.textures)-1 {
		if s.loop {
			s.frame = 0
		} else {
			return true
		}
	} else {
		s.frame++
	}
	return false
}
