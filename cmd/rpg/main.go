package main

import (
	"github.com/scottcagno/gng/cmd/rpg/game"
	"github.com/scottcagno/gng/cmd/rpg/ui2d"
)

func main() {
	game.Run(new(ui2d.UI2d))
}
