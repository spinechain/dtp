package tasknet

import (
	"database/sql"
	"errors"
	"spinedtp/util"
)

var taskDb *sql.DB

func OpenDB(filePath string) error {

	if taskDb != nil {
		return errors.New("database already open")
	}
	var create bool = false
	if !util.FileExists(filePath) {
		create = true
	}

	var err error
	taskDb, err = sql.Open("sqlite3", filePath)
	if err != nil {
		return err
	}

	if create {
		err := CreateTables()
		if err != nil {
			panic("Could not create tables!")
		}
	}
	return err
}

func CloseDB() {
	if taskDb != nil {
		util.PrintPurple("Closing database")
		taskDb.Close()
	}
}

func DropAllTables() error {
	sqlStmt := `
	drop table tasks;
	`
	_, err := taskDb.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
	drop table bids;
	`
	_, err = taskDb.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func CreateTables() error {

	DropAllTables()

	// propagated variable happens when we receive a task. We propagate it to all connected
	// clients. But if a new client connects, we don't use this mechanism to propagate tasks
	// to it, rather, it makes a request for tasks from a certain height.
	sqlStmt := `
	create table tasks (tid text not null unique primary key, command text, 
						created int, fee real, reward real, owner_id string, 
						height int, propagated int, local_worker_status int, local_work_provider_status int, 
						global_status int, bid_timeout int,
						task_hash string);
	delete from tasks;
	`
	_, err := taskDb.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `
	create table bids (bid_id text not null unique primary key, task_id text, 
						created int, fee real, bid_value real, bidder_id string, 
						geo string, arrival_route string, selected int);
	delete from bids;
	`
	_, err = taskDb.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}
