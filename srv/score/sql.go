package score

import "log"
import "github.com/szgut/www24/srv/core"

type TaskQueries struct {
	DB
	task string
}

func NewTaskQueries(path string, task string) TaskQueries {
	return TaskQueries{DB: ConnectDB(path), task: task}
}

func (self *TaskQueries) LastSnapshot() int {
	var snapshot int
	self.ScanQuery(&snapshot, "select max(snapshot) from score where task = ?", self.task)
	return snapshot
}

func (self *TaskQueries) WriteScores(scores map[core.Team]float64, snapshot int) {
	tx := self.Begin()
	defer tx.Commit()
	tx.Exec("delete from score where task = ? and snapshot = ?", self.task, snapshot)
	for team, score := range scores {
		tx.Exec("insert into score(id, task, snapshot, team, score) values(null, ?, ?, ?, ?)",
			self.task, snapshot, team.String(), score)
	}
}

func (self *TaskQueries) ReadScores(snapshot int) map[core.Team]float64 {
	rows := self.Query("select team, score from score where task = ? and snapshot = ?", self.task, snapshot)
	scores := make(map[core.Team]float64)
	for rows.Next() {
		var login string
		var score float64
		if err := rows.Scan(&login, &score); err != nil {
			log.Fatal(err)
		}
		scores[core.NewTeam(login)] = score
	}
	return scores
}

func (self *TaskQueries) Clear() {
	self.Exec("delete from score where task = ?", self.task)
}
