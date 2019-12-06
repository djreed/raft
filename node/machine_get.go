package node

import "github.com/djreed/raft/data"

func HandleGet(n *Node, get data.GetMessage) data.MessageList {
	if n.IsLeader() {
		// Get val, return
		val := n.Get(get.Key)
		core := n.NewMessageCoreId(get.Source, data.OK_MSG, get.MessageId)
		msg := data.GetResponse{
			MessageCore: core,
			Val:         val,
		}
		return MakeList(msg)
	} else {
		// Redirect
		core := n.NewMessageCoreId(get.Source, data.REDIRECT_MSG, get.MessageId)
		return MakeList(core)
	}
}
