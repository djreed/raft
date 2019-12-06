package node

import "github.com/djreed/raft/data"

func HandleElectionTimeout(n *Node) data.MessageList {
	n.BecomeCandidate()
	messages := MakeList()
	for _, nodeId := range n.Neighbors {
		core := n.NewMessageCore(nodeId, data.VOTE_MSG)
		termCore := n.NewTermCore()

		request := data.RequestVote{
			MessageCore:  core,
			TermCore:     termCore,
			CandidateId:  n.Id,
			LastLogIndex: n.State.LastLogIndex(),
			LastLogTerm:  n.State.LastLogTerm(),
		}

		OUT.Printf("(%v) %v\n", n.Id, request)

		messages = append(messages, request)
	}

	return messages
}
