package game

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// TODO:
// change the way key presses work
// make speed work properly
// combat
// cleanup? items?

// type GameUI interface {
// 	Draw(*Level)
// 	GetInput() *Input
// }

type Game struct {
	LevelChans []chan *Level
	InputChan  chan *Input
	Level      *Level
}

func NewGame(numWindows int, path string) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)
	level := loadLevelFromFile(path)
	return &Game{
		LevelChans: levelChans,
		InputChan:  inputChan,
		Level:      level,
	}
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	QuitGame
	CloseWindow
	Search
)

type Input struct {
	Typ          InputType
	LevelChannel chan *Level
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

type Character struct {
	Entity
	Hitpoints    int
	Strength     int
	Speed        float64
	ActionPoints float64
	SightRange   int
	Items        []*Item
	Helmet       *Item
	Weapon       *Item
}

type Player struct {
	Entity
}

type Level struct {
	Map      [][]Tile
	Player   Player
	Monsters map[Pos]*Monster
	Debug    map[Pos]bool
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
		Map:      make([][]Tile, len(levelLines)),
		Monsters: make(map[Pos]*Monster),
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
			case '@':
				level.Player.X = x
				level.Player.Y = y
				t = Pending
			case 'R':
				level.Monsters[Pos{x, y}] = NewRat(Pos{x, y})
				t = Pending
			case 'S':
				level.Monsters[Pos{x, y}] = NewSpider(Pos{x, y})
				t = Pending
			default:
				panic("Invalid character in map")
			}
			level.Map[y][x] = t
		}
	}

	// TODO: we should use bfs to find the first floor tile
	for y, row := range level.Map {
		for x, tile := range row {
			if tile == Pending {
				// SearchLoop:
				// 	for searchX := x - 1; searchX <= x+1; searchX++ {
				// 		for searchY := y - 1; searchY <= y+1; searchY++ {
				// 			searchTile := level.Map[searchY][searchX]
				// 			switch searchTile {
				// 			case DirtFloor:
				// 				level.Map[y][x] = DirtFloor
				// 				break SearchLoop
				// 			}
				// 		}
				// 	}
				// }
				level.Map[y][x] = level.bfsFloor(Pos{x, y})
			}
		}
	}
	return level
}

func inRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Map[0]) && pos.Y < len(level.Map) && pos.X >= 0 && pos.Y >= 0
}

func canWalk(level *Level, pos Pos) bool {
	if !inRange(level, pos) {
		return false
	}
	t := level.Map[pos.Y][pos.X]
	switch t {
	case StoneWall, ClosedDoor, Blank:
		// checkDoor(level, pos)
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

func (player *Player) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	if !exists {
		player.Pos = to
	}

}

func (game *Game) handleInput(input *Input) {
	level := game.Level
	p := game.Level.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		if canWalk(level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		if canWalk(game.Level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		if canWalk(game.Level, newPos) {
			level.Player.Move(newPos, level)
		} else {
			checkDoor(level, newPos)
		}
	case Search:
		// game.bfs(game.Level.Player.Pos)
		level.astar(level.Player.Pos, Pos{3, 2})
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChans {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		game.LevelChans = append(game.LevelChans[:chanIndex], game.LevelChans[chanIndex+1:]...)
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

func (level *Level) bfsFloor(start Pos) Tile {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true
	// level.Debug = visited

	for len(frontier) > 0 {
		current := frontier[0]

		currentTile := level.Map[current.Y][current.X]
		switch currentTile {
		case DirtFloor:
			return DirtFloor
		default:
		}

		frontier = frontier[1:]
		for _, next := range getNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
				// time.Sleep(100 * time.Millisecond)
			}
		}
	}

	return DirtFloor
}

func (level *Level) astar(start Pos, goal Pos) []Pos {
	frontier := make(pqueue, 0, 8)
	frontier = frontier.push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	var current Pos
	for len(frontier) > 0 {

		frontier, current = frontier.pop()

		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)

			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			return path
		}

		for _, next := range getNeighbors(level, current) {
			newCost := costSoFar[current] + 1 // always 1 for now
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.push(next, priority)
				cameFrom[next] = current
			}

		}

	}

	return nil
}

func (game *Game) Run() {

	fmt.Printf("Starting RPG...\n")

	// send game state updates
	for _, lchan := range game.LevelChans {
		lchan <- game.Level
	}

	// game loop
	for input := range game.InputChan {

		// handle input signals
		if input.Typ == QuitGame {
			return
		}

		// handle inputs
		game.handleInput(input)

		// update monsters
		for _, monster := range game.Level.Monsters {
			monster.Update(game.Level)
		}

		// all windows are closed
		if len(game.LevelChans) == 0 {
			return
		}

		// send gamestate updates
		for _, lchan := range game.LevelChans {
			lchan <- game.Level
		}
	}
}
