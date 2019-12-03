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
	*KeyRequest
}

type ReturnValue struct {
	Val VAL_TYPE `json:"value"`
}

type GetResponse struct {
	*MessageCore
	*ReturnValue
	// Type of "ok"
}

type GetFail struct {
	*MessageCore
	// Type of "fail"
}

///
// Client PUT
///

type PutMessageBody struct {
	Key KEY_TYPE `json:"key"`
	Val VAL_TYPE `json:"value"`
}

type PutMessage struct {
	*MessageCore
	*PutMessageBody
}

type PutResponse struct {
	*MessageCore
	// Type of "ok" or "fail"
}

///
// Request to non-leader Node
///

// TODO
type RedirectMessage struct {
	*MessageCore
	// Type of "redirect"
}
