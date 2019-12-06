package node

import "github.com/djreed/raft/data"

func HandleHeartbeatTimeout(n *Node) data.MessageList {
	messages := make(data.MessageList, 0)
	for idx, nodeId := range n.Neighbors {
		msgCore := n.NewMessageCore(nodeId, data.APPEND_MSG)
		termCore := n.NewTermCore()

		toSendIdx := n.State.IndexToSend(idx) // Index of new entries
		// replicatedIdx := n.State.IndexReplicated(idx) // TODO Why?

		prevLogIdx := toSendIdx - 1 // preceding
		prevLogTerm := n.State.Log[toSendIdx-1].Term
		entriesToSend := n.State.Log[toSendIdx:] // TODO batching
		leaderCommit := n.State.CommitIndex

		request := data.AppendEntries{
			MessageCore:  msgCore,
			TermCore:     termCore,
			PrevLogIndex: prevLogIdx, // Preceding
			PrevLogTerm:  prevLogTerm,
			Entries:      entriesToSend,
			LeaderCommit: leaderCommit,
		}

		messages = append(messages, request)
	}

	n.ResetHeartbeatTimeout()
	return messages
}
