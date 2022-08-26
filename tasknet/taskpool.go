package tasknet

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	util "spinedtp/util"
	"time"

	"github.com/lithammer/shortuuid/v3"
)

// The taskpool is where tasks are queued up that can be processed by the network. Each node can have a somewhat
// different taskpool, there is no requirement of consistency across the entire network. However, the task pool
// is synchronised as much as possible.

// Tasks are only maintained up to the space available. Tasks are ordered by fee. Tasks time-out after a fixed amoung
// of time. A task can be updated with a higher fee.

/*
From an online guide:
    First, a user initiates a transaction from a Dapp or Wallet, such as sending funds to another account or contract
    Then the user signs that transaction with their Wallet
    The Wallet sends the signed transaction to a node, often called a gateway node, to get onto the Ethereum network (think Infura or Pocket)
    That node will verify the transaction is valid and add it to its mempool
    Since the node is connected to a group of peers, it broadcasts the transaction to other nodes.
    These peer nodes will receive the transaction, validate it, move it into their own mempool, and broadcast to additional peers, essentially replicating the transaction across the network
    Miners, as a specific kind of node, also receive the transaction from peers, validate it, and attempt to add it to a block
    Eventually, a successful miner adds a block with the transaction to the chain
    The new block is broadcast over the network
    As all nodes receive the new block from their peers, they see the included transaction and remove it from their mempool

Also look at:
	https://blog.kaiko.com/an-in-depth-guide-into-how-the-mempool-works-c758b781c608
*/

type TaskPool struct {
	networkTasks   []*Task // tasks that come from the network
	myTasks        []*Task // tasks I created and sent to network
	acceptedTasks  []*Task // tasks I have been authorized to work on
	completedTasks []*Task
	highestIndex   uint64
}

const TASKPOOL_MAX_ITEMS = 1000
const TASKPOOL_MAX_KB = 1024 * 1024
const TASKPOOL_EXPIRY_DAYS = 7

var taskPool TaskPool

var taskForProcessingAvailable chan int
var taskForExecutionAvailable chan int
var shutDownTaskThread bool = false

func (task *Task) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}

func (task *TaskApproval) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}
func (task *TaskSubmission) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}
func (task *TaskAccept) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}

func (task *TaskBid) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}

// We take a request here and put it into the tasktransaction
// format. Then we drop it into our taskpool. It will be maintained
// there till the taskpool expires. We will also regularly propagate
// our taskpool to connected peers
func AddToNetworkTaskPool(task *Task) {

	for _, t := range taskPool.networkTasks {
		if t.ID == task.ID {
			util.PrintYellow("Task received we already have in taskpool")
			return
		}
	}

	// Add to the taskpool list
	taskPool.networkTasks = append(taskPool.networkTasks, task)

	IncHighestIndex(task.Index)
}

func CreateNewTask(taskCmd string) {

	var task Task
	task.Command = taskCmd
	task.ID = shortuuid.New()
	task.Created = time.Now()
	task.Fee = 0.00001
	task.Reward = 0.0001
	task.TaskOwnerID = networkSettings.MyPeerID
	task.FullyPropagated = false
	task.Index = taskPool.highestIndex + 1

	// Add to the taskpool list with my own tasks
	taskPool.myTasks = append(taskPool.myTasks, &task)

	IncHighestIndex(task.Index)

}

func NewTaskBidArrived(tb TaskBid) {

	if tb.TaskOwnerID == networkSettings.MyPeerID {
		// This is a bid for a task of mine

		task, err := FindInMyTaskPool(tb.TaskID)

		if err == nil {
			task.Bids = append(task.Bids, tb)
		}

	} else {
		// This is a bid for another peer that is not me. We route
		// it to the best connection we have

		RoutePacketOn()
	}
}

func IncHighestIndex(newVal uint64) {
	if newVal > taskPool.highestIndex {
		taskPool.highestIndex = newVal
	}
}

