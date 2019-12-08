package node

import (
	"encoding/json"

	"github.com/djreed/raft/data"
)

func (n *Node) StateMachineSteady() error {
	stateChange := false
	for {
		var responses []interface{}
		select {
		case rvr := <-n.RequestVoteResponses:
			ERR.Printf("(%v) -- Received RequestVoteResponse from [%v]: %v", n.Id, rvr.Source, rvr.VoteGranted)
			responses = HandleRequestVoteResponse(n, rvr)
			break

		case aer := <-n.AppendEntryResponses:
			responses, _ = HandleAppendEntriesResponse(n, aer, false)
			break

		case rv := <-n.RequestVotes:
			b, _ := json.Marshal(rv)
			ERR.Printf("(%v) -- Received RequestVote: %s", n.Id, string(b))

			responses = HandleRequestVote(n, rv)

			b2, _ := json.Marshal(responses[0])
			ERR.Printf("(%v) -- Sending RequestVoteResponse: %s", n.Id, string(b2))

			break

		case ae := <-n.AppendEntries:
			responses = HandleAppendEntries(n, ae)
			break

		case get := <-n.GetMessages:
			responses = HandleGet(n, get)
			break

		case put := <-n.PutMessages:
			responses = HandlePut(n, put)
			if n.IsLeader() {
				stateChange = true
			}
			break

		case <-n.ElectionTimeout:
			if n.IsLeader() {
				ERR.Print("\n\n\n\n\n\n\n\n\n\n\nWHAT THE FUCK\n\n\n\n\n\n\n\n\n\n\n\n\n")
				ERR.Print("\n\n\n\n\n\n\n\n\n\n\nWHAT THE FUCK\n\n\n\n\n\n\n\n\n\n\n\n\n")
				ERR.Print("\n\n\n\n\n\n\n\n\n\n\nWHAT THE FUCK\n\n\n\n\n\n\n\n\n\n\n\n\n")
				ERR.Print("\n\n\n\n\n\n\n\n\n\n\nWHAT THE FUCK\n\n\n\n\n\n\n\n\n\n\n\n\n")
				ERR.Print("\n\n\n\n\n\n\n\n\n\n\nWHAT THE FUCK\n\n\n\n\n\n\n\n\n\n\n\n\n")
			}
			ERR.Printf("(%v) -- !!! ELECTION TIMEOUT !!!", n.Id)
			responses = HandleElectionTimeout(n)
			ERR.Printf("(%v) -- Sending RequestVotes:", n.Id)
			for _, res := range responses {
				b, _ := json.Marshal(res)
				ERR.Printf("%v", string(b))
			}
			break

		case <-n.HeartbeatTimeout:
			ERR.Printf("(%v) -- HEARTBEAT TIMEOUT FROM !!! STEADY !!!", n.Id)
			responses, _ = HandleHeartbeatTimeout(n)
			break
		}

		if len(responses) > 0 {
			for _, response := range responses {
				n.SendMessage(response)
			}
		}

		if stateChange {
			ERR.Printf("(!!! %v !!!) -- STATE CHANGE STEADY -> COMMIT", n.Id)
			n.BeginCommit()
			n.StateMachineCommit()
			stateChange = false
			ERR.Printf("(!!! %v !!!) -- BACK IN STEADY", n.Id)
		}
	}
}

func (n *Node) StateMachineCommit() {
	stateChange := false
	for {
		var responses []interface{}
		select {
		case rvr := <-n.RequestVoteResponses:
			ERR.Printf("(%v) -- Received RequestVoteResponse from [%v]: %v", n.Id, rvr.Source, rvr.VoteGranted)

			responses = HandleRequestVoteResponse(n, rvr)
			stateChange = true // On election, drop back down to Steady
			break

		case aer := <-n.AppendEntryResponses:
			responses, stateChange = HandleAppendEntriesResponse(n, aer, true)
			break

		case rv := <-n.RequestVotes:
			b, _ := json.Marshal(rv)
			ERR.Printf("(%v) -- Received RequestVote: %s", n.Id, string(b))

			responses = HandleRequestVote(n, rv)

			b2, _ := json.Marshal(responses[0])
			ERR.Printf("(%v) -- Sending RequestVoteResponse: %s", n.Id, string(b2))

			stateChange = true // On election, drop back down to steady
			break

		case ae := <-n.AppendEntries:
			responses = HandleAppendEntries(n, ae)
			break

		case get := <-n.GetMessages:
			responses = HandleGet(n, get)
			break

		case <-n.ElectionTimeout:
			ERR.Printf("(%v) -- !!! ELECTION TIMEOUT !!!", n.Id)
			responses = HandleElectionTimeout(n)
			break

		case <-n.HeartbeatTimeout:
			ERR.Printf("(%v) -- HEARTBEAT TIMEOUT FROM !!! COMMIT !!!", n.Id)
			responses, _ = HandleHeartbeatTimeout(n) // TODO handle batching
			break
		}

		if len(responses) > 0 {
			for _, response := range responses {
				n.SendMessage(response)
			}
		}

		if stateChange {
			ERR.Printf("(!!! %v !!!) -- STATE CHANGE COMMIT -> STEADY", n.Id)
			n.EndCommit()
			return
		}
	}
}

func CreateResponseCore(n *Node, msgType data.MSG_TYPE, msg data.MessageCore) *data.MessageCore {
	return &data.MessageCore{
		Source:    n.Id,
		Dest:      msg.Source,
		Leader:    n.Leader,
		Type:      msgType,
		MessageId: msg.MessageId,
	}
}

func MakeList(data ...interface{}) []interface{} {
	return data
}
