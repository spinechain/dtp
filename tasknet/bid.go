package tasknet

import (
	"fmt"
	"sort"
	"spinedtp/util"
	"time"
)

// Waits till the expiry of the bid timeout for a particular task
func WaitForBidExpiry(task *Task) {

	task.BidTimeoutTimer = time.NewTimer(networkSettings.BidTimeoutSeconds * time.Second)
	<-task.BidTimeoutTimer.C
	task.Status = BiddingComplete
	taskForProcessingAvailable <- 1
}

func SelectWinningBids(task *Task) {

	fmt.Println("Found a task with bidding complete: " + task.Command)

	sort.SliceStable(task.Bids, func(i, j int) bool {
		return task.Bids[i].BidValue < task.Bids[j].BidValue
	})

	for i, bid := range task.Bids {

		if bid.BidValue > task.Reward {
			continue
		}

		for _, routePeer := range bid.GetReturnRoute() {

			peer := FindPeer(routePeer.ID)

			peer.AcceptBid(task, bid)
			break
		}

		if i >= int(networkSettings.AcceptedBidsPerTask) {
			break
		}

	}
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
