package node

import "github.com/djreed/raft/data"

func HandleGet(n *Node, get data.GetMessage) data.MessageList {
	val := n.Get(get.Key)
	core := n.CreateResponseCore(data.OK, *get.MessageCore)
	return MakeList(data.GetResponse{MessageCore: core, ReturnValue: &data.ReturnValue{val}})
}
