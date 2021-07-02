package game

import (
	"bufio"
	"os"
	"time"
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
	Search
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

type Pos struct {
	X, Y int
}

type Entity struct {
	Pos
}

type Player struct {
	Entity
}

type Level struct {
	Map    [][]Tile
	Player Player
	Debug  map[Pos]bool
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

func canWalk(level *Level, pos Pos) bool {
	t := level.Map[pos.Y][pos.X]
	switch t {
	case StoneWall, ClosedDoor, Blank:
		checkDoor(level, pos)
		return false
	default:
		return true
	}
}

func checkDoor(level *Level, pos Pos) {
	t := level.Map[pos.Y][pos.X]
	if t == ClosedDoor {
		level.Map[pos.Y][pos.X] = OpenDoor
	}
}

func handleInput(ui GameUI, level *Level, input *Input) {
	p := level.Player
	switch input.Typ {
	case Up:
		if canWalk(level, Pos{p.X, p.Y - 1}) {
			level.Player.Y--
		}
	case Down:
		if canWalk(level, Pos{p.X, p.Y + 1}) {
			level.Player.Y++
		}
	case Left:
		if canWalk(level, Pos{p.X - 1, p.Y}) {
			level.Player.X--
		}
	case Right:
		if canWalk(level, Pos{p.X + 1, p.Y}) {
			level.Player.X++
		}
	case Search:
		bfs(ui, level, level.Player.Pos)
	case Quit:
		return
	}
}

func getNeighbors(level *Level, current Pos) []Pos {
	left := Pos{current.X - 1, current.Y}
	right := Pos{current.X + 1, current.Y}
	up := Pos{current.X, current.Y - 1}
	down := Pos{current.X, current.Y + 1}

	neighbors := make([]Pos, 0, 4)
	if canWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if canWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if canWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if canWalk(level, down) {
		neighbors = append(neighbors, down)
	}
	return neighbors
}

func bfs(ui GameUI, level *Level, start Pos) {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool, 0)
	visited[start] = true
	level.Debug = visited
	for len(frontier) > 0 {
		current := frontier[0]
		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
				ui.Draw(level)
				time.Sleep(100 * time.Millisecond)
			}
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

		handleInput(ui, level, input)
	}
}
