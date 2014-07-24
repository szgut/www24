package backend

import "github.com/szgut/www24/srv/core"

type throttler struct {
	game  core.Game
	limit int
	used  map[core.Team]int
}

func (t *throttler) Execute(team core.Team, cmd core.Command) core.CommandResult {
	t.used[team]++
	if t.used[team] > t.limit {
		return core.CommandResult{core.CommandLimitReachedError(), nil}
	}
	return t.game.Execute(team, cmd)
}

func (t *throttler) Tick() {
	t.used = make(map[core.Team]int)
	t.game.Tick()
}

func Throttler(limit int, game core.Game) core.Game {
	return &throttler{game: game, limit: limit, used: make(map[core.Team]int)}
}
