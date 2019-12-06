package node

import "github.com/djreed/raft/data"

func HandlePut(n *Node, put data.PutMessage) data.MessageList {
	if n.IsLeader() {
		// Set val, return (TODO quorum and replication)
		n.Set(put.Key, put.Val)
		core := n.NewMessageCoreId(put.Source, data.OK_MSG, put.MessageId)
		return MakeList(core)
	} else {
		// Redirect
		core := n.NewMessageCoreId(put.Source, data.REDIRECT_MSG, put.MessageId)
		return MakeList(core)
	}
}