func FindInNetworkTaskPool(id string) (*Task, error) {

	for _, task := range taskPool.networkTasks {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, errors.New("task not found")
}

func FindInMyTaskPool(id string) (*Task, error) {

	for _, task := range taskPool.myTasks {
		if task.ID == id {
			return task, nil
		}
	}

	return nil, errors.New("task not found")
}

func RemoveIndex(s []*Task, index int) []*Task {
	return append(s[:index], s[index+1:]...)
}

// The input in here does not know strictly about who triggered it. So it just needs
// to search for what was requested so it can respond
func EngineHasCompletedTask(taskType string, taskCommand string, taskData string) {

	// var t network.TaskSubmission
	// network.SubmitTaskResult(&t)
	if taskType == "download" {
		for _, task := range taskPool.networkTasks {
			if task.Status == AcceptedForWork {
				if task.Command == taskType+" "+taskCommand {
					// At this point, we have completed a task here internally. We are to
					// propagate it back to the person who requested this task.

					util.PrintPurple("Changing task status to WorkComplete")

					// Set the task result in the task structure here
					task.Result = []byte(taskData)
					task.Status = WorkComplete

					// Trigger the thread to return the result to network
					taskForProcessingAvailable <- 1

				}
			}
		}
	}
}

func SendNewTaskToPeers() {

	for _, task := range taskPool.myTasks {

		if !task.FullyPropagated {
			packet, err := ConstructTaskPropagationPacket(task)
			if err != nil {
				continue
			}

			task.Status = WaitingForBids
			// Todo: probably needs to be removed
			// taskPool.tasks = RemoveIndex(taskPool.tasks, i)
			SendPacketToAllPeers(packet)
			task.FullyPropagated = true

			// Set a timeout
			go WaitForBidExpiry(task)

			break
		}
	}
}

// Waits till the expiry of the bid timeout for a particular task
func WaitForBidExpiry(task *Task) {

	task.BidTimeout = time.NewTimer(networkSettings.BidTimeoutSeconds * time.Second)
	<-task.BidTimeout.C
	task.Status = BiddingComplete
	taskForProcessingAvailable <- 1
}

// This function looks through all tasks, and based on their status, it decides what
// to do with them
func ProcessAvailableTasks() {

	for _, task := range taskPool.myTasks {
		switch task.Status {
		case BiddingComplete:
			util.PrintPurple("Found a task with bidding period complete")
			SelectWinningBids(task)
			task.Status = BidsSelected
		}
	}

	for _, task := range taskPool.networkTasks {
		switch task.Status {
		case Received:
			util.PrintYellow("Found a new unprocessed task: " + task.Command)

			task.Status = Bid

			// We now bid for this task. We need to find the route that this
			// packet came in through to respond through the same one
			// TODO: this is likely not needed here
			for _, routePeer := range task.ArrivalRoute {

				peer := FindPeer(routePeer.ID)

				peer.BidForTask(task)
				break
			}

		case WorkComplete:
			util.PrintYellow("Found a completed task. Submitting: " + task.Command)
			for _, routePeer := range task.ArrivalRoute {

				peer := FindPeer(routePeer.ID)

				util.PrintYellow("Submitting task to " + peer.ID)
				peer.SubmitTaskResult(task)

				break
			}
		}
	}

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
}

// This means that we bid for a task and we have been accepted as one of
// those to execute the task
func TaskAcceptanceReceived(tt *TaskAccept) {
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

}

func ShutDownTaskPool() {

	shutDownTaskThread = true
	taskForProcessingAvailable <- 1
}

func ProcessAcceptedTasks() {

	for {

		<-taskForExecutionAvailable

		for _, task := range taskPool.acceptedTasks {

			var tt TaskSubmission
			tt.ID = shortuuid.New()
			tt.Created = time.Now()

			if task.Status == AcceptedForWork {
				util.PrintYellow("Executing Task: " + task.Command)
				tt.Submission = []byte("This would be my submission")

				if networkSettings.TaskReadyForProcessing != nil {
					networkSettings.TaskReadyForProcessing(task.Command)
				} else {
					fmt.Println("No callback available for task processing")
				}
			}

		}

		if shutDownTaskThread {
			util.PrintYellow("Task execution Thread Shutdown")
			return
		}
	}

}

// This thread waits for new tasks to come into the network
func ProcessTasks() {

	// Create new channel to wait for tasks
	taskForProcessingAvailable = make(chan int)
	taskForExecutionAvailable = make(chan int)

	// Start the thread that will do the actual work (execution of each task)
	go ProcessAcceptedTasks()

	for {

		<-taskForProcessingAvailable
		ProcessAvailableTasks()

		if shutDownTaskThread {
			fmt.Println("Task processing Thread Shutdown")
			return
		}
	}

}

func ReorganiseTaskPool() {

}

func GetNetworkTaskList() []string {
	var tasks []string

	for _, t := range taskPool.networkTasks {
		tasks = append(tasks, t.Command)
	}

	return tasks
}

func GetMyTaskList() []string {
	var tasks []string

	for _, t := range taskPool.myTasks {
		tasks = append(tasks, t.Command)
	}

	return tasks
}
