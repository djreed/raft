package data

const (
	// Interface (Client) Messages
	GET      = MSG_TYPE("get")
	PUT      = MSG_TYPE("put")
	OK       = MSG_TYPE("ok")
	FAIL     = MSG_TYPE("fail")
	REDIRECT = MSG_TYPE("redirect")
)

///
// Key / Value
///
type KEY_TYPE string
type VAL_TYPE string

///
// Client GET
///

type KeyRequest struct {
	Key KEY_TYPE `json:"key"`
}

type GetMessage struct {
	*MessageCore
	*MessageIdBase
	*KeyRequest
}

type ReturnValue struct {
	Value VAL_TYPE `json:"key"`
}

type GetResponse struct {
	*MessageCore
	*MessageIdBase
	*ReturnValue
}

type GetFail struct {
	*MessageCore
	*MessageIdBase
}

///
// Client PUT
///

type PutMessageBody struct {
	Key   KEY_TYPE `json:"key"`
	Value VAL_TYPE `json:"value"`
}

type PutMessage struct {
	*MessageCore
	*MessageIdBase
	*PutMessageBody
}

type PutResponse struct {
	*MessageCore
	*MessageIdBase
}

///
// Request to non-leader Node
///

type RedirectMessage struct {
	*MessageCore
	*MessageIdBase
}
