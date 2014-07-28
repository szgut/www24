package score

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
		tx.Exec("insert into score set score = ? where team = ?, snapshot = ?", score, team.String(), snapshot)
	}
}

func (self *TaskQueries) ReadScores(snapshot int) map[core.Team]float64 {
	return nil
}

func (self *TaskQueries) Clear() {
	self.Exec("delete from score where task = ?", self.task)
}
