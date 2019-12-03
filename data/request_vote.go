package data

type RequestVoteCore struct {
	// NOTE: Term ID from TermCore

	// The Node ID of the Candidate
	CandidateId NODE_ID `json:"candidateId"`

	// The Index of the Candidate's last log entry
	LastLogIndex ENTRY_INDEX `json:"lastLogIndex"`

	// The Term of the Candidate's last log entry
	LastLogTerm TERM_ID `json:"lastLogTerm"`
}

type RequestVote struct {
	*MessageCore
	*TermCore
	*RequestVoteCore
}

type RequestVoteResponseData struct {
	// Whether the vote has been granted to the Candidate
	VoteGranted bool `json:"voteGranted"`
}

type RequestVoteResponse struct {
	*MessageCore
	*TermCore
	*RequestVoteResponseData
}
