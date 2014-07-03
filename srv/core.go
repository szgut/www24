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

func AuthenticationFailedError() *CommandError {
	return &CommandError{1, "bad login or password"}
}
func UnknownCommandError() *CommandError {
	return &CommandError{2, "unknown command"}
}
func BadFormatError() *CommandError {
	return &CommandError{3, "bad format"}
}
func CommandLimitReachedError() *CommandError {
	return &CommandError{6, "commands limit reached, forced waiting activated"}
}
