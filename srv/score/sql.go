package score

import "github.com/szgut/www24/srv/core"

//TODO: transactions
//TODO: embedding, task

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

func (self *TaskQueries) UpdateScores(scores map[core.Team]float64, snapshot int) {
	// TODO
	for team, score := range scores {
		self.Exec("update score_teams set score = ? where team = ?, snapshot = ?", score, team.String(), snapshot)
	}
}

func (self *TaskQueries) TakeSnapshot(snapshot int) {
	self.Exec("insert into score(id, team, task, score, snapshot)"+
		"select null, team, task, score, ?, from score where task = ? and snapshot = 0", snapshot, self.task)
}

func (self *TaskQueries) Initialize() {
	self.Exec("delete from score where task = ?", self.task)
}

func (self *TaskQueries) InitTeams(teams []core.Team) {
	for _, team := range teams {
		self.Exec("insert into score_teams(team, task) values(?,?)", team.String(), self.task)
	}
}
