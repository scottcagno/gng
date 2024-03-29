package game

type Monster struct {
	Pos
	Rune         rune
	Name         string
	Hitpoints    int
	Strength     int
	Speed        float32
	ActionPoints float32
	// Character
}

func NewRat(p Pos) *Monster {
	return &Monster{p, 'R', "Rat", 5, 5, 2.0, 0}
}

func NewSpider(p Pos) *Monster {
	return &Monster{p, 'S', "Spider", 10, 10, 1.0, 0}
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos

	apInt := int(m.ActionPoints)

	positions := level.astar(m.Pos, playerPos)
	moveIndex := 1
	for i := 0; i < apInt; i++ {
		// must be > 1 because the 1st position is the monsters current position
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]

	// todo: check if the tile being moved is valid
	// todo: if player is in the way, attack player
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
	}
}

/*
func NewRat(p Pos) *Monster {
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'R'
	monster.Name = "Rat"
	monster.Hitpoints = 5
	monster.Strength = 0
	monster.Speed = 2.0
	monster.ActionPoints = 0.0
	monster.SightRange = 10
	return monster
}

func NewSpider(p Pos) *Monster {
	monster := &Monster{}
	monster.Pos = p
	monster.Rune = 'S'
	monster.Name = "Spider"
	monster.Hitpoints = 1000
	monster.Strength = 5
	monster.Speed = 1.0
	monster.ActionPoints = 0.0
	monster.SightRange = 10
	return monster
}

func (m *Monster) Kill(level *Level) {
	delete(level.Monsters, m.Pos)
	groundItems := level.Items[m.Pos]
	for _, item := range m.Items {
		item.Pos = m.Pos
		groundItems = append(groundItems, item)
	}
	level.Items[m.Pos] = groundItems
}

func (m *Monster) Update(level *Level) {
	m.ActionPoints += m.Speed
	playerPos := level.Player.Pos
	apInt := int(m.ActionPoints)
	positions := level.astar(m.Pos, playerPos)

	//Do we have any path to the goal?
	if len(positions) == 0 {
		m.Pass()
		return
	}

	moveIndex := 1
	for i := 0; i < apInt; i++ {
		if moveIndex < len(positions) {
			m.Move(positions[moveIndex], level)
			moveIndex++
			m.ActionPoints--
		}
	}
}

func (m *Monster) Pass() {
	m.ActionPoints -= m.Speed
}

func (m *Monster) Move(to Pos, level *Level) {
	_, exists := level.Monsters[to]
	// TODO check if the tile being moved to is valid
	if !exists && to != level.Player.Pos {
		delete(level.Monsters, m.Pos)
		level.Monsters[to] = m
		m.Pos = to
		return
	}

	if to == level.Player.Pos {
		level.Attack(&m.Character, &level.Player.Character)
		if m.Hitpoints <= 0 {
			delete(level.Monsters, m.Pos)
		}
		if level.Player.Hitpoints <= 0 {
			panic("ded")
		}
	}

}
*/
