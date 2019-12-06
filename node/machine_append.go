package node

import "github.com/djreed/raft/data"

// TODO if we get a term > currentTerm, convert to Follower, set term to higher
func HandleAppendEntries(n *Node, appendEntries data.AppendEntries) data.MessageList {
	return MakeList(n.NewMessageCore(appendEntries.Source, data.FAIL_MSG))
}

// func CreateAppendEntriesResponse() data.AppendEntriesResponse {
//
// 	/*
// 	  type AppendEntriesResponse struct {
// 	  	*MessageCore
// 	  	*TermCore
// 	  	*AppendEntriesResponseData
// 	  }
// 	*/
//
// }

func HandleAppendEntriesResponse(n *Node, appendRes data.AppendEntriesResponse) data.MessageList {
	return MakeList(n.NewMessageCore(appendRes.Source, data.FAIL_MSG))

}
