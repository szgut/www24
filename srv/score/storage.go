package score

import "fmt"
import "log"
import "github.com/szgut/www24/srv/core"

type Database interface {
	Storage
	Initialize(teams []core.Team)
	ReadScores()
	StartRound() int
}

type Storage interface {
	Scored(team core.Team, change float64) error
	SyncScores()
	TakeSnapshot()
}

type storage struct {
	queries TaskQueries
	scores  map[core.Team]float64
}

func (self *storage) Scored(team core.Team, change float64) error {
	_, ok := self.scores[team]
	if !ok {
		return fmt.Errorf("No such team: %s", team.String())
	}
	self.scores[team] += change
	return nil
}

func (self *storage) SyncScores() {
	self.queries.WriteScores(self.scores, 0)
}

func (self *storage) TakeSnapshot() {
	self.queries.WriteScores(self.scores, self.queries.LastSnapshot()+1)
}

func (self *storage) Initialize(teams []core.Team) {
	log.Println("Initializing database")
	self.scores = make(map[core.Team]float64)
	for _, team := range teams {
		self.scores[team] = 0
	}
	self.queries.Clear()
	self.SyncScores()
}

func (self *storage) ReadScores() {
	self.scores = self.queries.ReadScores(self.queries.LastSnapshot())
}

func (self *storage) StartRound() int {
	return self.queries.LastSnapshot() + 1
}

func NewStorage(dbPath string, task string) Database {
	return &storage{queries: NewTaskQueries(dbPath, task)}
}
