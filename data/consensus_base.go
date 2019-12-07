package data

const (
	// Consensus Algorithm Messages
	APPEND_MSG = MSG_TYPE("appendRequest")
	VOTE_MSG   = MSG_TYPE("voteRequest")

	// Responses to Consensus Messages
	APPEND_RES_MSG = MSG_TYPE("appendResponse")
	VOTE_RES_MSG   = MSG_TYPE("voteResponse")
)

const (
	UNKNOWN_LEADER = NODE_ID("FFFF")
)

type TERM_ID int

type TermCore struct {
	TermId TERM_ID `json:"term"`
}
