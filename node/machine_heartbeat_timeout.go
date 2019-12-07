package node

import "github.com/djreed/raft/data"

func HandleHeartbeatTimeout(n *Node) (data.MessageList, bool) {
	messages := make(data.MessageList, 0)
	for idx, nodeId := range n.Neighbors {
		msgCore := n.NewMessageCore(nodeId, data.APPEND_MSG)
		termCore := n.NewTermCore()

		toSendIdx := n.State.IndexToSend(idx) // Index of new entries
		// replicatedIdx := n.State.IndexReplicated(idx) // TODO Why?

		prevLogIdx := toSendIdx - 1 // Preceding
		prevLogEntry, present := n.State.GetLogEntry(prevLogIdx)
		prevLogTerm := data.TERM_ID(0)
		if present {
			prevLogTerm = prevLogEntry.Term
		}

		request := data.AppendEntries{
			MessageCore:  msgCore,
			TermCore:     termCore,
			Entries:      make([]data.LogEntry, 0),
			PrevLogIndex: prevLogIdx,
			PrevLogTerm:  prevLogTerm,
			LeaderCommit: n.State.CommitIndex,
		}

		// Is data being sent
		if 0 < toSendIdx && toSendIdx <= n.State.LastLogIndex() {
			// Is non-committed data being sent
			if toSendIdx > n.State.CommitIndex {
				n.AddReplicationMid(request.MessageId)
			}
			entriesToSend := []data.LogEntry{n.State.Log[toSendIdx-1]} // TODO batching
			request.Entries = entriesToSend
		}

		messages = append(messages, request)
	}

	n.ResetHeartbeatTimeout()
	return messages, n.PendingCommits()
}
