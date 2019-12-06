package node

import "github.com/djreed/raft/data"

func HandleElectionTimeout(n *Node) data.MessageList {
	n.BecomeCandidate()
	messages := MakeList()
	for _, nodeId := range n.Neighbors {
		core := n.NewMessageCore(nodeId, data.VOTE_MSG)
		termCore := n.NewTermCore()
		voteCore := CreateRequestVoteCore(n)

		/*
		   Source    NODE_ID    `json:"src"`
		   Dest      NODE_ID    `json:"dst"`
		   Leader    NODE_ID    `json:"leader"`
		   Type      MSG_TYPE   `json:"type"`
		   MessageId MESSAGE_ID `json:"MID"`
		*/

		request := data.RequestVote{
			core,
			termCore,
			voteCore,
		}

		messages = append(messages, request)
	}

	return messages
}

func CreateRequestVoteCore(n *Node) *data.RequestVoteCore {
	// type RequestVoteCore struct {
	// 	// NOTE: Term ID from TermCore
	//
	// 	// The Node ID of the Candidate
	// 	CandidateId NODE_ID `json:"candidateId"`
	//
	// 	// The Index of the Candidate's last log entry
	// 	LastLogIndex ENTRY_INDEX `json:"lastLogIndex"`
	//
	// 	// The Term of the Candidate's last log entry
	// 	LastLogTerm TERM_ID `json:"lastLogTerm"`
	// }

	return &data.RequestVoteCore{
		CandidateId:  n.Id,
		LastLogIndex: 10,
		LastLogTerm:  10,
		// TODO LastLogIndex and LastLogTerm
	}
}
