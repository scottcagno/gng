package main

import "github.com/veandco/go-sdl2/sdl"

type vulnToBullets struct {
	container *element
	animator  *animator
}

func newVulnToBullets(container *element) *vulnToBullets {
	return &vulnToBullets{
		container: container,
		animator:  container.getComponent(&animator{}).(*animator),
	}
}

func (v *vulnToBullets) onDraw(renderer *sdl.Renderer) error {
	return nil
}

func (v *vulnToBullets) onUpdate() error {
	if v.animator.finished && v.animator.current == "destroy" {
		v.container.active = false
	}
	return nil
}

func (v *vulnToBullets) onCollision(other *element) error {
	if other.tag == "bullet" {
		v.animator.setSequence("destroy")
	}
	return nil
}
