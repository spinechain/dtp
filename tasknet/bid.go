package tasknet

import (
	"database/sql"
	"errors"
	"spinedtp/util"
	"strconv"
	"time"

	"github.com/lithammer/shortuuid/v3"
)

// /// BID TIMEOUT
// Waits till the expiry of the bid timeout for a particular task
func WaitForBidExpiry(task *Task) {

	task.BidTimeoutTimer = time.NewTimer(NetworkSettings.BidTimeout * time.Second)
	<-task.BidTimeoutTimer.C
	task.GlobalStatus = StatusBiddingComplete
	OpenTaskPool.UpdateTaskStatus(task, StatusBiddingComplete, task.LocalWorkerStatus, StatusBiddingPeriodExpired)
	taskForProcessingAvailable <- 1

	util.PrintWhite("Bid timeout for task " + task.ID + " expired")
}

func (bid *TaskBid) ScanMyBid(rows *sql.Rows) error {

	var created string
	err := rows.Scan(&bid.ID, &bid.TaskID, &created, &bid.Fee, &bid.BidValue, &bid.Geo, &bid.Selected)
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	bid.Created, _ = time.Parse("2006-01-02 15:04:05-07:00", created)

	return err
}

func (bid *TaskBid) Scan(rows *sql.Rows) error {

	var created string
	var arrival_route string
	err := rows.Scan(&bid.ID, &bid.TaskID, &created, &bid.Fee, &bid.BidValue, &bid.BidderID, &bid.Geo, &arrival_route, &bid.Selected)
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
	stmt, err := taskDb.Prepare("update bids_sent set selected=? where bid_id=?")
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

func CreateBidForTask(db *sql.DB, bid *TaskBid) error {
	task := OpenTaskPool.GetTask(bid.TaskID)
	if task == nil {
		return errors.New("task not found")
	}

	// insert the bid to db
	stmt, err := db.Prepare("INSERT INTO bids_sent(bid_id, task_id, created, fee, bid_value, geo, selected) values(?,?,?,?,?,?,?)")
	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	var arrivalRoute string
	for i := 0; i < len(bid.ArrivalRoute); i++ {
		arrivalRoute = bid.ArrivalRoute[i].ID + ";" + arrivalRoute
	}

	_, err = stmt.Exec(bid.ID, task.ID, bid.Created, bid.Fee,
		bid.BidValue, bid.Geo, 0)

	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	return err
}

func ProcessBidForMyTask(db *sql.DB, bid *TaskBid) error {

	task := OpenTaskPool.GetTask(bid.TaskID)
	if task == nil {
		return errors.New("task not found")
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
			util.PrintRed("Double bid for task found")
			return errors.New("bid exists")
		}
	}

	// insert the bid to db
	stmt, err = db.Prepare("INSERT INTO bids(bid_id, task_id, created, fee, bid_value, bidder_id, geo, arrival_route, selected) values(?,?,?,?,?,?,?,?,?)")
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
		arrivalRoute, 0)

	if err != nil {
		util.PrintRed(err.Error())
		return err
	}

	return nil
}

func GetMyBids(filter string, args ...any) ([]*TaskBid, error) {
	// query

	full_query := "SELECT * FROM bids_sent " + filter

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
		err = bid.ScanMyBid(rows)
		if err != nil {
			util.PrintRed(err.Error())
			return nil, err
		}

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

		SendTaskBidApproved(task, &bid)

		i++
	}

	return nil
}
