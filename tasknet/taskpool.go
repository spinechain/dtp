package tasknet

import (
	"database/sql"
	"errors"
	"fmt"
	"spinedtp/util"
	"time"

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

	if !util.FileExists(filePath) {
		create = true
	}

	var err error
	t.db, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}

	if create {
		err := t.CreateTable()
		if err != nil {
			panic("Could not create tasks table!")
		}
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
	create table tasks (tid text not null unique primary key, command text, 
						created int, fee real, reward real, owner_id string, 
						height int, propagated int, local_status int, global_status int, bid_timeout int,
						task_hash string);
	delete from tasks;
	`
	_, err := t.db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) AddTask(task *Task) error {

	// insert
	stmt, err := t.db.Prepare("INSERT INTO tasks(tid, command, created, fee, reward, owner_id, height, propagated, local_status, global_status, bid_timeout, task_hash) values(?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(task.ID, task.Command, task.Created, task.Fee,
		task.Reward, task.TaskOwnerID, task.Index,
		task.FullyPropagated, task.LocalStatus, task.GlobalStatus,
		task.BidEndTime, task.TaskHash)

	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println(id)

	if t.OnTaskAdded != nil {
		t.OnTaskAdded(task.ID, task.Command)
	}

	return nil
}

func (t *Taskpool) RemoveAllTasks() error {

	// delete
	stmt, err := t.db.Prepare("delete from tasks")
	if err != nil {
		return err
	}

	res, err := stmt.Exec()
	if err != nil {
		return err
	}

	_, err = res.RowsAffected()
	if err != nil {
		return err
	}

	t.Tasks = nil

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

func (t *Taskpool) UpdateTaskStatus(task *Task, newGlobalStatus GlobalTaskStatus, newLocalStatus LocalTaskStatus) error {

	task.LocalStatus = newLocalStatus
	task.GlobalStatus = newGlobalStatus
	// update
	stmt, err := t.db.Prepare("update tasks set local_status=?, global_status=? where tid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(task.LocalStatus, task.GlobalStatus, task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) UpdateTask(task *Task) error {

	// update
	stmt, err := t.db.Prepare("update tasks set command=?, propagated=? where tid=?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(task.Command, task.FullyPropagated, task.ID)
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

func (t *Taskpool) GetTasks(filter string, args ...any) ([]*Task, error) {
	// query

	full_query := "SELECT * FROM tasks " + filter

	stmt, err := t.db.Prepare(full_query)
	if err != nil {
		return t.Tasks, err
	}

	var rows *sql.Rows
	if args != nil {
		rows, err = stmt.Query(args...)
	} else {
		rows, err = stmt.Query()
	}

	if err != nil {
		return t.Tasks, err
	}

	t.Tasks = nil

	for rows.Next() {

		var task Task

		var created string
		var bid_end_time string
		err = rows.Scan(&task.ID, &task.Command, &created,
			&task.Fee, &task.Reward, &task.TaskOwnerID,
			&task.Index, &task.FullyPropagated, &task.LocalStatus, &task.GlobalStatus,
			&bid_end_time, &task.TaskHash)
		if err == nil {

			// date, error := time.Parse("2006-01-02", dateString)
			task.BidEndTime, _ = time.Parse("2006-01-02 15:04:05.0000000-07:00", bid_end_time)
			task.Created, _ = time.Parse("2006-01-02 15:04:05.0000000-07:00", created)
			t.Tasks = append(t.Tasks, &task)
		}
	}

	rows.Close()

	return t.Tasks, nil
}

func (t *Taskpool) IncHighestIndex(newVal uint64) {
	if newVal > t.highestIndex {
		t.highestIndex = newVal
	}
}

func (t *Taskpool) GetTaskApprovedForClient() ([]*Task, error) {
	return t.GetAllTasks()
}

// We take a request here and put it into the tasktransaction
// format. Then we drop it into our taskpool. It will be maintained
// there till the taskpool expires. We will also regularly propagate
// our taskpool to connected peers
func (t *Taskpool) AddToTaskPool(task *Task) {

	// We check if we have it
	tasks, err := t.GetTasks("where tid=?", task.ID)

	if tasks == nil || err != nil || len(tasks) == 0 {
		// we do not have this task in our db. We can add it directly
		t.AddTask(task)
		return
	}

	existingTask := tasks[0]
	OpenTaskPool.IncHighestIndex(task.Index)

	// We have the task. We need to update the status
	if existingTask.LocalStatus == StatusNewFromLocal && task.LocalStatus == StatusNewFromNetwork {
		t.UpdateTaskStatus(task, task.GlobalStatus, task.LocalStatus)
	}

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
