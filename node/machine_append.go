package node

import "github.com/djreed/raft/data"

// TODO if we get a term > currentTerm, convert to Follower, set term to higher
func HandleAppendEntries(n *Node, appendEntries data.AppendEntries) data.MessageList {
	return MakeList(n.CreateResponseCore(data.FAIL, *appendEntries.MessageCore))
}

func HandleAppendEntriesResponse(n *Node, appendRes data.AppendEntriesResponse) data.MessageList {
	return MakeList(n.CreateResponseCore(data.FAIL, *appendRes.MessageCore))

}
