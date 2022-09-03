package taskworkers

import (
	"database/sql"
	"fmt"
	"spinedtp/util"

	_ "github.com/mattn/go-sqlite3"
)

// This stores all known workers (globally) as well as how to reach them

type TaskWorkers struct {
	db                 *sql.DB
	TaskWorkers        []*TaskWorker
	OnTaskWorkersAdded func(string) // ID
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

func TaskWorkerSqlColumnString() string {
	s := "wid string not null primary key, " +
		"address text, port int, taskdone int, reputation real, " +
		"tasksinqueue int, avg_completion_time real, min_fee real, " +
		"deadness int, capabilities string, last_active int, " +
		"last_24_hrs int, best_connections string"

	return s
}

func TaskWorkerSqlColumnStringValues() string {

	s := "wid, " +
		"address, port, taskdone, reputation, " +
		"tasksinqueue, avg_completion_time, min_fee, " +
		"deadness, capabilities, last_active, " +
		"last_24_hrs, best_connections"

	return s
}

func (t *TaskWorkers) Start(filePath string, create bool) error {

	if !util.FileExists(filePath) {
		create = true
	}

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

func (t *TaskWorkers) DropAllTables() error {
	sqlStmt := `
	drop table taskworkers;
	`
	_, err := t.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (t *TaskWorkers) CreateTable() error {

	t.DropAllTables()

	sqlStmt := fmt.Sprintf("create table taskworkers (%s);", TaskWorkerSqlColumnString())
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

func (t *TaskWorkers) AddTaskWorker(tw *TaskWorker) error {

	s := fmt.Sprintf("INSERT INTO taskworkers(%s) values(?,?,?,?,?,?,?,?,?,?,?,?,?)", TaskWorkerSqlColumnStringValues())

	// insert
	stmt, err := t.db.Prepare(s)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(tw.ID, tw.Address, tw.Port, tw.TasksDone, tw.Reputation,
		tw.TasksInQueue, tw.AvgCompletionTime, tw.MinimumFee,
		tw.Deadness, tw.Capabilities, tw.LastActive, tw.TasksDoneLast24Hours,
		tw.BestConnections)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(id)

	if t.OnTaskWorkersAdded != nil {
		t.OnTaskWorkersAdded(tw.ID)
	}

	return nil
}

func (t *TaskWorkers) GetTaskWorkers(filter string) ([]*TaskWorker, error) {
	// query
	rows, err := t.db.Query("SELECT * FROM taskworkers" + filter)
	if err != nil {
		return t.TaskWorkers, err
	}

	for rows.Next() {

		var task TaskWorker

		err = rows.Scan(&task.ID, &task.Address, &task.Port, &task.TasksDone, &task.Reputation, &task.TasksInQueue,
			&task.AvgCompletionTime, &task.MinimumFee, &task.Deadness, &task.Capabilities,
			&task.LastActive, &task.TasksDoneLast24Hours, &task.BestConnections)
		if err == nil {
			t.TaskWorkers = append(t.TaskWorkers, &task)
		}
	}

	rows.Close()

	return t.TaskWorkers, nil
}

func (t *TaskWorkers) GetAllTaskWorkers() ([]*TaskWorker, error) {
	return t.GetTaskWorkers("")
}
