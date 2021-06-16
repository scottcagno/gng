package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const basicEnemySize = 100

func newBasicEnemy(renderer *sdl.Renderer, position vector) *element {
	enemy := &element{
		position: position,
		rotation: 180,
		active:   true,
	}
	sr := newSpriteRenderer(enemy, renderer, "cmd/space/sprites/square.bmp")
	enemy.addComponent(sr)

	vuln := newVulnToBullets(enemy)
	enemy.addComponent(vuln)

	enemy.collisions = append(enemy.collisions, newCollisionCircle(enemy.position, basicEnemySize))

	return enemy
}
