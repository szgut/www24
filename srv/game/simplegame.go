package game

import "log"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

type simpleGame struct {
	round int
	turn int
	ss score.Storage
}

func (self *simpleGame) Execute(team core.Team, cmd core.Command) core.CommandResult {
	if cmd.Name != "DUPA" {
		return core.CommandResult{nil, []interface{}{cmd.Name, cmd.Params}}
	} else {
		self.ss.Scored(team, 10)
		return core.CommandResult{&core.CommandError{Desc: "spadaj"}, nil}
	}
}

func (self *simpleGame) Tick() {
	self.turn++
	if self.turn == 10 {
		self.turn = 0
		self.round++
		self.ss.TakeSnapshot()
	}
	self.ss.SyncScores()
	log.Printf("tick %d/%d\n", self.turn, self.round)
}

func newSimpleGame(params Params, startRound int, ss score.Storage) core.Game {
	log.Println("new game:", startRound, params)
	return &simpleGame{ss: ss, round: startRound}
}
