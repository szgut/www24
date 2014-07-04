package main

import "log"

type SimpleGame struct {
	round int
}

func (game *SimpleGame) Execute(team Team, cmd Command) (params []interface{}, err *CommandError) {
	if cmd.Name != "DUPA" {
		return []interface{}{cmd.Name, cmd.Params}, nil
	} else {
		return nil, &CommandError{Desc: "spadaj"}
	}
}

func (game *SimpleGame) Tick() {
	game.round++
	log.Println("round", game.round)
}
