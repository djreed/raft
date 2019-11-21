package data

/*
  All messages must contain the following:
  "src": "<ID>",
  "dst": "<ID>",
  "leader": "<ID>",
  "type": "redirect"
*/

// NODE_ID is a hex ID, or FFFF if unknown
type NODE_ID string

type MessageCore struct {
	Source NODE_ID  `json:"src"`
	Dest   NODE_ID  `json:"dst"`
	Leader NODE_ID  `json:"leader"`
	Type   MSG_TYPE `json:"type"`
}

// MESSAGE_ID is a uniquely generated string
type MESSAGE_ID string

type MessageIdBase struct {
	MessageId MESSAGE_ID `json:"MID"`
}

type MSG_TYPE string

const (
	// Interface (Client) Messages
	GET      = MSG_TYPE("get")
	PUT      = MSG_TYPE("put")
	OK       = MSG_TYPE("ok")
	FAIL     = MSG_TYPE("fail")
	REDIRECT = MSG_TYPE("redirect")

	// Consensus Algorithms
	APPEND = MSG_TYPE("append")
	VOTE   = MSG_TYPE("vote")

	// Append Types
	PROMISE = MSG_TYPE("promise")
	COMMIT  = MSG_TYPE("commit")
)
