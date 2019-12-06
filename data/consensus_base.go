package data

const (
	// Consensus Algorithm Messages
	APPEND = MSG_TYPE("appendRequest")
	VOTE   = MSG_TYPE("voteRequest")

	// Responses to Consensus Messages
	APPEND_RES = MSG_TYPE("appendResponse")
	VOTE_RES   = MSG_TYPE("voteResponse")
)

const (
	UNKNOWN_LEADER = NODE_ID("FFFF")
)

type TERM_ID uint32

type TermCore struct {
	TermId TERM_ID `json:"term"`
}
