package game

import "log"
import "math/rand"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

func newSimpleGame(params Params, startRound int, ss score.Storage) core.Game {
	log.Println("new game:", startRound, params)
	game := simpleGame{ss: ss, round: startRound}
	game.Ticker = NewTicker(&game)
	game.Router = NewRouter(map[string]interface{}{
		"SCORE":        game.cmdScore,
		"CAT":          game.cmdCat,
		"ADD":          game.cmdAdd,
		"MUL":          game.cmdMul,
		"TURN":         game.Ticker.CmdTurn,
		"SUPER_SECRET": game.cmdSuperSecret,
	})
	return &game
}

type simpleGame struct {
	*Router
	*Ticker
	round  int
	ss     score.Storage
	secret int
}

func (self *simpleGame) NextRoundLength() int {
	return rand.Intn(10) + 1
}

func (self *simpleGame) roundStart() {
	log.Printf("starting round %d\n", self.round)
}

func (self *simpleGame) RoundEnd() {
	self.ss.TakeSnapshot()
	self.round++
	self.roundStart()
}

func (self *simpleGame) Turn(turn, length int) {
	log.Printf("turn %d/%d of round %d\n", turn, length, self.round)
	if turn == 1 {
		self.ss.SyncScores()
	}
	self.secret = rand.Intn(1000) + 1
}

func (self *simpleGame) cmdScore(team core.Team, score float64) core.CommandResult {
	if score == float64(self.secret) {
		self.ss.Scored(team, score)
		return core.NewOkResult()
	} else {
		return core.NewErrResult(101, "too greedy")
	}
}

func (self *simpleGame) cmdCat(team core.Team, a string, b string) core.CommandResult {
	return core.NewOkResult([]interface{}{a + b})
}

func (self *simpleGame) cmdAdd(team core.Team, a, b int) core.CommandResult {
	return core.NewOkResult([]interface{}{a + b})
}

func (self *simpleGame) cmdMul(team core.Team, a, b float64) core.CommandResult {
	return core.NewOkResult([]interface{}{a*b + 1})
}

func (self *simpleGame) cmdSuperSecret(team core.Team) core.CommandResult {
	return core.NewOkResult([]interface{}{self.secret})
}
