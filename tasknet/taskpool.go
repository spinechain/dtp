package tasknet

import (
	"database/sql"
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
	Tasks        []*Task
	OnTaskAdded  func(string, string) // ID, Description
	highestIndex uint64
}

func (t *Taskpool) Start(filePath string, create bool) error {

	return OpenDB(filePath)
}

func (t *Taskpool) Stop() {
	CloseDB()
}

func (t *Taskpool) AddTask(task *Task) error {

	// insert
	stmt, err := taskDb.Prepare("INSERT INTO tasks(tid, command, created, fee, reward, owner_id, height, propagated, local_worker_status, local_work_provider_status, global_status, bid_timeout, task_hash) values(?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(task.ID, task.Command, task.Created, task.Fee,
		task.Reward, task.TaskOwnerID, task.Index,
		task.FullyPropagated, task.LocalWorkerStatus,
		task.LocalWorkProviderStatus, task.GlobalStatus,
		task.BidEndTime, task.TaskHash)

	if err != nil {
		return err
	}

	if t.OnTaskAdded != nil {
		t.OnTaskAdded(task.ID, task.Command)
	}

	return nil
}

func (t *Taskpool) RemoveAllTasks() error {

	stmt, err := taskDb.Prepare("delete from tasks")
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

	stmt, err = taskDb.Prepare("delete from bids")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = taskDb.Prepare("delete from bids_sent")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (t *Taskpool) RemoveTask(taskID string) error {

	// delete
	stmt, err := taskDb.Prepare("delete from tasks where tid=?")
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

func (t *Taskpool) UpdateTaskStatus(task *Task, newGlobalStatus GlobalTaskStatus, newLocalWorkerStatus LocalTaskStatus, newLocalWorkProviderStatus LocalTaskStatus) error {

	task.LocalWorkerStatus = newLocalWorkerStatus
	task.LocalWorkProviderStatus = newLocalWorkProviderStatus
	task.GlobalStatus = newGlobalStatus

	// update
	stmt, err := taskDb.Prepare("update tasks set local_worker_status=?, local_work_provider_status=?, global_status=? where tid=?")
	if err != nil {
		util.PrintRed("UpdateTaskStatus: " + err.Error())
		return err
	}

	_, err = stmt.Exec(task.LocalWorkerStatus, task.LocalWorkProviderStatus, task.GlobalStatus, task.ID)
	if err != nil {
		util.PrintRed("UpdateTaskStatus Exec: " + err.Error())
		return err
	}

	if NetworkCallbacks.OnTaskStatusUpdate != nil {
		NetworkCallbacks.OnTaskStatusUpdate()
	}

	return nil
}

func (t *Taskpool) UpdateTask(task *Task) error {

	// update
	stmt, err := taskDb.Prepare("update tasks set command=?, propagated=? where tid=?")
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

func (t *Taskpool) GetTask(task_id string) *Task {
	task, err := t.GetTasks("where tid=?", task_id)
	if err != nil || len(task) < 1 {
		return nil
	}

	return task[0]
}

func (t *Taskpool) GetTasks(filter string, args ...any) ([]*Task, error) {
	// query

	full_query := "SELECT * FROM tasks " + filter

	stmt, err := taskDb.Prepare(full_query)
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

	defer rows.Close()

	for rows.Next() {

		var task Task

		var created string
		var bid_end_time string
		err = rows.Scan(&task.ID, &task.Command, &created,
			&task.Fee, &task.Reward, &task.TaskOwnerID,
			&task.Index, &task.FullyPropagated, &task.LocalWorkerStatus, &task.LocalWorkProviderStatus, &task.GlobalStatus,
			&bid_end_time, &task.TaskHash)
		if err == nil {

			// date, error := time.Parse("2006-01-02", dateString)
			task.BidEndTime, _ = time.Parse("2006-01-02 15:04:05.0000000-07:00", bid_end_time)
			task.Created, _ = time.Parse("2006-01-02 15:04:05.0000000-07:00", created)
			t.Tasks = append(t.Tasks, &task)
		}
	}

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
	fmt.Print("Existing task: " + existingTask.ID + " " + existingTask.Command)

	OpenTaskPool.IncHighestIndex(task.Index)

	// Review what's going on down here

	// We are
	// We update the local statii of the existing task in the db
	// TODO: this looks risky, as we are taking this info from network
	// existingTask := tasks[0]

	// Bug here. Tasks exists and is routed, but we change local status to new from network
	// t.UpdateTaskStatus(task, existingTask.GlobalStatus, task.LocalWorkerStatus, task.LocalWorkProviderStatus)

	// We have the task. We need to update the status
	//if existingTask.LocalStatus == StatusNewFromLocal && task.LocalStatus == StatusNewFromNetwork {
	//	t.UpdateTaskStatus(task, task.GlobalStatus, task.LocalStatus)
	//}
}

func ReorganiseTaskPool() {

}
