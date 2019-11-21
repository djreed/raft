package data

type GetMessageBody struct {
	Key KEY_TYPE `json:"key"`
}

type GetMessage struct {
	*MessageCore
	*MessageIdBase
	*GetMessageBody
}

type GetResponseBody struct {
	Value VAL_TYPE `json:"key,omitempty"`
	// TODO confirm whether empty string is a valid VAL_TYPE
}

type GetResponse struct {
	*MessageCore
	*MessageIdBase
	*GetResponseBody
}
