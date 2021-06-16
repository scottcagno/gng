package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
)

const basicEnemySize = 100

func newBasicEnemy(renderer *sdl.Renderer, position vector) *element {
	enemy := &element{
		position: position,
		rotation: 0,
		active:   true,
	}

	idleSeq, err := newSequence("cmd/space/sprites/enemy/idle", 5, true, renderer)
	if err != nil {
		panic(fmt.Errorf("creating idle sequence: %v\n", err))
	}

	destroySeq, err := newSequence("cmd/space/sprites/enemy/destroy", 10, false, renderer)
	if err != nil {
		panic(fmt.Errorf("creating destroy sequence: %v\n", err))
	}

	sequences := map[string]*sequence{
		"idle":    idleSeq,
		"destroy": destroySeq,
	}

	an := newAnimator(enemy, sequences, "idle")
	enemy.addComponent(an)

	vuln := newVulnToBullets(enemy)
	enemy.addComponent(vuln)

	enemy.collisions = append(enemy.collisions, newCollisionCircle(enemy.position, basicEnemySize))

	return enemy
}
