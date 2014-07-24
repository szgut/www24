package main

import "log"

type SimpleGame struct {
	turn int
}

func (game *SimpleGame) Execute(team Team, cmd Command) CommandResult {
	if cmd.Name != "DUPA" {
		return CommandResult{nil, []interface{}{cmd.Name, cmd.Params}}
	} else {
		return CommandResult{&CommandError{Desc: "spadaj"}, nil}
	}
}

func (game *SimpleGame) Tick() {
	game.turn++
	log.Println("turn", game.turn)
}
