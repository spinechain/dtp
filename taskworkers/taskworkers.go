package taskworkers

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// This stores all known workers (globally) as well as how to reach them

type TaskWorkers struct {
	db          *sql.DB
	TaskWorkers []*TaskWorker
}

type TaskWorker struct {
	ID                   string
	Address              string
	Port                 int
	TasksDone            int
	Reputation           float64
	TasksInQueue         int
	AvgCompletionTime    float64
	MinimumFee           float64
	Deadness             uint
	Capabilities         string
	LastActive           int
	TasksDoneLast24Hours int
	BestConnections      string // IDs of the nodes that can easily reach it
}

func (t *TaskWorkers) Start(filePath string, create bool) error {
	var err error
	t.db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}

	if create {
		t.CreateTable()
	}
	return err
}

func (t *TaskWorkers) CreateTable() error {
	sqlStmt := `
	create table taskworkers (tid integer not null primary key, description text);
	delete from taskworkers;
	`
	_, err := t.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (t *TaskWorkers) Stop() {
	if t.db != nil {
		t.db.Close()
	}
}

func (t *TaskWorkers) AddTaskWorker() {

}
