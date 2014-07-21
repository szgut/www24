package main

type throttler struct {
	game  Game
	limit int
	used  map[Team]int
}

func (t *throttler) Execute(team Team, cmd Command) CommandResult {
	t.used[team]++
	if t.used[team] > t.limit {
		return CommandResult{CommandLimitReachedError(), nil}
	}
	return t.game.Execute(team, cmd)
}

func (t *throttler) Tick() {
	t.used = make(map[Team]int)
	t.game.Tick()
}

func Throttler(limit int, game Game) Game {
	return &throttler{game: game, limit: limit, used: make(map[Team]int)}
}
