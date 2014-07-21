package main

import "log"

type SimpleGame struct {
	round int
}

func (game *SimpleGame) Execute(team Team, cmd Command) CommandResult {
	if cmd.Name != "DUPA" {
		return CommandResult{nil, []interface{}{cmd.Name, cmd.Params}}
	} else {
		return CommandResult{&CommandError{Desc: "spadaj"}, nil}
	}
}

func (game *SimpleGame) Tick() {
	game.round++
	log.Println("round", game.round)
}
