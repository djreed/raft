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
	OUT.Printf("(%s) Listening to Socket: %s\n", n.Id, n.Socket.RemoteAddr())

	// Read from Connection
	_, decoder := JSONStreams(n.Socket)

	for {
		var baseMsg data.UnknownMessage
		decoder.Decode(&baseMsg)

		byteData, _ := json.Marshal(baseMsg)

		// Decode JSON into correct message type
		// Send along corresponding channel

		var decodeErr error
		// Ravioli Ravioli Give Me The PANIC: NIL POINTER
		switch data.MSG_TYPE(baseMsg["type"].(string)) {
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
			OUT.Panicf("(!!! %s !!!) Unknown message type: %s\n", n.Id, baseMsg["type"])
		}

		if decodeErr != nil {
			OUT.Panicf("(!!! %s !!!) %s\n", n.Id, decodeErr)
		} else {
			OUT.Printf("(RECEIVED %s) -- %s\n", n.Id, string(byteData))
		}
	}
}

func (n *Node) SendMessage(msg interface{}) {
	// Read from Connection
	encoder, _ := JSONStreams(n.Socket)
	err := encoder.Encode(msg)
	if err != nil {
		OUT.Panicf("(!!! %s !!!) -- %s\n", n.Id, err)
	} else {
		byteData, _ := json.Marshal(msg)
		OUT.Printf("(SENDING %s) -- %s", n.Id, string(byteData))
	}
}
