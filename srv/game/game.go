package game

import "log"
import "github.com/szgut/www24/srv/core"

type SimpleGame struct {
	turn int
}

func (game *SimpleGame) Execute(team core.Team, cmd core.Command) core.CommandResult {
	if cmd.Name != "DUPA" {
		return core.CommandResult{nil, []interface{}{cmd.Name, cmd.Params}}
	} else {
		return core.CommandResult{&core.CommandError{Desc: "spadaj"}, nil}
	}
}

func (game *SimpleGame) Tick() {
	game.turn++
	log.Println("turn", game.turn)
}
