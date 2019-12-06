package node

import (
	"math/rand"
	"time"

	"github.com/djreed/raft/data"
)

func (n *Node) BecomeFollower(leader data.NODE_ID) {
	n.SetRole(data.FOLLOWER)
	n.SetLeader(leader)
	n.ResetVotes()
	n.ResetElectionTimeout()
	n.UnsetHeartbeatTimeout()
}

func (n *Node) BecomeCandidate() {
	n.SetRole(data.CANDIDATE)
	n.SetLeader(data.UNKNOWN_LEADER)
	n.State.IncrementTerm()
	n.ResetVotes()
	n.VoteForSelf()
	n.ResetElectionTimeout()
	n.UnsetHeartbeatTimeout()
}

func (n *Node) BecomeLeader() {
	n.SetRole(data.LEADER)
	n.ResetVotes()
	n.SetLeader(n.Id)
	n.UnsetElectionTimeout()
	n.ResetHeartbeatTimeout()
}

// Timeout reset

func NewElectionTimeout() <-chan time.Time {
	randomScale := time.Duration(time.Duration(rand.Int31n(electRange)) * time.Millisecond)
	return time.After(electBase + randomScale)
}

func NewHeartbeatTimeout() <-chan time.Time {
	return time.After(electBase / heartbeatScale)
}

// Votes

func (n *Node) ResetVotes() {
	n.Votes = 0
}

func (n *Node) VoteForSelf() {
	n.State.SetVotedFor(n.Id)
	n.IncrementVotes()
}

func (n *Node) IncrementVotes() {
	n.Votes += 1
}

func (n *Node) VoteQuorum() bool {
	return n.Votes > (len(n.Neighbors)+1)/2
}

// Roles

func (n *Node) SetRole(role data.NODE_STATE) {
	n.Role = role
}

// Leader

func (n *Node) SetLeader(leader data.NODE_ID) {
	n.Leader = leader
}

func (n *Node) IsLeader() bool {
	return n.Id == n.Leader
}

// Timeout Resets

func (n *Node) ResetElectionTimeout() {
	n.ElectionTimeout = NewElectionTimeout()
}

func (n *Node) UnsetElectionTimeout() {
	n.ElectionTimeout = nil
}

func (n *Node) ResetHeartbeatTimeout() {
	n.HeartbeatTimeout = NewHeartbeatTimeout()
}

func (n *Node) UnsetHeartbeatTimeout() {
	n.HeartbeatTimeout = nil
}
