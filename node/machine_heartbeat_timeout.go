package node

import "github.com/djreed/raft/data"

func HandleHeartbeatTimeout(n *Node) (data.MessageList, bool) {
	pendingCommit := n.PendingCommits()

	messages := make(data.MessageList, 0)
	for idx, nodeId := range n.Neighbors {
		msgCore := n.NewMessageCore(nodeId, data.APPEND_MSG)
		termCore := n.NewTermCore()

		toSendIdx := n.State.IndexToSend(idx) // Index of new entries
		// replicatedIdx := n.State.IndexReplicated(idx) // TODO Why?

		prevLogIdx := toSendIdx - 1 // Preceding
		entry, present := n.State.GetLogEntry(prevLogIdx)
		prevLogTerm := n.State.CurrentTerm
		if present {
			prevLogTerm = entry.Term
		}
		leaderCommit := n.State.CommitIndex

		request := data.AppendEntries{
			MessageCore:  msgCore,
			TermCore:     termCore,
			Entries:      make([]data.LogEntry, 0),
			PrevLogIndex: prevLogIdx,
			PrevLogTerm:  prevLogTerm,
			LeaderCommit: leaderCommit,
		}

		if 0 <= toSendIdx && toSendIdx <= n.State.LastLogIndex() {
			entriesToSend := []data.LogEntry{n.State.Log[toSendIdx]} // TODO batching
			request.Entries = entriesToSend
		}

		messages = append(messages, request)
	}

	n.ResetHeartbeatTimeout()
	return messages, pendingCommit
}
