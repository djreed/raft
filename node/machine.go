package node

import "github.com/djreed/raft/data"

func (n *Node) StateMachine() error {
	for {
		var shouldRespond bool
		var response interface{}
		select {
		case msg := <-n.RequestVotes: // TODO
			shouldRespond, response = HandleRequestVote(n, msg)
			break

		case msg := <-n.AppendEntries: // TODO
			shouldRespond, response = HandleAppendEntries(n, msg)
			break

		case msg := <-n.RequestVoteResponses: // TODO
			shouldRespond, response = HandleRequestVoteResponse(n, msg)
			break

		case msg := <-n.AppendEntryResponses: // TODO
			shouldRespond, response = HandleAppendEntriesResponse(n, msg)
			break

		case msg := <-n.GetMessages: // TODO
			shouldRespond, response = HandleGet(n, msg)
			break

		case msg := <-n.PutMessages: // TODO
			shouldRespond, response = HandlePut(n, msg)
			break

			// case msg := <-n.ElectionTimeout: // TODO
			//
			// case msg := <-n.HeartbeatTimeout: // TODO

		}

		if shouldRespond {
			n.SendMessage(response)
		}
	}
}

func (n *Node) CreateResponseCore(msgType data.MSG_TYPE, msg data.MessageCore) data.MessageCore {
	return data.MessageCore{
		Source:    n.Id,
		Dest:      msg.Source,
		Leader:    n.Id, // n.Leader, // TODO TODO TODO
		Type:      msgType,
		MessageId: msg.MessageId,
	}
}

func HandleRequestVote(n *Node, vote data.RequestVote) (bool, interface{}) {
	return true, n.CreateResponseCore(data.FAIL, *vote.MessageCore)
}

func HandleAppendEntries(n *Node, appendEntries data.AppendEntries) (bool, interface{}) {
	return true, n.CreateResponseCore(data.FAIL, *appendEntries.MessageCore)
}

func HandleRequestVoteResponse(n *Node, voteRes data.RequestVoteResponse) (bool, interface{}) {
	return true, n.CreateResponseCore(data.FAIL, *voteRes.MessageCore)

}

func HandleAppendEntriesResponse(n *Node, appendRes data.AppendEntriesResponse) (bool, interface{}) {
	return true, n.CreateResponseCore(data.FAIL, *appendRes.MessageCore)

}

func HandleGet(n *Node, get data.GetMessage) (bool, interface{}) {
	val := n.Get(get.Key)
	core := n.CreateResponseCore(data.OK, *get.MessageCore)
	return true, data.GetResponse{MessageCore: &core, ReturnValue: &data.ReturnValue{val}}
}

func HandlePut(n *Node, put data.PutMessage) (bool, interface{}) {
	n.Set(put.Key, put.Val)
	core := n.CreateResponseCore(data.OK, *put.MessageCore)
	return true, core
}
