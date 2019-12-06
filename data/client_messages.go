package data

const (
	// Interface (Client) Messages
	GET_MSG      = MSG_TYPE("get")
	PUT_MSG      = MSG_TYPE("put")
	OK_MSG       = MSG_TYPE("ok")
	FAIL_MSG     = MSG_TYPE("fail")
	REDIRECT_MSG = MSG_TYPE("redirect")
)

///
// Key / Value
///
type KEY_TYPE string
type VAL_TYPE string

///
// Client GET
///

type GetMessage struct {
	*MessageCore
	Key KEY_TYPE `json:"key"`
}

type GetResponse struct {
	*MessageCore
	Val VAL_TYPE `json:"value"`
	// Type of "ok"
}

type GetFail struct {
	*MessageCore
	// Type of "fail"
}

///
// Client PUT
///

type PutMessage struct {
	*MessageCore
	Key KEY_TYPE `json:"key"`
	Val VAL_TYPE `json:"value"`
}

type PutResponse struct {
	*MessageCore
	// Type of "ok" or "fail"
}

///
// Request to non-leader Node
///

type RedirectMessage struct {
	*MessageCore
	// Type of "redirect"
}
