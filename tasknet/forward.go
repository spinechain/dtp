package tasknet

import "spinedtp/util"

// This file contains the routines that forward our packages to other peers

func RouteTaskOn(task *Task) {

	util.PrintBlue("Routing Task On: " + task.ID + " (" + task.Command + ") ")

	// update status so we never deal with this task again
	OpenTaskPool.UpdateTaskStatus(task, task.GlobalStatus, StatusRoutedToNetwork, task.LocalWorkProviderStatus)
	task.MarkAsPropagated(OpenTaskPool)

	// We send to clients, except clients that were already on route or the task owner
	for _, peer := range Peers {

		// Check if this peer is in the arrival route
		alreadyOnRoute := false
		for _, routePeer := range task.ArrivalRoute {
			if routePeer.ID == peer.ID {
				alreadyOnRoute = true
				break
			}
		}

		if peer.ID != task.TaskOwnerID && !alreadyOnRoute && peer.Connected {

			packet, err := ConstructTaskPropagationPacket(task)
			if err != nil {
				continue
			}

			peer.SendPacket(packet)
		}
	}
}

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

func RouteTaskSubmissionOn(ts *TaskSubmission) {
	util.PrintBlue("Routing TaskSubmission On: " + ts.ID + " (TaskID:" + ts.TaskID + ") ")

	// check if the target is connected directly to us
	for _, peer := range Peers {

		if peer.ID == ts.TaskOwnerID {

			packet, err := ConstructTaskSubmissionPacket(ts, ts.GetReturnRoute())
			if err == nil {
				peer.SendPacket(packet)
				return
			}
		}
	}

	util.PrintRed("We received a task submission that is not directly connected to us. Not routing yet.")
}
