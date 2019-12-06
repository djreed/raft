package node

import "github.com/djreed/raft/data"

// TODO if we get a term > currentTerm, convert to Follower, set term to higher
func HandleRequestVote(n *Node, vote data.RequestVote) data.MessageList {
	responseCore := CreateResponseCore(n, data.OK_MSG, *vote.MessageCore)

	responseData := &data.RequestVoteResponseData{
		VoteGranted: false,
	}

	response := data.RequestVoteResponse{
		responseCore,
		vote.TermCore,
		responseData,
	}

	// 1. Reply false if term < currentTerm (§5.1)
	if vote.TermId >= n.State.CurrentTerm {
		//2. If votedFor is null or candidateId...
		if n.State.VotedFor == "" || n.State.VotedFor == vote.CandidateId {
			// and candidate’s log is at least as up-to-date as receiver’s log...
			if UpToDate(n, vote.LastLogIndex, vote.LastLogTerm) {
				responseData.VoteGranted = true
				VoteFor(n, vote.CandidateId) // grant vote (§5.2, §5.4)
			}
		}
	}

	return MakeList(response)
}

func HandleRequestVoteResponse(n *Node, voteRes data.RequestVoteResponse) data.MessageList {
	return MakeList(CreateResponseCore(n, data.FAIL_MSG, *voteRes.MessageCore))
}

func VoteFor(n *Node, candidate data.NODE_ID) {
	n.State.VotedFor = candidate
}
