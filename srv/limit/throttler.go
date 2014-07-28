package limit

import "github.com/szgut/www24/srv/core"

type throttler struct {
	game  core.Game
	limit int
	used  map[core.Team]int
}

func (self *throttler) Execute(team core.Team, cmd core.Command) core.CommandResult {
	self.used[team]++
	if self.used[team] > self.limit {
		return core.CommandResult{core.CommandLimitReachedError(), nil}
	}
	return self.game.Execute(team, cmd)
}

func (self *throttler) Tick() {
	self.used = make(map[core.Team]int)
	self.game.Tick()
}

func Throttler(limit int, game core.Game) core.Game {
	return &throttler{game: game, limit: limit, used: make(map[core.Team]int)}
}
