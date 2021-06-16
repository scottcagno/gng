package main

import "github.com/veandco/go-sdl2/sdl"

type vulnToBullets struct {
	container *element
}

func newVulnToBullets(container *element) *vulnToBullets {
	return &vulnToBullets{
		container: container,
	}
}

func (v *vulnToBullets) onDraw(renderer *sdl.Renderer) error {
	return nil
}

func (v *vulnToBullets) onUpdate() error {
	return nil
}

func (v *vulnToBullets) onCollision(other *element) error {
	if other.tag == "bullet" {
		v.container.active = false
	}
	return nil
}
