package data

const (
	// Consensus Algorithm Messages
	APPEND = MSG_TYPE("append")
	VOTE   = MSG_TYPE("vote")

	// Append Types
	PROMISE = MSG_TYPE("promise")
	COMMIT  = MSG_TYPE("commit")
)

type TERM_ID uint32

type TermCore struct {
	TermId TERM_ID `json:"term"`
}
