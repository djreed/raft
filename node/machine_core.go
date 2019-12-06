package node

import "github.com/djreed/raft/data"

func (n *Node) StateMachine() error {
	for {
		var shouldRespond bool
		var responses []interface{}
		select {
		case rv := <-n.RequestVotes:
			responses = HandleRequestVote(n, rv)
			break

		case ae := <-n.AppendEntries:
			responses = HandleAppendEntries(n, ae)
			break

		case rvr := <-n.RequestVoteResponses:
			responses = HandleRequestVoteResponse(n, rvr)
			break

		case aer := <-n.AppendEntryResponses:
			responses = HandleAppendEntriesResponse(n, aer)
			break

		case get := <-n.GetMessages:
			responses = HandleGet(n, get)
			break

		case put := <-n.PutMessages:
			responses = HandlePut(n, put)
			break

		case <-n.ElectionTimeout:
			responses = HandleElectionTimeout(n)
			break

			// case <-n.HeartbeatTimeout: // TODO

		}

		if shouldRespond {
			for _, response := range responses {
				n.SendMessage(response)
			}
		}
	}
}

func CreateResponseCore(n *Node, msgType data.MSG_TYPE, msg data.MessageCore) *data.MessageCore {
	return &data.MessageCore{
		Source:    n.Id,
		Dest:      msg.Source,
		Leader:    n.Id, // n.Leader, // TODO TODO TODO
		Type:      msgType,
		MessageId: msg.MessageId,
	}
}

func UpToDate(n *Node, lastLogIndex data.ENTRY_INDEX, lastLogTerm data.TERM_ID) bool {
	// TODO
	return true
}

func MakeList(data ...interface{}) []interface{} {
	vals := make([]interface{}, 0)
	return append(vals, data)
}
