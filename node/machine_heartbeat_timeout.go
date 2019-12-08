package node

import "github.com/djreed/raft/data"

const BATCH_SIZE = 10

func HandleHeartbeatTimeout(n *Node) (data.MessageList, bool) {
	messages := make(data.MessageList, 0)

	// ERR.Printf("(%v) -- Heartbeat State: lastLogIdx(%v) | nextIdxToSend(%+v)",
	// 	n.Id, n.State.LastLogIndex(), n.State.NextIndex)

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

		// Are we sending data that's within our logs?
		dataToSend := sendStartIdx <= n.State.LastLogIndex()
		if dataToSend {
			toSendCount := Min(BATCH_SIZE, len(n.State.Log[sendStartIdx-1:]))

			request.Entries = n.State.Log[sendStartIdx-1 : (sendStartIdx-1)+data.ENTRY_INDEX(toSendCount)]
			messagesSent := len(request.Entries)

			// Is non-committed data being sent?
			lastIndexSent := int(sendStartIdx) + messagesSent - 1
			if lastIndexSent > int(n.State.CommitIndex) {
				n.AddReplicationMid(request.MessageId)
			}
		}

		n.RecordAppendMessage(request)
		messages = append(messages, request)
	}

	n.ResetHeartbeatTimeout()
	return messages, n.PendingCommits()
}
