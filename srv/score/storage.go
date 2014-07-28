package score

import "fmt"
import "github.com/szgut/www24/srv/core"

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

}

func InitializeDatabase(dbPath string, task string, teams []core.Team) {
	db := ConnectDB(dbPath)
	db.Exec("delete from score_teams where task = ?", task)
	for _, team := range teams {
		db.Exec("insert into score_teams(team, task) values(?,?)", team.String(), task)
	}
}

func NewStorage(dbPath string, task string) Storage {
	return &storage{queries: NewTaskQueries(dbPath, task)}
}
