package taskpool

import (
	"database/sql"
	"fmt"
	"spinedtp/tasknet"

	_ "github.com/mattn/go-sqlite3"
)

// This will maintain the db api with all the tasks in the ENTIRE network.
// Each node is a client and a server

// The header of a spine packet
type Taskpool struct {
	Type  string // done or outstanding
	db    *sql.DB
	Tasks []*tasknet.Task
}

func (t *Taskpool) Start(filePath string, create bool) error {
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

func (t *Taskpool) Stop() {
	if t.db != nil {
		t.db.Close()
	}
}

func (t *Taskpool) CreateTable() error {
	sqlStmt := `
	create table tasks (tid integer not null primary key, description text);
	delete from tasks;
	`
	_, err := t.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) AddTask(taskID string, taskString string) error {

	// insert
	stmt, err := t.db.Prepare("INSERT INTO tasks(tid, description) values(?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(taskID, taskString)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(id)

	return nil
}

func (t *Taskpool) RemoveTask(taskID string) error {

	// delete
	stmt, err := t.db.Prepare("delete from tasks where tid=?")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(taskID)
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) UpdateTask(tid string, taskdesc string) error {

	// update
	stmt, err := t.db.Prepare("update tasks set tid=? where description=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(tid, taskdesc)
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) GetAllTasks() ([]*tasknet.Task, error) {
	// query
	rows, err := t.db.Query("SELECT * FROM tasks")
	if err != nil {
		return t.Tasks, err
	}

	var tid int
	var description string

	for rows.Next() {
		err = rows.Scan(&tid, &description)
		if err == nil {
			fmt.Println(tid)
			fmt.Println(description)
		}
	}

	rows.Close()

	return t.Tasks, nil
}

func (t *Taskpool) GetTaskApprovedForClient() ([]*tasknet.Task, error) {
	return t.GetAllTasks()
}
