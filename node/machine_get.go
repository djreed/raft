package node

import "github.com/djreed/raft/data"

func HandleGet(n *Node, get data.GetMessage) data.MessageList {
	if n.IsLeader() {
		// Get val, return
		val := n.Get(get.Key)
		core := n.NewMessageCoreId(get.Source, data.OK_MSG, get.MessageId)
		body := &data.ReturnValue{
			Val: val,
		}
		msg := data.GetResponse{
			core,
			body,
		}
		return MakeList(msg)
	} else {
		// Redirect
		core := n.NewMessageCoreId(get.Source, data.REDIRECT_MSG, get.MessageId)
		return MakeList(core)
	}
}
