package node

import "github.com/djreed/raft/data"

// TODO if we get a term > currentTerm, convert to Follower, set term to higher
func HandleAppendEntries(n *Node, appendEntries data.AppendEntries) data.MessageList {
	n.ResetElectionTimeout()
	n.SetLeader(appendEntries.Leader)
	// TODO
	return nil
}

func HandleAppendEntriesResponse(n *Node, appendRes data.AppendEntriesResponse) data.MessageList {
	// TODO
	return nil
}
