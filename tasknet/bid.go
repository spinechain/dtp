package tasknet

import (
	"database/sql"
	"errors"
	"spinedtp/util"
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v3"
)

func (bid *TaskBid) Scan(rows *sql.Rows) error {

	var created string
	var arrival_route string
	err := rows.Scan(&bid.ID, &bid.TaskID, &created, &bid.Fee, &bid.BidValue, &bid.BidderID, &bid.Geo, &arrival_route, &bid.Selected, &bid.MyBid)
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	bid.Created, _ = time.Parse("2006-01-02 15:04:05-07:00", created)

	return err
}

func (bid *TaskBid) MarkAsAccepted() error {

	bid.Selected = 1
	return bid.UpdateBid(bid)
}

func (bid *TaskBid) UpdateBid(b *TaskBid) error {

	// update
	stmt, err := taskDb.Prepare("update bids set selected=? where bid_id=?")
	if err != nil {
		util.PrintRed("UpdatedBid: " + err.Error())
		return err
	}

	_, err = stmt.Exec(b.Selected, bid.ID)
	if err != nil {
		util.PrintRed("UpdatedBid: " + err.Error())
		return err
	}

	return nil
}

// Waits till the expiry of the bid timeout for a particular task
func WaitForBidExpiry(task *Task) {

	task.BidTimeoutTimer = time.NewTimer(NetworkSettings.BidTimeoutSeconds * time.Second)
	<-task.BidTimeoutTimer.C
	task.GlobalStatus = StatusBiddingComplete
	OpenTaskPool.UpdateTaskStatus(task, StatusBiddingComplete, task.LocalWorkerStatus, StatusBiddingPeriodExpired)
	taskForProcessingAvailable <- 1
}

func CreateTaskBid(task *Task) *TaskBid {
	var t TaskBid
	t.BidValue = task.Reward - 0.0001
	t.BidderID = NetworkSettings.MyPeerID
	t.Fee = 0
	t.Geo = "US"
	t.ID = shortuuid.New()
	t.TaskOwnerID = task.TaskOwnerID
	t.TaskID = task.ID
	t.Created = time.Now()
	t.MyBid = 17

	return &t
}

func AddBid(db *sql.DB, bid *TaskBid, isMyBid bool) error {

	task := OpenTaskPool.GetTask(bid.TaskID)
	if task == nil {
		return errors.New("task not found")
	}

	if isMyBid {
		if bid.MyBid != 17 {
			panic("You must set MyBid to 17 if this is your own bid. Security reasons.")
		}
	}

	// Check that the same person is not bidding for the same task twice
	full_query := "SELECT count(*) FROM bids where bidder_id=? and task_id=?"
	stmt, err := db.Prepare(full_query)
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	rows, err := stmt.Query(bid.BidderID, task.ID)
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	defer rows.Close()

	for rows.Next() {

		var item_count string
		err = rows.Scan(&item_count)
		if err != nil {
			return err
		}

		cnt, _ := strconv.Atoi(item_count)
		if cnt != 0 {
			return errors.New("bid exists")
		}
	}

	// insert the bid to db
	stmt, err = db.Prepare("INSERT INTO bids(bid_id, task_id, created, fee, bid_value, bidder_id, geo, arrival_route, selected, my_bid) values(?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	var arrivalRoute string
	for i := 0; i < len(bid.ArrivalRoute); i++ {
		arrivalRoute = bid.ArrivalRoute[i].ID + ";" + arrivalRoute
	}

	_, err = stmt.Exec(bid.ID, task.ID, bid.Created, bid.Fee,
		bid.BidValue, bid.BidderID, bid.Geo,
		arrivalRoute, 0, bid.MyBid)

	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	return nil
}

func GetBids(filter string, args ...any) ([]*TaskBid, error) {
	// query

	full_query := "SELECT * FROM bids " + filter

	stmt, err := taskDb.Prepare(full_query)
	if err != nil {
		util.PrintRed(err.Error())
		return nil, err
	}

	var rows *sql.Rows
	if args != nil {
		rows, err = stmt.Query(args...)
	} else {
		rows, err = stmt.Query()
	}

	defer rows.Close()

	if err != nil {
		util.PrintRed(err.Error())
		return nil, err
	}

	var bids []*TaskBid

	for rows.Next() {

		var bid TaskBid
		bid.Scan(rows)

		bids = append(bids, &bid)
	}

	return bids, nil
}

func SelectWinningBids(task *Task) error {

	// Check that the same person is not bidding for the same task twice
	full_query := "SELECT * FROM bids where task_id=? ORDER BY bid_value ASC"
	stmt, err := taskDb.Prepare(full_query)
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	rows, err := stmt.Query(task.ID)
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	defer rows.Close()

	var i int
	for rows.Next() {

		var bid TaskBid
		bid.Scan(rows)

		if i >= int(NetworkSettings.AcceptedBidsPerTask) {
			break
		}

		for _, peer := range Peers {
			peer.AcceptBid(task, &bid)
		}

		i++
	}

	return nil
}

func TaskSubmissionReceived(tt *TaskSubmission) {

	if tt.ID != NetworkSettings.MyPeerID {
		util.PrintYellow("Task submission for another client received.")
		// TODO: Route this on
		return
	}

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

	if tt.ID != NetworkSettings.MyPeerID {
		util.PrintYellow("Task acceptance for another client received.")
		// TODO: Route this on
		return
	}

	util.PrintYellow("Received new task acceptance")

	// We need to find the task in our taskpool. If it's not there, we should
	// not do it
	task := OpenTaskPool.GetTask(tt.TaskID)
	if task == nil {
		util.PrintRed("Invalid task found!")
	}

	// Let's check if we bid on it
	ourBid, err := GetBids("where bidder_id=? and task_id=?", NetworkSettings.MyPeerID, tt.TaskID)
	if err == nil && len(ourBid) >= 0 {
		// In this case, we really did bid for this
		// TODO: check that if someone sends us a bid telling us that it is our ID, that we do not
		// accept it
		OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, StatusApprovedForMe, task.LocalWorkProviderStatus)
		ourBid[0].MarkAsAccepted()
		taskForExecutionAvailable <- 1

	}
}
