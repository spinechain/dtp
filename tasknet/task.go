package tasknet

import (
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v3"
)

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

type GlobalTaskStatus int
type LocalTaskStatus int

// Global status on a task. To be propagated
const (
	StatusWaitingForBids GlobalTaskStatus = iota
	StatusBiddingComplete
	StatusAcceptedWorkers
	StatusWorkComplete
	StatusCompletedAndPaid
	StatusTimeoutAndDead
)

// Local status, not to be propagated, but kept in our db
const (
	StatusNewFromNetwork LocalTaskStatus = iota
	StatusNewFromLocal
	StatusRoutedToNetwork
	StatusWaitingForBidsForMe
	StatusBiddingPeriodExpired
	StatusWaitingForExecution
	StatusInExecutionPhase
	StatusTimeout
	StatusSentBid
	StatusNotGoingToBid
	StatusSubmittedResults
	StatusApprovedForMe
	StatusExecuting
	StatusNotSelectedNoPay
	StatusSuccessfullAndPaid
)

type TaskResult struct {
	Data     []byte
	MimeType string
}

type Task struct {
	ID                      string    // The globally unique ID of this task
	Command                 string    // The actual request
	Created                 time.Time // when the creator created it
	Fee                     float64   // fee for putting it in the network
	Reward                  float64   // reward for whoever solves the task
	TaskOwnerID             string    // node that created this
	Index                   uint64    // Non-reliable index that indicates roughly where this transaction is in global transaction pool
	PropagatedTo            []string  // the peers I have sent it to
	FullyPropagated         bool      // Set to true if we won't send this to any other clients
	GlobalStatus            GlobalTaskStatus
	LocalWorkerStatus       LocalTaskStatus // indicates the status of this task for us as a worker
	LocalWorkProviderStatus LocalTaskStatus // indicates the status of this task for us if we initiated this task for the network
	Bids                    []TaskBid
	BidTimeoutTimer         *time.Timer
	BidEndTime              time.Time
	ArrivalRoute            []*Peer
	Results                 []TaskResult
	TaskHash                string // to prevent changes
}

func (task *Task) GlobalStatusAsString() string {

	switch task.GlobalStatus {
	case StatusWaitingForBids:
		return "Waiting for Bids"
	case StatusBiddingComplete:
		return "Bidding Complete"
	case StatusAcceptedWorkers:
		return "Accepted Workers"
	case StatusWorkComplete:
		return "Work Complete"
	case StatusCompletedAndPaid:
		return "Completed - Paid"
	case StatusTimeoutAndDead:
		return "Timeout - Dead"
	}

	return "Unknown Status" + strconv.Itoa(int(task.GlobalStatus))
}

func (task *Task) LocalStatusAsString(status LocalTaskStatus) string {

	switch status {
	case StatusNewFromLocal:
		return "New (mine)"
	case StatusBiddingPeriodExpired:
		return "Bidding over"
	case StatusNewFromNetwork:
		return "New (remote)"
	case StatusWaitingForBidsForMe:
		return "Waiting for Bids"
	case StatusSentBid:
		return "Sent Bid"
	case StatusWaitingForExecution:
		return "Waiting for Execution"
	case StatusNotGoingToBid:
		return "Not bidding"
	case StatusApprovedForMe:
		return "Approved"
	case StatusSubmittedResults:
		return "Submitted Results"
	case StatusNotSelectedNoPay:
		return "Task Cancelled"
	case StatusSuccessfullAndPaid:
		return "Accepted - Paid"
	case StatusTimeout:
		return "Timeout"
	case StatusInExecutionPhase:
		return "In Execution Phase"
	}

	return "Unknown Status: " + strconv.Itoa(int(status))
}

func (task *Task) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}

func (tba *TaskBidApproval) GetReturnRoute() []*Peer {
	return tba.ArrivalRoute
}
func (task *TaskSubmission) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}
func (tc *TaskCompleted) GetReturnRoute() []*Peer {
	return tc.ArrivalRoute
}

func (task *TaskBid) GetReturnRoute() []*Peer {
	return task.ArrivalRoute
}

func (task *Task) MarkAsPropagated(t *Taskpool) {
	task.FullyPropagated = true

	t.UpdateTask(task)
}

func CreateNewTask(taskCmd string) *Task {

	if len(NetworkSettings.MyPeerID) < 3 {
		panic("Node Id should not be so short")
	}

	var task Task
	task.Command = taskCmd
	task.ID = shortuuid.New()
	task.Created = time.Now()
	task.BidEndTime = time.Now().AddDate(0, 0, 1)
	task.Fee = 0.00001
	task.GlobalStatus = StatusWaitingForBids
	task.LocalWorkerStatus = StatusNewFromLocal
	task.LocalWorkProviderStatus = StatusNewFromLocal
	task.Reward = 0.0001
	task.TaskOwnerID = NetworkSettings.MyPeerID
	task.FullyPropagated = false
	task.Index = OpenTaskPool.highestIndex + 1

	OpenTaskPool.IncHighestIndex(task.Index)

	return &task
}

func RemoveIndex(s []*Task, index int) []*Task {
	return append(s[:index], s[index+1:]...)
}
