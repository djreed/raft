package main

import (
	"os"

	"github.com/djreed/raft/data"
	"github.com/djreed/raft/node"
)

func main() {
	if len(os.Args) < 3 {
		panic("Must specify Node ID and at least one Neighbor")
	}
	// Get args as IPs and relations
	nodeId := data.NODE_ID(os.Args[1])
	neighbors := AsNodeIds(os.Args[2:])

	node := node.NewNode(nodeId, neighbors)
	node.StartNode()
}

func AsNodeIds(ids []string) []data.NODE_ID {
	nodeIds := make([]data.NODE_ID, len(ids))
	for idx, nodeId := range ids {
		nodeIds[idx] = data.NODE_ID(nodeId)
	}
	return nodeIds
}
