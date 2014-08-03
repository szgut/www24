package game

import "log"
import "strconv"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

type simpleGame struct {
	round int
	turn  int
	ss    score.Storage
}

func (self *simpleGame) Execute(team core.Team, cmd core.Command) core.CommandResult {
	if cmd.Name == "DUPA" {
		return core.NewErrResult(core.CommandError{105, "spadaj"})
	} else if cmd.Name == "SCORE" {
		if len(cmd.Params) != 1 {
			return core.NewErrResult(core.BadFormatError())
		}
		score, err := strconv.ParseFloat(cmd.Params[0], 64);
		if err != nil {
			return core.NewErrResult(core.BadFormatError())
		}
		self.ss.Scored(team, score)
		return core.NewOkResult()
	} else {
		return core.NewOkResult([]interface{}{cmd.Name, cmd.Params})
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
