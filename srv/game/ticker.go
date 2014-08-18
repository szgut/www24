package game

import "github.com/szgut/www24/srv/core"

type Ticker struct {
	length int
	turn   int
	game   TickerCallback
}

type TickerCallback interface {
	NextRoundLength() int
	Turn(turn, length int)
	RoundEnd()
}

func NewTicker(game TickerCallback) *Ticker {
	return &Ticker{game: game, length: game.NextRoundLength()}
}

func (self *Ticker) Tick() {
	self.turn++
	if self.turn > self.length {
		self.game.RoundEnd()
		self.turn = 1
		self.length = self.game.NextRoundLength()
	}
	self.game.Turn(self.turn, self.length)
}

func (self *Ticker) CmdTurn(team core.Team) core.CommandResult {
	return core.NewOkResult([]interface{}{self.turn, self.length})
}
