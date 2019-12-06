package node

import "github.com/djreed/raft/data"

func HandleElectionTimeout(n *Node) data.MessageList {
	/*
		  • On conversion to candidate, start election:
		    • Increment currentTerm
			  • Vote for self
			  • Reset election timer
		  • Send RequestVote RPCs to all other servers
		  • If votes received from majority of servers: become leader
		  • If AppendEntries RPC received from new leader: convert to follower
		  • If election timeout elapses: start new election
	*/
	BecomeCandidate(n)

	return nil
}
