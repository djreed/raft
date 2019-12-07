package node

import "github.com/djreed/raft/data"

func HandlePut(n *Node, put data.PutMessage) data.MessageList {
	if n.IsLeader() {
		logEntry := data.LogEntry{
			Key:    put.Key,
			Value:  put.Val,
			Type:   data.PUT,
			Term:   n.State.CurrentTerm,
			Sender: put.Source,
			MID:    put.MessageId,
		}

		n.State.AppendLog(logEntry)
		return nil
	} else {
		// Redirect
		core := n.NewMessageCoreId(put.Source, data.REDIRECT_MSG, put.MessageId)
		if n.Leader == data.UNKNOWN_LEADER {
			core.Type = data.FAIL_MSG // No known leader, unable to redirect
		}
		return MakeList(core)
	}
}
