package node

import "github.com/djreed/raft/data"

// TODO if we get a term > currentTerm, convert to Follower, set term to higher
func HandleRequestVote(n *Node, vote data.RequestVote) data.MessageList {
	n.HandleTermUpdate(vote.TermId, data.UNKNOWN_LEADER)

	responseCore := CreateResponseCore(n, data.VOTE_RES_MSG, *vote.MessageCore)

	response := data.RequestVoteResponse{
		MessageCore: responseCore,
		TermCore:    vote.TermCore,
		VoteGranted: false,
	}

	// 1. Reply false if term < currentTerm (§5.1)
	if vote.TermId >= n.State.CurrentTerm {
		// 2. If votedFor is null or candidateId...
		if n.State.VotedFor == "" || n.State.VotedFor == vote.CandidateId {
			// and candidate’s log is at least as up-to-date as receiver’s log...
			if n.TargetUpToDate(vote.LastLogIndex, vote.LastLogTerm) {
				response.VoteGranted = true
				n.State.SetVotedFor(vote.CandidateId) // grant vote (§5.2, §5.4)
				n.ResetElectionTimeout()
			}
		}
	}

	return MakeList(response)
}

// TODO handle term > currentTerm
func HandleRequestVoteResponse(n *Node, voteRes data.RequestVoteResponse) data.MessageList {
	// If the term of the response is greater than our own,
	// we know that we're not getting VoteGranted
	if n.HandleTermUpdate(voteRes.TermId, data.UNKNOWN_LEADER) {
		return nil
	}

	if voteRes.VoteGranted {
		n.IncrementVotes()
		if n.VoteQuorum() {
			ERR.Printf("(%v) !!! I AM NOW THE LEADER, BOW BEFORE ME !!!", n.Id)
			n.BecomeLeader()
		}
	}

	return nil
}
