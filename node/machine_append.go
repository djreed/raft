package node

import (
	"github.com/djreed/raft/data"
)

func HandleAppendEntries(n *Node, appendEntries data.AppendEntries) data.MessageList {
	n.SetLeader(appendEntries.Leader)
	n.HandleTermUpdate(appendEntries.TermId, appendEntries.Leader)
	if !n.IsLeader() {
		n.ResetElectionTimeout()
	}

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
		endOfLogIdx := int(n.State.LastLogIndex())
		n.State.CommitIndex = data.ENTRY_INDEX(Min(int(appendEntries.LeaderCommit), endOfLogIdx))
	}

	return MakeList(response)
}

func HandleAppendEntriesResponse(n *Node, appendRes data.AppendEntriesResponse, isCommitting bool) (messageList data.MessageList, stateChange bool) {
	if n.HandleTermUpdate(appendRes.TermId, data.UNKNOWN_LEADER) {
		stateChange = true
		return
	} // Past here, our term == appendRes.TermId

	if appendRes.Success {
		if isCommitting && !n.AlreadyReplicated(appendRes.Source) && n.IsReplicationMessage(appendRes.MessageId) {
			// ERR.Printf("(%v) Replication from [%s] for '%s'\n", n.Id, appendRes.Source, appendRes.MessageId)
			n.SetReplicated(appendRes.Source)
			n.IncrementReplications()
		}

		/// INDEX MANAGEMENT
		// How many entries did we send in the AppendEntries this is responding to
		sentCount := data.ENTRY_INDEX(0)
		originalMessage := n.AppendMessages[appendRes.MessageId]
		// if originalMessage != nil {
		sentCount = data.ENTRY_INDEX(len(originalMessage.Entries))
		// }

		n.DeleteAppendMessage(appendRes.MessageId)

		lastReplicatedIdx := n.LastReplicatedIdx(appendRes.Source) // Which index is known to be replicated
		sendStartIdx := n.SendStartIdx(appendRes.Source)           // Which index did we start sending at

		n.State.MatchIndex[n.NeighborIndex(appendRes.Source)] = lastReplicatedIdx + sentCount
		n.State.NextIndex[n.NeighborIndex(appendRes.Source)] = sendStartIdx + sentCount
		/// INDEX MANAGEMENT

		/// POTENTIAL RESPONSES
		if isCommitting && n.ReplicationQuorum() {
			messageList = make(data.MessageList, 0)
			// ERR.Printf("(%v) ___Quorum Reached___", n.Id)

			// Can commit all messages on log
			n.State.ApplyAll()

			// Replicated PUT(s) to quorum, can respond to client(s)
			// NOTE, everything _after_ commitIdx, hence no -1
			for _, replicatedPut := range n.State.Log[n.State.CommitIndex:] {
				core := n.NewMessageCoreId(replicatedPut.Sender, data.OK_MSG, replicatedPut.MID)
				response := data.PutResponse{MessageCore: core}
				messageList = append(messageList, response)
			}

			n.State.CommitAll()
		}
		/// POTENTIAL RESPONSES
	} else {
		// Failure -- decrement index to search for log difference point
		n.State.NextIndex[n.NeighborIndex(appendRes.Source)]--
	}

	stateChange = n.ReplicationQuorum()
	return
}
