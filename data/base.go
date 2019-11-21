package data

/*
  All Messages contain AT LEAST:
  - src
  - dest
  - leader
  - type
*/
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
