package node

import (
	"net"
	"time"

	"github.com/djreed/raft/data"
	"github.com/djreed/raft/logging"
)

var LOG = logging.LOG

const (
	CHAN_BUFFER = 128
)

// Node contains Network state
type Node struct {
	Id data.NODE_ID

	Neighbors []data.NODE_ID

	// Single UNIX domain socket, emulating Ethernet
	Socket *net.UnixConn

	Data      map[data.KEY_TYPE]data.VAL_TYPE
	RaftState data.RaftState

	// A Channel for each data type we need to handle
	RequestVotes  chan data.RequestVote
	AppendEntries chan data.AppendEntries
	GetMessages   chan data.GetMessage
	PutMessages   chan data.PutMessage

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
		RequestVotes:  make(chan data.RequestVote, CHAN_BUFFER),
		AppendEntries: make(chan data.AppendEntries, CHAN_BUFFER),
		GetMessages:   make(chan data.GetMessage, CHAN_BUFFER),
		PutMessages:   make(chan data.PutMessage, CHAN_BUFFER),
		RaftState:     initialRaftState,
	}

	go initialNode.HandleConn(unixSock)

	return initialNode
}

func OpenSocket(id data.NODE_ID) *net.UnixConn {
	unixAddr, err := net.ResolveUnixAddr("unixpacket", string(id))
	if err != nil {
		LOG.Panic(err)
	}

	conn, err := net.DialUnix("unixpacket", unixAddr, unixAddr)
	if err != nil {
		LOG.Panic(err)
	}

	return conn
}

func (n *Node) HandleConn(sock *net.UnixConn) error {
	// Read from Connection
	// Decode JSON into correct message type
	// Send along corresponding channel

	return nil
}
