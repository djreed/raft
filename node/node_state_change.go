package node

import (
	"math/rand"
	"time"

	"github.com/djreed/raft/data"
)

func BecomeFollower(n *Node) {
	// TODO
}

/*
  • On conversion to candidate, start election:
    • Increment currentTerm
    • Vote for self
    • Reset election timer
*/
func BecomeCandidate(n *Node) {
	ResetVotes(n)
	VoteForSelf(n)
	SetRole(n, data.CANDIDATE)
	SetLeader(n, data.UNKNOWN_LEADER)
	ResetElectionTimeout(n)
	// TODO potentially reset additional state
}

func BecomeLeader(n *Node) {
	// TODO
}

const (
	electBase = time.Duration(200 * time.Millisecond)
)

// 200-300ms
func NewElectionTimeout() <-chan time.Time {
	randomScale := time.Duration(time.Duration(rand.Int()) * time.Millisecond)
	return time.After(randomScale + electBase)
}

func NewHeartbeatTimeout() <-chan time.Time {
	return time.After(electBase / 20)
}

func ResetVotes(n *Node) {
	n.Votes = 0
}

func VoteForSelf(n *Node) {
	VoteFor(n, n.Id)
	n.Votes += 1
}

func SetRole(n *Node, role data.NODE_STATE) {
	n.Role = role
}

func SetLeader(n *Node, leader data.NODE_ID) {
	n.Leader = leader
}

func ResetElectionTimeout(n *Node) {
	n.ElectionTimeout = NewElectionTimeout()
}

func IncrementTerm(n *Node) {
	n.State.CurrentTerm++
}
