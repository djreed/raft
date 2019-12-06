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

	Role data.NODE_STATE

	// Single UNIX domain socket, emulating Ethernet
	Socket Socket

	Data  map[data.KEY_TYPE]data.VAL_TYPE
	State data.RaftState

	// Quorum Tracking
	Votes int // How many neighbors have voted for me

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
		State:         initialRaftState,
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
		OUT.Panic(err)
	}

	return conn
}

func (n *Node) StartNode() error {
	go n.HandleConn()
	n.BecomeFollower(data.UNKNOWN_LEADER)
	return n.StateMachine()
}

func (n *Node) NewMessageCore(dest data.NODE_ID, msgType data.MSG_TYPE) *data.MessageCore {
	return n.NewMessageCoreId(dest, msgType, data.NewMessageId())
}

func (n *Node) NewMessageCoreId(dest data.NODE_ID, msgType data.MSG_TYPE, mid data.MESSAGE_ID) *data.MessageCore {
	return &data.MessageCore{
		Source:    n.Id,
		Dest:      dest,
		Leader:    n.Leader,
		Type:      msgType,
		MessageId: mid,
	}
}

func (n *Node) NewTermCore() *data.TermCore {
	return &data.TermCore{
		TermId: n.State.CurrentTerm,
	}
}

func (n *Node) Get(key data.KEY_TYPE) data.VAL_TYPE {
	return n.Data[key]
}

func (n *Node) Set(key data.KEY_TYPE, val data.VAL_TYPE) {
	n.Data[key] = val
}
