package main

import (
	"runtime"

	"github.com/scottcagno/gng/cmd/rpg/game"
	"github.com/scottcagno/gng/cmd/rpg/ui2d"
)

// using the tutorials at: https://www.youtube.com/watch?v=ZxG90DtNs3k&list=PLDZujg-VgQlZUy1iCqBbe5faZLMkA3g2x&index=29

func main() {
	numWindows := 1
	rpgGame := game.NewGame(numWindows, "cmd/rpg/game/maps/level1.map")
	for i := 0; i < numWindows; i++ {
		go func(i int) {
			runtime.LockOSThread()
			ui := ui2d.NewUI(rpgGame.InputChan, rpgGame.LevelChans[i])
			ui.Run()
		}(i)
	}
	rpgGame.Run()
}
