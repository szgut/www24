package game

import "log"
import "math"
import "math/rand"
import "time"
import "strconv"
import "os"
import "fmt"
import "github.com/szgut/www24/srv/core"
import "github.com/szgut/www24/srv/score"

var gameLogger *log.Logger

func init() {
	f, err := os.OpenFile("world", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("error opening file:", err)
	}
	gameLogger = log.New(f, "", log.Ltime|log.Lmicroseconds)
}

const STAR_SCORE_SCALE = 1.1

var (
	UNIT_COST int
	//BASE_COST int

	MAX_MANA  int
	SCAN_MANA int

	MAX_ROUNDS   int
	MAX_DIAMETER int
)

func readParam(params Params, name string) int {
	value, ok := params[name]
	if !ok {
		log.Fatal("no game param", name)
	}
	i, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	return int(i)
}

func setGameParams(params Params) {
	UNIT_COST = readParam(params, "unit_cost")
	//BASE_COST = readParam(params, "base_cost")
	MAX_MANA = readParam(params, "max_mana")
	SCAN_MANA = readParam(params, "scan_mana")

	MAX_ROUNDS = readParam(params, "max_rounds")
	MAX_DIAMETER = readParam(params, "max_diameter")
}

func newStarGame(params Params, startRound int, teams []core.Team, ss score.Storage) core.Game {
	log.Println("new game:", startRound, params)
	rand.Seed(time.Now().UTC().UnixNano())
	setGameParams(params)
	game := StarGame{ss: ss, round: startRound, logins: teams}
	game.Router = NewRouter(map[string]interface{}{
		"TURN":           game.cmdTurn,
		"DESCRIBE_WORLD": game.cmdDescribeWorld,
		"ME":             game.cmdMe,
		"ENEMIES":        game.cmdEnemies,
		"ATTACK":         game.cmdAttack,
		"SCAN":           game.cmdScan,
		"BUY_PEON":       game.cmdBuyPeon,
		"BUY_ARMY":       game.cmdBuyArmy,
		"ATTACKS":        game.cmdAttacks,
	})
	game.roundStart()
	return &game
}

type StarGame struct {
	*Router
	ss     score.Storage
	round  int
	turn   int
	length int

	logins    []core.Team
	unit_cost int
	max_mana  int
	scan_mana int

	roundScoreScale float64
	diameter        int
	players         map[core.Team]*Player
	attacks         map[int][]Attack

	//basesUsed       map[core.Team]int
	attackResults map[core.Team][]AttackResult
}

type AttackResult struct {
	Attack
	killed int
}

type Attack struct {
	from     core.Team
	to       core.Team
	strength int
}

func (self Attack) String() string {
	return fmt.Sprintf("%s->%s(%d)", self.from.String(), self.to.String(), self.strength)
}

type Player struct {
	peons int
	army  int
	//bases    int
	minerals int
	mana     int
}

func NewPlayer() *Player {
	return &Player{peons: 4, army: 0, minerals: UNIT_COST, mana: SCAN_MANA}
}

func (self *Player) totalUnits() int {
	return self.peons + self.army
}

func (self *StarGame) Tick() {
	self.turn++
	if self.turn > self.length {
		self.turn = 1
		self.roundEnd()
	}
	self.turnStart()
}

func random(a, b int) int {
	return rand.Intn(b-a) + a
}

func (self *StarGame) roundStart() {
	self.roundScoreScale = math.Pow(STAR_SCORE_SCALE, float64(self.round-1))
	self.length = random(500, MAX_ROUNDS)
	self.diameter = random(1, MAX_DIAMETER)
	log.Printf("starting round %d, diameter=%d\n", self.round, self.diameter)
	self.attacks = make(map[int][]Attack)
	self.players = make(map[core.Team]*Player)
	for _, team := range self.logins {
		self.players[team] = NewPlayer()
	}
}

func (self *StarGame) roundEnd() {
	self.ss.TakeSnapshot()
	self.round++
	self.roundStart()
}

func (self *StarGame) turnStart() {
	if self.turn%4 == 0 {
		self.ss.SyncScores()
	}
	log.Printf("turn %d/%d of round %d, diameter %d\n", self.turn, self.length, self.round, self.diameter)

	self.attackResults = make(map[core.Team][]AttackResult)
	for _, attack := range self.attacks[self.turn] {
		self.executeAttack(attack)
	}
	delete(self.attacks, self.turn)

	for _, player := range self.players {
		player.mana = min(player.mana+1, MAX_MANA)
		player.minerals += player.peons
	}

	gameLogger.Printf("turn %d/%d\n", self.turn, self.length)
	for _, team := range self.logins {
		gameLogger.Printf("%15v\t%+v\n", team, *self.players[team])
	}
}

