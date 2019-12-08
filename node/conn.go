package node

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"sync"

	"github.com/djreed/raft/data"
)

var mut sync.Mutex

func JSONStreams(c io.ReadWriter) (*json.Encoder, *json.Decoder) {
	encoder := json.NewEncoder(c)
	decoder := json.NewDecoder(c)
	return encoder, decoder
}

func (n *Node) HandleConn() {
	ERR.Printf("(%s) Listening to Socket: %s\n", n.Id, n.Socket.RemoteAddr())

	// Read from Connection
	_, decoder := JSONStreams(n.Socket)

	for {
		var baseMsg data.UnknownMessage

		// var rawBytes = make([2048]bytes{}, 2048)

		if err := decoder.Decode(&baseMsg); err != nil || baseMsg == nil {
			buf := decoder.Buffered()
			OUT.Printf("(%v) FAILED TO DECODE: %v", n.Id, err)
			b, _ := ioutil.ReadAll(buf)
			OUT.Printf("!!!!!\n%s\n!!!", string(b))
			continue
		}

		byteData, _ := json.Marshal(baseMsg)
		var messageCore data.MessageCore
		if err := json.Unmarshal(byteData, &messageCore); err != nil {
			OUT.Printf("(%v) FAILED TO DECODE MESSAGE CORE: %v", n.Id, err)
			OUT.Printf("!!!\n%s\n!!!", string(byteData))
			continue
		}
		ERR.Printf("(RECEIVED %s) -- %s\n", n.Id, string(byteData))

		messageType := messageCore.Type
		OUT.Printf("(%v) RECV %s", n.Id, messageType)

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
			ERR.Panicf("(!!! %s !!!) Unknown message type: %s\n", n.Id, baseMsg["type"])
		}

		if decodeErr != nil {
			ERR.Panicf("(!!! %s !!!) %s\n", n.Id, decodeErr)
		}
	}
}

func (n *Node) SendMessage(msg interface{}) {
	// Read from Connection
	mut.Lock()
	encoder, _ := JSONStreams(n.Socket)
	err := encoder.Encode(msg)
	if err != nil {
		ERR.Panicf("(!!! %s !!!) -- %s\n", n.Id, err)
	} else {
		byteData, _ := json.Marshal(msg)
		var m data.MessageCore
		json.Unmarshal(byteData, &m)
		OUT.Printf("(%v) SEND %s", n.Id, m.Type)
	}
	mut.Unlock()
}
