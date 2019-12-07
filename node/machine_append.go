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
		Success:     false,
	}

	if appendEntries.TermId < n.State.CurrentTerm { // #1
		return MakeList(response)
	}

	if 0 <= appendEntries.PrevLogIndex { // Valid index lower bound
		val, exists := n.State.GetLogEntry(appendEntries.PrevLogIndex)
		if !exists || val.Term != appendEntries.PrevLogTerm { // #2
			return MakeList(response)
		}

		startingIdx := appendEntries.PrevLogIndex + 1
		n.State.Log = append(n.State.Log[:startingIdx], appendEntries.Entries...) // #3 and #4

		if appendEntries.LeaderCommit > n.State.CommitIndex { // #5
			lastIdx := int(n.State.LastLogIndex()) // TODO off by one?
			n.State.CommitIndex = data.ENTRY_INDEX(Min(int(appendEntries.LeaderCommit), lastIdx))
		}

	}

	response.Success = true

	return MakeList(response)
}

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func HandleAppendEntriesResponse(n *Node, appendRes data.AppendEntriesResponse, isCommitting bool) (data.MessageList, bool) {
	// TODO Handle term mismatch
	if n.HandleTermUpdate(appendRes.TermId, appendRes.Leader) {
		return nil, true
	}

	sendStart := n.State.MatchIndex[n.NeighborIndex(appendRes.Source)]

	if appendRes.Success {
		if isCommitting && !n.AlreadyReplicated(appendRes.Source) {
			n.SetReplicated(appendRes.Source)
			n.IncrementReplications()
		}

		// TODO these are almost certainly wrong lmao
		n.State.MatchIndex[n.NeighborIndex(appendRes.Source)] = CalculateSentIndex(n, appendRes.Source)
		n.State.NextIndex[n.NeighborIndex(appendRes.Source)] = CalculateIndexToSend(n, appendRes.Source)
	} else {
		// Failure -- decrement index
		n.State.NextIndex[n.NeighborIndex(appendRes.Source)]--
	}

	// TODO TODO TODO RESPOND TO PUTS IF QUORUM
	messages := make(data.MessageList, 0)
	if n.ReplicationQuorum() {
		for _, replicatedPut := range n.State.Log[sendStart:] {
			core := n.NewMessageCoreId(replicatedPut.Sender, data.OK_MSG, replicatedPut.MID)
			response := data.PutResponse{MessageCore: core}
			messages = append(messages, response)
		}
	}

	return messages, n.ReplicationQuorum()
}

func CalculateSentIndex(n *Node, recv data.NODE_ID) data.ENTRY_INDEX {
	nodeIdx := n.NeighborIndex(recv)
	return n.State.IndexReplicated(nodeIdx) + 1
}

func CalculateIndexToSend(n *Node, recv data.NODE_ID) data.ENTRY_INDEX {
	nodeIdx := n.NeighborIndex(recv)
	return n.State.IndexToSend(nodeIdx) + 1
}
