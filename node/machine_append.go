package node

import (
	"github.com/djreed/raft/data"
)

// TODO if we get a term > currentTerm, convert to Follower, set term to higher
func HandleAppendEntries(n *Node, appendEntries data.AppendEntries) data.MessageList {
	n.HandleTermUpdate(appendEntries.TermId, appendEntries.Leader)
	n.ResetElectionTimeout()
	n.SetLeader(appendEntries.Leader) // TODO validate in proper world

	core := n.NewMessageCoreId(appendEntries.Source, data.APPEND_RES_MSG, appendEntries.MessageId)
	termCore := n.NewTermCore()
	response := data.AppendEntriesResponse{
		MessageCore: core,
		TermCore:    termCore,
		Success:     true,
	}

	if appendEntries.TermId < n.State.CurrentTerm { // #1
		response.Success = false
		return MakeList(response)
	}

	prevLogEntry, exists := n.State.GetLogEntry(appendEntries.PrevLogIndex)
	if appendEntries.PrevLogIndex > 0 && (!exists || prevLogEntry.Term != appendEntries.PrevLogTerm) { // #2
		response.Success = false
		return MakeList(response)
	}

	newLogsIdx := appendEntries.PrevLogIndex + 1                               // Starting at previous + 1
	n.State.Log = append(n.State.Log[:newLogsIdx-1], appendEntries.Entries...) // #3 and #4

	if appendEntries.LeaderCommit > n.State.CommitIndex { // #5
		endOfLogIdx := int(n.State.LastLogIndex()) // TODO off by one?
		n.State.CommitIndex = data.ENTRY_INDEX(Min(int(appendEntries.LeaderCommit), endOfLogIdx))
	}

	return MakeList(response)
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func HandleAppendEntriesResponse(n *Node, appendRes data.AppendEntriesResponse, isCommitting bool) (messageList data.MessageList, stateChange bool) {
	// TODO Handle term mismatch
	if n.HandleTermUpdate(appendRes.TermId, appendRes.Leader) {
		stateChange = true
		return
	} // Past here, our term == appendRes.TermId

	if appendRes.Success {
		if isCommitting && !n.AlreadyReplicated(appendRes.Source) && n.IsReplicationMessage(appendRes.MessageId) {
			ERR.Printf("(!!! %v !!!) Got Replication from [%s] for '%s'\n", n.Id, appendRes.Source, appendRes.MessageId)
			n.SetReplicated(appendRes.Source)
			n.IncrementReplications()
		}

		// TODO these are almost certainly wrong lmao
		lastReplicatedIdx := n.State.MatchIndex[n.NeighborIndex(appendRes.Source)] // Which index is known to be replicated
		sendStartIdx := n.State.NextIndex[n.NeighborIndex(appendRes.Source)]       // Which index did we start sending at
		sentCount := data.ENTRY_INDEX(0)                                           // How many messages were sent
		if n.State.LastLogIndex() > lastReplicatedIdx {
			// If there are local log entries that haven't been replicated
			// to a given node, we've sent data in the append this response is for
			sentCount = n.State.LastLogIndex() - lastReplicatedIdx
		}

		n.State.MatchIndex[n.NeighborIndex(appendRes.Source)] = lastReplicatedIdx + sentCount
		n.State.NextIndex[n.NeighborIndex(appendRes.Source)] = sendStartIdx + sentCount
	} else {
		// Failure -- decrement index to search for log difference point
		n.State.NextIndex[n.NeighborIndex(appendRes.Source)]--
	}

	messageList = make(data.MessageList, 0)
	replicatedQuorum := n.ReplicationQuorum()
	if replicatedQuorum {
		// Can commit all messages on log
		n.State.CommitAll()

		// Replicated PUT(s) to quorum, can respond to client(s)
		knownReplicatedIdx := n.State.MatchIndex[n.NeighborIndex(appendRes.Source)]
		for _, replicatedPut := range n.State.Log[knownReplicatedIdx-1:] {
			core := n.NewMessageCoreId(replicatedPut.Sender, data.OK_MSG, replicatedPut.MID)
			response := data.PutResponse{MessageCore: core}
			messageList = append(messageList, response)
		}
	}

	stateChange = replicatedQuorum
	return
}

func KnownReplicatedIdx(n *Node, recv data.NODE_ID) data.ENTRY_INDEX {
	nodeIdx := n.NeighborIndex(recv)
	lastKnownReplicatedIdx := n.State.IndexReplicated(nodeIdx)
	if lastKnownReplicatedIdx > 0 {
		return lastKnownReplicatedIdx + 1 // TODO changes in batched world
	}
	return lastKnownReplicatedIdx
}

func NextIndexToSend(n *Node, recv data.NODE_ID) data.ENTRY_INDEX {
	nodeIdx := n.NeighborIndex(recv)
	lastSentIdx := n.State.IndexToSend(nodeIdx)
	if lastSentIdx > 1 {
		ERR.Printf("(%v) Incrementing index to send to [%s] from %d -> %d\n", n.Id, recv, lastSentIdx, lastSentIdx+1)
		return lastSentIdx + 1 // TODO changes in batched world
	}
	return lastSentIdx
}
