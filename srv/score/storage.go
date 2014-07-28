package score

import "fmt"
import "github.com/szgut/www24/srv/core"

type Storage interface {
	Scored(team core.Team, change float64) error
	SyncScores()
	TakeSnapshot()
}

type storage struct {
	db     DB
	task   string
	scores map[core.Team]float64
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
	for team, score := range self.scores {
		self.db.Exec("update score_teams set score = ? where team = ?", score, team.String())
	}
}

func (self *storage) TakeSnapshot() {

}

func InitializeDatabase(dbPath string, task string, teams []core.Team) {
	db := ConnectDB(dbPath)
	db.Exec("delete from score_teams where task = ?", task)
	for _, team := range teams {
		db.Exec("insert into score_teams(team, task) values(?,?)", team.String(), task)
	}
}

func NewStorage(dbPath string, task string) Storage {
	return &storage{db: ConnectDB(dbPath), task: task}
}
