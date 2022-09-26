package tasknet

import "spinedtp/util"

// This file contains the routines that forward our packages to other peers

func RouteTaskBidOn(tb *TaskBid) {
	util.PrintBlue("Routing TaskBid On: " + tb.ID + " (TaskID:" + tb.TaskID + ") ")

	// check if the target is connected directly to us
	for _, peer := range Peers {

		if peer.ID == tb.TaskOwnerID {

			packet, err := ConstructTaskBidPacket(tb, tb.GetReturnRoute())
			if err == nil {
				peer.SendPacket(packet)
				return
			}
		}
	}

	util.PrintRed("We received a bid that is not directly connected to us. Not routing yet.")
}

func RouteTaskBidApprovalOn(tb *TaskBidApproval) {
	util.PrintBlue("Routing TaskBid On: " + tb.ID + " (TaskID:" + tb.TaskID + ") ")

	// chek if the target is connected directly to us
	for _, peer := range Peers {

		if peer.ID == tb.BidderID {

			packet, err := ConstructTaskBidApprovalPacket(tb, tb.GetReturnRoute())
			if err == nil {
				peer.SendPacket(packet)
				return
			}
		}
	}

	util.PrintRed("We received a bid approval that is not directly connected to us. Not routing yet.")
}
