package main

import (
	"github.com/scottcagno/gng/cmd/rpg/game"
	"github.com/scottcagno/gng/cmd/rpg/ui2d"
	"runtime"
)

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
