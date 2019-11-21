package data

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
