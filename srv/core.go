package main

type Team struct {
	login string
}

func (team Team) String() string {
	return team.login
}

type Command struct {
	Name   string
	Params []string
}

type CommandResult struct {
	Err    *CommandError
	Params []interface{}
}

type CommandError struct {
	Id   int
	Desc string
}

type Game interface {
	Execute(team Team, cmd Command) CommandResult
	Tick()
}

func (err *CommandError) ShouldWait() bool {
	return err != nil && err.Id == 6
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
