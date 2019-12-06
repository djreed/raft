package data

import (
	"crypto/rand"
	"encoding/base64"
)

/*
  All Messages contain AT LEAST:
  - src
  - dest
  - leader
  - type
*/
type NODE_ID string

type MessageCore struct {
	Source    NODE_ID    `json:"src"`
	Dest      NODE_ID    `json:"dst"`
	Leader    NODE_ID    `json:"leader"`
	Type      MSG_TYPE   `json:"type"`
	MessageId MESSAGE_ID `json:"MID"`
}

// Anonymous incoming message
type UnknownMessage = map[string]interface{}

// MESSAGE_ID is a uniquely generated string
type MESSAGE_ID string

type MSG_TYPE string

type MessageList = []interface{}

const MID_LEN = 8

func NewMessageId() MESSAGE_ID {
	str, _ := GenerateRandomString(MID_LEN)
	return MESSAGE_ID(str)
}

// Source: https://flaviocopes.com/go-random/

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
