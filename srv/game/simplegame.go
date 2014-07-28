package game

import "log"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

type simpleGame struct {
	turn int
}

func (game *simpleGame) Execute(team core.Team, cmd core.Command) core.CommandResult {
	if cmd.Name != "DUPA" {
		return core.CommandResult{nil, []interface{}{cmd.Name, cmd.Params}}
	} else {
		return core.CommandResult{&core.CommandError{Desc: "spadaj"}, nil}
	}
}

func (game *simpleGame) Tick() {
	game.turn++
	log.Println("turn", game.turn)
}

func newSimpleGame(config *core.Config, ss score.Storage) core.Game {
	return new(simpleGame)
}
