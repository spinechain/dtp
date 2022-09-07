package tasknet

import (
	"spinedtp/util"
	"time"
)

// Waits till the expiry of the bid timeout for a particular task
func WaitForBidExpiry(task *Task) {

	task.BidTimeoutTimer = time.NewTimer(NetworkSettings.BidTimeoutSeconds * time.Second)
	<-task.BidTimeoutTimer.C
	task.GlobalStatus = StatusBiddingComplete
	OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, task.LocalStatus)
	taskForProcessingAvailable <- 1
}

func SelectWinningBids(task *Task) error {

	// Check that the same person is not bidding for the same task twice
	full_query := "SELECT * FROM bids where task_id=? ORDER BY bid_value ASC"
	stmt, err := OpenTaskPool.db.Prepare(full_query)
	if err != nil {
		return err
	}

	rows, err := stmt.Query(task.ID)
	if err != nil {
		return err
	}

	var i int
	for rows.Next() {

		// bid_id, task_id, created, fee, bid_value, bidder_id, geo, arrival_route, selected
		var bid TaskBid
		var created string
		var arrival_route string
		err = rows.Scan(&bid.ID, &bid.TaskID, &created, &bid.Fee, &bid.BidValue, &bid.BidderID, &bid.Geo, &arrival_route, &bid.Selected)
		if err != nil {
			return err
		}

		bid.Created, _ = time.Parse("2006-01-02 15:04:05-07:00", created)
		if bid.BidValue > task.Reward {
			continue
		}

		if i >= int(NetworkSettings.AcceptedBidsPerTask) {
			break
		}

		for _, peer := range Peers {
			peer.AcceptBid(task, bid)
		}

		i++
	}

	return nil
}

func TaskSubmissionReceived(tt *TaskSubmission) {
	util.PrintPurple("Received task submission with value: " + string(tt.Submission))

	/*
		for _, task := range taskPool.myTasks {
			if task.ID == tt.TaskID && task.Status == BidsSelected {

				// Once we find it, we move it to our pool for tasks we are working on
				task.Status = Completed
				taskPool.completedTasks = append(taskPool.completedTasks, task)

				taskFile := filepath.Join(networkSettings.DataFolder, task.ID)

				err := os.WriteFile(taskFile, tt.Submission, 0644)
				if err != nil {
					fmt.Println()
				}
				break
			}
		}
	*/
}

// This means that we bid for a task and we have been accepted as one of
// those to execute the task
func TaskAcceptanceReceived(tt *TaskAccept) {

	/*
		util.PrintYellow("Received new task acceptance")

		// We need to find the task in our taskpool. If it's not there, we should
		// not do it

		for _, task := range taskPool.networkTasks {
			if task.ID == tt.TaskID {
				// Once we find it, we move it to our pool for tasks we are working on
				task.Status = AcceptedForWork
				taskPool.acceptedTasks = append(taskPool.acceptedTasks, task)
				taskForExecutionAvailable <- 1
				break
			}
		}
	*/

}
