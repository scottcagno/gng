package game

import (
	"bufio"
	"os"
)

type GameUI interface {
	Draw(*Level)
}

type Tile rune

const (
	StoneWall Tile = '#'
	DirtFloor Tile = '.'
	Door      Tile = '|'
	Blank     Tile = 0
)

type Level struct {
	Map [][]Tile
}

func loadLevelFromFile(filename string) *Level {
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	levelLines := make([]string, 0)
	longestRow := 0
	var index int
	for scanner.Scan() {
		levelLines = append(levelLines, scanner.Text())
		if len(levelLines[index]) > longestRow {
			longestRow = len(levelLines[index])
		}
		index++
	}
	level := &Level{
		Map: make([][]Tile, len(levelLines)),
	}
	for i := range level.Map {
		level.Map[i] = make([]Tile, longestRow)
	}

	for y := 0; y < len(level.Map); y++ {
		line := levelLines[y]
		for x, c := range line {
			var t Tile
			switch c {
			case ' ', '\t', '\n', '\r':
				t = Blank
			case '#':
				t = StoneWall
			case '|':
				t = Door
			case '.':
				t = DirtFloor
			}
			level.Map[y][x] = t
		}
	}
	return level
}

func Run(ui GameUI) {
	level := loadLevelFromFile("cmd/rpg/game/maps/level1.map")
	ui.Draw(level)
}
