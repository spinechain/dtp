package tasknet

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lithammer/shortuuid"
	_ "github.com/mattn/go-sqlite3"
)

// The taskpool is where tasks are queued up that can be processed by the network. Each node can have a somewhat
// different taskpool, there is no requirement of consistency across the entire network. However, the task pool
// is synchronised as much as possible.

// Tasks are only maintained up to the space available. Tasks are ordered by fee. Tasks time-out after a fixed amoung
// of time. A task can be updated with a higher fee.

const TASKPOOL_MAX_ITEMS = 1000
const TASKPOOL_MAX_KB = 1024 * 1024
const TASKPOOL_EXPIRY_DAYS = 7

// This will maintain the db api with all the tasks in the ENTIRE network.
// Each node is a client and a server

type Taskpool struct {
	Type         string // done or outstanding
	db           *sql.DB
	Tasks        []*Task
	OnTaskAdded  func(string, string) // ID, Description
	highestIndex uint64
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

func (t *Taskpool) DropAllTables() error {
	sqlStmt := `
	drop table tasks;
	`
	_, err := t.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) CreateTable() error {

	t.DropAllTables()

	// propagated variable happens when we receive a task. We propagate it to all connected
	// clients. But if a new client connects, we don't use this mechanism to propagate tasks
	// to it, rather, it makes a request for tasks from a certain height.
	sqlStmt := `
	create table tasks (tid text not null primary key, description text, propagated int);
	delete from tasks;
	`
	_, err := t.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) AddMyTask(taskString string) error {
	tid := shortuuid.New()
	return t.AddTask(tid, taskString)
}

func (t *Taskpool) AddTaskStructure(task *Task) error {

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

	if t.OnTaskAdded != nil {
		t.OnTaskAdded(taskID, taskString)
	}

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

// Get tasks that we have not yet fully propagated
func (t *Taskpool) GetTasksForPropagation() ([]*Task, error) {
	return t.GetTasks("where propagated=0")
}

func (t *Taskpool) GetAllTasks() ([]*Task, error) {
	return t.GetTasks("")
}

func (t *Taskpool) GetTasks(filter string) ([]*Task, error) {
	// query
	rows, err := t.db.Query("SELECT * FROM tasks" + filter)
	if err != nil {
		return t.Tasks, err
	}

	var tid string
	var description string

	t.Tasks = nil

	for rows.Next() {
		err = rows.Scan(&tid, &description)
		if err == nil {
			var task Task
			task.Command = description
			task.ID = tid
			t.Tasks = append(t.Tasks, &task)
		}
	}

	rows.Close()

	return t.Tasks, nil
}

func (t *Taskpool) GetTaskApprovedForClient() ([]*Task, error) {
	return t.GetAllTasks()
}

// We take a request here and put it into the tasktransaction
// format. Then we drop it into our taskpool. It will be maintained
// there till the taskpool expires. We will also regularly propagate
// our taskpool to connected peers
func (t *Taskpool) AddToNetworkTaskPool(task *Task) {

	/*
		for _, t := range taskPool.networkTasks {
			if t.ID == task.ID {
				util.PrintYellow("Task received we already have in taskpool")
				return
			}
		}

		// Add to the taskpool list
		taskPool.networkTasks = append(taskPool.networkTasks, task)

		IncHighestIndex(task.Index)
	*/
}

func FindInNetworkTaskPool(id string) (*Task, error) {

	/*
		for _, task := range taskPool.networkTasks {
			if task.ID == id {
				return task, nil
			}
		}
	*/
	return nil, errors.New("task not found")
}

func FindInMyTaskPool(id string) (*Task, error) {

	/*
		for _, task := range taskPool.myTasks {
			if task.ID == id {
				return task, nil
			}
		}
	*/
	return nil, errors.New("task not found")
}

/*
func GetNetworkTaskList() []string {
	var tasks []string

	for _, t := range taskPool.networkTasks {
		tasks = append(tasks, t.Command)
	}

	return tasks
}
*/

/*
func GetMyTaskList() []string {
	var tasks []string

	for _, t := range taskPool.myTasks {
		tasks = append(tasks, t.Command)
	}

	return tasks
}
*/

func ReorganiseTaskPool() {

}
