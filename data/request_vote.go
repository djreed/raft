package data

type RequestVote struct {
	*MessageCore
	*TermCore

	// The Node ID of the Candidate
	CandidateId NODE_ID `json:"candidateId"`

	// The Index of the Candidate's last log entry
	LastLogIndex ENTRY_INDEX `json:"lastLogIndex"`

	// The Term of the Candidate's last log entry
	LastLogTerm TERM_ID `json:"lastLogTerm"`
}

type RequestVoteResponse struct {
	*MessageCore
	*TermCore

	VoteGranted bool `json:"voteGranted"`
}
