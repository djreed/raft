package node

import "github.com/djreed/raft/data"

func (n *Node) StateMachine() error {
	for {
		var responses []interface{}
		select {
		case <-n.ElectionTimeout:
			OUT.Printf("(%v) -- !!! ELECTION TIMEOUT !!!", n.Id)
			responses = HandleElectionTimeout(n)

			// case <-n.HeartbeatTimeout: // TODO

		case rvr := <-n.RequestVoteResponses:
			responses = HandleRequestVoteResponse(n, rvr)

		case aer := <-n.AppendEntryResponses:
			responses = HandleAppendEntriesResponse(n, aer)

		case rv := <-n.RequestVotes:
			responses = HandleRequestVote(n, rv)

		case ae := <-n.AppendEntries:
			responses = HandleAppendEntries(n, ae)

		case get := <-n.GetMessages:
			responses = HandleGet(n, get)

		case put := <-n.PutMessages:
			responses = HandlePut(n, put)

		}

		for _, response := range responses {
			n.SendMessage(response)
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

func UpToDate(n *Node, lastLogIndex data.ENTRY_INDEX, lastLogTerm data.TERM_ID) bool {
	// TODO
	return true
}

func MakeList(data ...interface{}) []interface{} {
	return data
}
