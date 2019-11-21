package data

///
// Current Term
///

// Term the message corresponds to
type TERM_ID uint32

type TermCore struct {
	TermId TERM_ID `json:"term"`
}

///
// Log Append
///

type APPEND_ID uint32

type UpdateValue struct {
	AppendId APPEND_ID `json:"id"`
	Key      KEY_TYPE  `json:"key"`
	Value    VAL_TYPE  `json:"val"`
	Type     MSG_TYPE  `json:"type"`
}

type UpdatePayload struct {
	Updates []UpdateValue `json:"updates"`
}

type AppendMessage struct {
	*MessageCore
	*MessageIdBase
	*TermCore
	*UpdatePayload
}

///
// Leader Election Votes
///

type VoteMessageBody struct {
	VoteId NODE_ID `json:"vote"`
}

type VoteMessage struct {
	*MessageCore
	*MessageIdBase
	*TermCore
	*VoteMessageBody
}
