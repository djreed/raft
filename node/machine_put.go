package node

import "github.com/djreed/raft/data"

func HandlePut(n *Node, put data.PutMessage) data.MessageList {
	n.Set(put.Key, put.Val)
	core := CreateResponseCore(n, data.OK_MSG, *put.MessageCore)
	return MakeList(core)
}
