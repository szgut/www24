package game

import "log"
import "math"
import "math/rand"
import "time"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

const MAX_SIZE = 100
const SCAN_RANGE = 2
const SCORE_SCALE = 1.05

func newFieldsGame(params Params, startRound int, ss score.Storage) core.Game {
	log.Println("new game:", startRound, params)
	rand.Seed(time.Now().UTC().UnixNano())
	game := FieldsGame{ss: ss, round: startRound}
	game.Ticker = NewTicker(&game)
	game.Router = NewRouter(map[string]interface{}{
		"TURN":           game.Ticker.CmdTurn,
		"DESCRIBE_WORLD": game.cmdDescribeWorld,
		"BUY":            game.cmdBuy,
		"SELL":           game.cmdSell,
		"LAST_PURCHASES": game.cmdLastReservations,
		"SCAN":           game.cmdScan,
	})
	game.roundStart()
	return &game
}

type FieldsGame struct {
	*Router
	*Ticker
	round int
	ss    score.Storage

	owner           [][]core.Team
	n               int
	m               int
	soldSq          map[Square]bool
	roundScoreScale float64

	reservations     map[core.Team]Square
	lastReservations map[core.Team]Square
	//sold             map[core.Team]bool
}

type Square struct {
	x int
	y int
}

func (self *FieldsGame) NextRoundLength() int {
	// TODO
	return rand.Intn(1000) + 5
}

func (self *FieldsGame) roundStart() {
	self.roundScoreScale = math.Pow(SCORE_SCALE, float64(self.round))
	self.n = rand.Intn(MAX_SIZE) + 5
	self.m = rand.Intn(MAX_SIZE) + 5
	log.Printf("starting round %d, %dx%d\n", self.round, self.n, self.m)
	self.owner = make([][]core.Team, self.n)
	for i := range self.owner {
		self.owner[i] = make([]core.Team, self.m)
	}
	self.soldSq = make(map[Square]bool)
}

func (self *FieldsGame) RoundEnd() {
	self.ss.TakeSnapshot()
	self.round++
	self.roundStart()
}

func (self *FieldsGame) Turn(turn, length int) {
	log.Printf("turn %d/%d of round %d\n", turn, length, self.round)
	self.ss.SyncScores()

	self.lastReservations = self.reservations
	self.reservations = make(map[core.Team]Square)
	//self.sold = make(map[core.Team]bool)
}

func (self *FieldsGame) cmdDescribeWorld(team core.Team) core.CommandResult {
	return core.NewOkResult([]interface{}{self.n, self.m})
}

func (self *FieldsGame) isOk(x, y int) bool {
	return 1 <= x && x <= self.n && 1 <= y && y <= self.m
}

func (self *FieldsGame) isFree(x, y int) bool {
	return self.owner[x][y] == (core.Team{})
}

func (self *FieldsGame) cmdBuy(team core.Team, x, y int) core.CommandResult {
	if !self.isOk(x, y) {
		return core.NewErrResult(101, "not in range")
	}
	if _, done := self.reservations[team]; done {
		return core.NewErrResult(104, "purchase already made in this turn")
	}
	if !self.isFree(x, y) {
		return core.NewErrResult(102, "field already bought")
	}
	self.owner[x][y] = team
	self.reservations[team] = Square{x, y}
	return core.NewOkResult()
}

func (self *FieldsGame) cmdSell(team core.Team, x1, y1, x2, y2 int) core.CommandResult {
	if !self.isOk(x1, y1) || !self.isOk(x2, y2) {
		return core.NewErrResult(101, "not in range")
	}
	/*if self.sold[team] {
		return core.NewErrResult(103, "already sold in this turn")
	}
	self.sold[team] = true*/
	top := Square{min(x1, x2), min(y1, y2)}
	bottom := Square{max(x2, x2), max(y1, y2)}
	sum := 0
	for x := top.x; x <= bottom.x; x++ {
		for y := top.y; y <= bottom.y; y++ {
			if self.owner[x][y] != team {
				return core.NewErrResult(105, "field not owned")
			}
			sum++
		}
	}
	self.ss.Scored(team, float64(sum*sum)*self.roundScoreScale)
	for x := top.x; x <= bottom.x; x++ {
		for y := top.y; y <= bottom.y; y++ {
			self.soldSq[Square{x, y}] = true
		}
	}
	return core.NewOkResult()
}

func (self *FieldsGame) cmdLastReservations(team core.Team) core.CommandResult {
	lines := [][]interface{}{[]interface{}{len(self.lastReservations)}}
	for team, point := range self.lastReservations {
		lines = append(lines, []interface{}{team, point.x, point.y})
	}
	return core.NewOkResult(lines...)
}

func (self *FieldsGame) cmdScan(team core.Team, x, y int) core.CommandResult {
	owner := func(x, y int) string {
		if self.isFree(x, y) {
			return "-"
		} else {
			return self.owner[x][y].String()
		}
	}
	lines := [][]interface{}{}
	for i := x - SCAN_RANGE; i <= x+SCAN_RANGE; i++ {
		line := []interface{}{}
		for j := y - SCAN_RANGE; j <= y+SCAN_RANGE; j++ {
			if !self.isOk(i, j) {
				return core.NewErrResult(101, "not in range")
			}
			line = append(line, owner(i, j))
		}
		lines = append(lines, line)
	}
	return core.NewOkResult(lines...)
}
