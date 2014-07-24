package score

import "fmt"
import "github.com/szgut/www24/srv/core"

type ScoreStorage interface {
	Scored(team core.Team, change float64) error
	FlushScores()
	TakeSnapshot(round int)
}

type scoreStorage struct {
	db     DB
	task   string
	scores map[core.Team]float64
}

func InitializeDatabase(dbPath string, task string, teams []core.Team) {
	db := ConnectDB(dbPath)
	db.Exec("delete from score_teams where task = ?", task)
	for _, team := range teams {
		db.Exec("insert into score_teams(team, task) values(?,?)", team.String(), task)
	}
}

func NewScoreStorage(dbPath string, task string) ScoreStorage {
	return &scoreStorage{db: ConnectDB(dbPath), task: task}
}

func (ss *scoreStorage) Scored(team core.Team, change float64) error {
	_, ok := ss.scores[team]
	if !ok {
		return fmt.Errorf("Scored: attempt to increase score of nonexistent team %s", team)
	}
	ss.scores[team] += change
	return nil
}

func (ss *scoreStorage) UpdateScores(scores map[core.Team]float64) {
	for team, score := range scores {
		ss.db.Exec("update score_teams set score = ? where team = ?", score, team.String())
	}
}

func (ss *scoreStorage) TakeSnapshot(round int) {

}

func (ss *scoreStorage) readSnapshot() {

}

func (ss *scoreStorage) FlushScores() {

}
