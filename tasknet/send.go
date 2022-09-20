package tasknet

import "spinedtp/util"

// This will create a new task in our local task pool.
// Immediately after, it will propagate open tasks in
// our task pool into the next nodes. It will favour
// the latest task we just droppped of course, but it will
// also handle other open tasks it may have received
// from other nodes.
func SendTaskToNetwork(text string) {

	task := CreateNewTask(text)
	OpenTaskPool.AddTask(task)
	CheckForNewTasks()
}

func SendTaskAcceptance(task *Task, bid *TaskBid) {

	util.PrintBlue("Sending Task Acceptance for Task: " + task.ID + " (" + task.Command + ") to " + PeerIDToDescription(bid.BidderID))

	// Get the target client ID
	targetID := bid.BidderID

	// If we are connected directly to the target client, then send it directly
	// without sending to any other clients
	foundTarget := false
	for _, peer := range Peers {
		if peer.ID == targetID {
			peer.AcceptBid(task, bid)
			foundTarget = true
			break
		}
	}

	if !foundTarget {
		for _, peer := range Peers {
			peer.AcceptBid(task, bid)
		}
	}
}

// This function will need to be improved a lot. This is because the submissions can be quite large.
// Sending the result through all peers and potentially over the entire network is not going to be good
// Solutions:
//  1. The acceptance packet should contain all peers that the task-giver is connected to.
//  2. The acceptance packet naturally contains the route over which it came
//  3. The worker sends the result to the route it arrived from. It also provides the connection list
//     The next peer checks if it can reach the task giver. If not, it gives up and informs worker.
//  4. In this situation, the worker sends to all the connected peers. If no peer has any of the connected
//     ones connected to it, then they give up. This way we need maximum of two hops between worker and
//     the task giver. This is better for network fairness, so everyone has a chance to get jobs.
//  5. If the above method still congests the network, we will use 'judges' who will provide IP routes. They
//     can also do the transfer for a fee. Generally, there should be a fee for those transferring results.
func SendTaskSubmission(task *Task, submissionData *[]byte) {

	util.PrintBlue("Sending Task Submission for Task: " + task.ID + " (" + task.Command + ") ")

	// Get the target client ID
	targetID := task.TaskOwnerID

	var tr TaskResult
	tr.Data = *submissionData
	tr.MimeType = "image" // for now, we will differentiate later

	// If we are connected directly to the target client, then send it directly
	// without sending to any other clients
	foundTarget := false
	for _, peer := range Peers {
		if peer.ID == targetID {

			task.Results = append(task.Results, tr)
			peer.SubmitTaskResult(task)
			foundTarget = true
			break
		}
	}

	// Otherwise send to all connected clients
	if !foundTarget {
		for _, peer := range Peers {
			task.Results = append(task.Results, tr)
			peer.SubmitTaskResult(task)
		}
	}

}