func fight(a, b *int) {
	if *a < *b {
		a, b = b, a
	}
	*a -= *b * *b / *a
	*b = 0
}

func (self *StarGame) executeAttack(attack Attack) {
	gameLogger.Println("executing attack", attack)
	astr := attack.strength
	defender := self.players[attack.to]
	unitsBefore := defender.totalUnits()
	fight(&astr, &defender.army)
	defender.peons = max(0, defender.peons-int(math.Pow(float64(astr), 1.25)))
	killed := unitsBefore - defender.totalUnits()
	gameLogger.Printf("player %s kills %d scoring %f\n", attack.from.String(), killed, self.score(killed))
	self.ss.Scored(attack.from, self.score(killed))
	self.attackResults[attack.to] = append(self.attackResults[attack.to], AttackResult{Attack: attack, killed: killed})
}

func (self *StarGame) score(size int) float64 {
	return float64(size) * self.roundScoreScale
}

func (self *StarGame) cmdTurn(team core.Team) core.CommandResult {
	return core.NewOkResult([]interface{}{self.turn, self.length})
}

func (self *StarGame) cmdDescribeWorld(team core.Team) core.CommandResult {
	return core.NewOkResult([]interface{}{self.diameter})
}

func (self *StarGame) cmdMe(team core.Team) core.CommandResult {
	player := self.players[team]
	return core.NewOkResult([]interface{}{player.peons, player.army, player.minerals, player.mana})
}

func (self *StarGame) cmdEnemies(team core.Team) core.CommandResult {
	enemies := []interface{}{}
	for _, login := range self.logins {
		if login != team {
			enemies = append(enemies, login.String())
		}
	}
	return core.NewOkResult([]interface{}{len(enemies)}, enemies)
}

func (self *StarGame) scheduleAttack(attack Attack) {
	gameLogger.Println("scheduling attack", attack)
	when := self.turn + self.diameter
	self.attacks[when] = append(self.attacks[when], attack)
}

func (self *StarGame) cmdAttack(team core.Team, targetStr string, strength int) core.CommandResult {
	if strength <= 0 {
		return core.ErrResult(core.BadFormatError())
	}
	target := core.NewTeam(targetStr)
	if team == target {
		return core.NewErrResult(102, "unable to attack yourself")
	}
	_, ok := self.players[target]
	if !ok {
		return core.NewErrResult(101, "no such player")
	}
	player := self.players[team]
	if player.army < strength {
		return core.NewErrResult(104, "not enough army")
	}
	player.army -= strength
	self.scheduleAttack(Attack{from: team, to: target, strength: strength})
	return core.NewOkResult()
}

func (self *StarGame) cmdScan(team core.Team, target string) core.CommandResult {
	enemy, ok := self.players[core.NewTeam(target)]
	if !ok {
		return core.NewErrResult(101, "no such player")
	}
	player := self.players[team]
	if player.mana < SCAN_MANA {
		return core.NewErrResult(103, "not enough energy")
	}
	player.mana -= SCAN_MANA
	return core.NewOkResult([]interface{}{enemy.peons, enemy.army})
}

func (self *StarGame) cmdBuyPeon(team core.Team, n int) core.CommandResult {
	return self.cmdBuy(team, false, n)
}

func (self *StarGame) cmdBuyArmy(team core.Team, n int) core.CommandResult {
	return self.cmdBuy(team, true, n)
}

func (self *StarGame) cmdBuy(team core.Team, army bool, n int) core.CommandResult {
	if n <= 0 {
		return core.ErrResult(core.BadFormatError())
	}
	player := self.players[team]
	if player.minerals < UNIT_COST*n {
		return core.NewErrResult(105, "you have not enough minerals")
	}
	player.minerals -= UNIT_COST * n
	if army {
		player.army += n
	} else {
		player.peons += n
	}
	return core.NewOkResult()
}

func (self *StarGame) cmdAttacks(team core.Team) core.CommandResult {
	lines := [][]interface{}{[]interface{}{len(self.attackResults[team])}}
	for _, ar := range self.attackResults[team] {
		lines = append(lines, []interface{}{ar.from.String(), ar.strength, ar.killed})
	}
	return core.NewOkResult(lines...)
}

//func (self *StarGame) cmdBuyBase(team)
