package game

import (
	"bufio"
	"os"
)

type GameUI interface {
	Draw(*Level)
	GetInput() *Input
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	Quit
)

type Input struct {
	Typ InputType
}

type Tile rune

const (
	StoneWall  Tile = '#'
	DirtFloor  Tile = '.'
	ClosedDoor Tile = '|'
	OpenDoor   Tile = '/'
	Blank      Tile = 0
	Pending    Tile = -1
)

type Entity struct {
	X, Y int
}

type Player struct {
	Entity
}

type Level struct {
	Map    [][]Tile
	Player Player
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
				t = ClosedDoor
			case '/':
				t = OpenDoor
			case '.':
				t = DirtFloor
			case 'P':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			default:
				panic("Invalid character in map")
			}
			level.Map[y][x] = t
		}
	}

	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
			SearchLoop:
				for searchX := x - 1; searchX <= x+1; searchX++ {
					for searchY := y - 1; searchY <= y+1; searchY++ {
						searchTile := level.Map[searchY][searchX]
						switch searchTile {
						case DirtFloor:
							level.Map[y][x] = DirtFloor
							break SearchLoop
						}
					}
				}
			}
		}
	}
	return level
}

func canWalk(level *Level, x, y int) bool {
	t := level.Map[y][x]
	switch t {
	case StoneWall, ClosedDoor, Blank:
		checkDoor(level, x, y)
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, x, y int) {
	t := level.Map[y][x]
	if t == ClosedDoor {
		level.Map[y][x] = OpenDoor
	}
}

func handleInput(level *Level, input *Input) {
	p := level.Player
	switch input.Typ {
	case Up:
		if canWalk(level, p.X, p.Y-1) {
			level.Player.Y--
		}
	case Down:
		if canWalk(level, p.X, p.Y+1) {
			level.Player.Y++
		}
	case Left:
		if canWalk(level, p.X-1, p.Y) {
			level.Player.X--
		}
	case Right:
		if canWalk(level, p.X+1, p.Y) {
			level.Player.X++
		}
	}
}

func Run(ui GameUI) {
	level := loadLevelFromFile("cmd/rpg/game/maps/level1.map")

	// game loop
	var input *Input
	for {
		ui.Draw(level)
		input = ui.GetInput()

		if input != nil && input.Typ == Quit {
			return
		}

		handleInput(level, input)
	}
}
