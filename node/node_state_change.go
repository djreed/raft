package node

import (
	"math/rand"
	"time"

	"github.com/djreed/raft/data"
)

func (n *Node) BecomeFollower(leader data.NODE_ID) {
	n.SetRole(data.FOLLOWER)
	n.ResetVotes()
	n.SetLeader(data.UNKNOWN_LEADER)
	n.ResetElectionTimeout()
}

/*
  • On conversion to candidate, start election:
    • Increment currentTerm
    • Vote for self
    • Reset election timer
*/
func (n *Node) BecomeCandidate() {
	n.SetRole(data.CANDIDATE)
	n.ResetVotes()
	n.VoteForSelf()
	n.SetLeader(data.UNKNOWN_LEADER)
	n.ResetElectionTimeout()
	// TODO potentially reset additional state
}

func (n *Node) BecomeLeader() {
	n.ResetVotes()
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

func (n *Node) ResetVotes() {
	n.Votes = 0
}

func (n *Node) VoteForSelf() {
	VoteFor(n, n.Id)
	n.Votes += 1
}

func (n *Node) SetRole(role data.NODE_STATE) {
	n.Role = role
}

func (n *Node) SetLeader(leader data.NODE_ID) {
	n.Leader = leader
}

func (n *Node) ResetElectionTimeout() {
	n.ElectionTimeout = NewElectionTimeout()
}

func (n *Node) IncrementTerm() {
	n.State.CurrentTerm++
}
