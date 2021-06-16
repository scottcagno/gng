package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"reflect"
)

type vector struct {
	x, y float64
}

type component interface {
	onUpdate() error
	onDraw(renderer *sdl.Renderer) error
	onCollision(other *element) error
}

type element struct {
	position   vector
	rotation   float64
	active     bool
	tag        string
	collisions []circle
	components []component
}

func (e *element) addComponent(comp component) {
	typ := reflect.TypeOf(comp)
	for _, existing := range e.components {
		if reflect.TypeOf(existing) == typ {
			panic(fmt.Sprintf("attempt to add new component with existing type %v\n", typ))
		}
	}
	e.components = append(e.components, comp)
}

func (e *element) getComponent(withType component) component {
	typ := reflect.TypeOf(withType)
	for _, existing := range e.components {
		if reflect.TypeOf(existing) == typ {
			return existing
		}
	}
	panic(fmt.Sprintf("no components with type %v\n", typ))
}

func (e *element) draw(renderer *sdl.Renderer) error {
	for _, comp := range e.components {
		err := comp.onDraw(renderer)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *element) update() error {
	for _, comp := range e.components {
		err := comp.onUpdate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *element) collision(other *element) error {
	for _, comp := range e.components {
		err := comp.onCollision(other)
		if err != nil {
			return err
		}
	}
	return nil
}

var elements []*element
