package main

type Team struct {
	login string
}

type Command struct {
	Name   string
	Params []string
}

type CommandError struct {
	Id   int
	Desc string
}

type Game interface {
	Execute(team Team, cmd Command) (params []interface{}, err *CommandError)
	Tick()
}
