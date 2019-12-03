package node

import (
	"net"
	"time"

	"github.com/djreed/raft/data"
	"github.com/djreed/raft/logging"
)

var OUT = logging.OUT
var ERR = logging.ERR

const (
	CHAN_BUFFER = 128
)

type Socket = net.Conn

// Node contains Network state
type Node struct {
	Id data.NODE_ID

	Neighbors []data.NODE_ID
	Leader    data.NODE_ID

	// Single UNIX domain socket, emulating Ethernet
	Socket Socket

	Data      map[data.KEY_TYPE]data.VAL_TYPE
	RaftState data.RaftState

	// A Channel for each data type we need to handle
	RequestVotes         chan data.RequestVote
	AppendEntries        chan data.AppendEntries
	RequestVoteResponses chan data.RequestVoteResponse
	AppendEntryResponses chan data.AppendEntriesResponse

	// Client messages
	GetMessages chan data.GetMessage
	PutMessages chan data.PutMessage

	// On timeout, start new election cycle
	ElectionTimeout <-chan time.Time

	// On timeout, send empty AppendEntries
	HeartbeatTimeout <-chan time.Time
}

func NewNode(id data.NODE_ID, neighbors []data.NODE_ID) Node {
	unixSock := OpenSocket(id)

	initialRaftState := data.NewRaftState(len(neighbors))

	initialNode := Node{
		Id:            id,
		Neighbors:     neighbors,
		Socket:        unixSock,
		Data:          make(map[data.KEY_TYPE]data.VAL_TYPE),
		RaftState:     initialRaftState,
		RequestVotes:  make(chan data.RequestVote, CHAN_BUFFER),
		AppendEntries: make(chan data.AppendEntries, CHAN_BUFFER),
		GetMessages:   make(chan data.GetMessage, CHAN_BUFFER),
		PutMessages:   make(chan data.PutMessage, CHAN_BUFFER),
	}

	return initialNode
}

func OpenSocket(id data.NODE_ID) Socket {
	conn, err := net.Dial("unixpacket", string(id))
	if err != nil {
		ERR.Panic(err)
	}

	return conn
}

func (n *Node) StartNode() error {
	go n.HandleConn()
	return n.StateMachine()
}

func (n *Node) Get(key data.KEY_TYPE) data.VAL_TYPE {
	return n.Data[key]
}

func (n *Node) Set(key data.KEY_TYPE, val data.VAL_TYPE) {
	n.Data[key] = val
}

func (n *Node) IsLeader() bool {
	return n.Id == n.Leader
}
