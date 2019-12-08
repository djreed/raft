package node

import (
	"encoding/json"
	"io"

	"github.com/djreed/raft/data"
)

func JSONStreams(c io.ReadWriter) (*json.Encoder, *json.Decoder) {
	encoder := json.NewEncoder(c)
	decoder := json.NewDecoder(c)
	return encoder, decoder
}

func (n *Node) HandleConn() {
	ERR.Printf("(%s) Listening to Socket: %s\n", n.Id, n.Socket.RemoteAddr())

	for {
		rawBytes := make([]byte, 2048)
		read, err := n.Socket.Read(rawBytes)
		if err != nil {
			// ERR.Printf("(%v) Failed to decode raw bytes:\n%s", n.Id, string(rawBytes))
			continue
		}

		byteData := rawBytes[:read]

		var messageCore data.MessageCore
		if err := json.Unmarshal(byteData, &messageCore); err != nil {
			// ERR.Printf("(%v) Failed to decode message core: %v", n.Id, err)
			// ERR.Printf("!!!\n%s", string(byteData))
			continue
		}
		// ERR.Printf("(RECEIVED %s) -- %s\n", n.Id, string(byteData))

		messageType := messageCore.Type

		// ERR.Printf("(%v) RECV %s", n.Id, messageType)

		// Decode JSON into correct message type
		// Send along corresponding channel
		var decodeErr error
		switch messageType {
		case data.GET_MSG:
			var getMsg data.GetMessage
			decodeErr = json.Unmarshal(byteData, &getMsg)
			n.GetMessages <- getMsg
			break

		case data.PUT_MSG:
			var putMsg data.PutMessage
			decodeErr = json.Unmarshal(byteData, &putMsg)
			n.PutMessages <- putMsg
			break

		case data.APPEND_MSG:
			var appendMsg data.AppendEntries
			decodeErr = json.Unmarshal(byteData, &appendMsg)
			n.AppendEntries <- appendMsg
			break

		case data.VOTE_MSG:
			var voteMsg data.RequestVote
			decodeErr = json.Unmarshal(byteData, &voteMsg)
			n.RequestVotes <- voteMsg
			break

		case data.APPEND_RES_MSG:
			var appendResponse data.AppendEntriesResponse
			decodeErr = json.Unmarshal(byteData, &appendResponse)
			n.AppendEntryResponses <- appendResponse
			break

		case data.VOTE_RES_MSG:
			var voteResponse data.RequestVoteResponse
			decodeErr = json.Unmarshal(byteData, &voteResponse)
			n.RequestVoteResponses <- voteResponse
			break

		default:
			ERR.Panicf("(!!! %s !!!) Unknown message type: %s\n", n.Id, messageType)
		}

		if decodeErr != nil {
			ERR.Panicf("(!!! %s !!!) %s\n", n.Id, decodeErr)
		}
	}
}

func (n *Node) SendMessage(msg interface{}) {
	// Read from Connection
	encoder, _ := JSONStreams(n.Socket)
	err := encoder.Encode(msg)
	if err != nil {
		ERR.Panicf("(%s) Failed to encode: %s\n", n.Id, err)
	} else {
		// byteData, _ := json.Marshal(msg)
		// var m data.MessageCore
		// json.Unmarshal(byteData, &m)
		// ERR.Printf("(%v) SEND %s", n.Id, m.Type)
	}
}
