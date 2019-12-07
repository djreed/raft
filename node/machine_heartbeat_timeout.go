package node

import "github.com/djreed/raft/data"

func HandleHeartbeatTimeout(n *Node) (data.MessageList, bool) {
	messages := make(data.MessageList, 0)
	for _, nodeId := range n.Neighbors {
		msgCore := n.NewMessageCore(nodeId, data.APPEND_MSG)
		termCore := n.NewTermCore()

		sendStartIdx := n.SendStartIdx(nodeId) // Starting index of logs being sent

		prevLogIdx := sendStartIdx - 1 // Preceding index to logs being sent
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

		ERR.Printf("!!! (%v) !!! -- lastLogIdx(%v) | sendStart(%v) | nextIdxToSend(%+v)",
			n.Id, n.State.LastLogIndex(), sendStartIdx, n.State.NextIndex)

		// Are we sending data that's within our logs?
		dataToSend := sendStartIdx <= n.State.LastLogIndex()
		if dataToSend {
			entriesToSend := n.State.Log[sendStartIdx-1:] // TODO batching
			request.Entries = entriesToSend

			messagesSent := len(entriesToSend)

			// Is non-committed data being sent
			lastIndexSent := int(sendStartIdx) + messagesSent - 1 // TODO validate minus 1
			if lastIndexSent > int(n.State.CommitIndex) {         // TODO batching
				n.AddReplicationMid(request.MessageId)
			}
		}

		n.RecordAppendMessage(request)
		messages = append(messages, request)
	}

	n.ResetHeartbeatTimeout()
	return messages, n.PendingCommits()
}
