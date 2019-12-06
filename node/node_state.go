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
	n.IncrementTerm()
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

const (
	electBase = time.Duration(150 * time.Millisecond)
)

// 150-250ms
func NewElectionTimeout() <-chan time.Time {
	randomScale := time.Duration(time.Duration(rand.Int()) * time.Millisecond)
	return time.After(randomScale + electBase)
}

func NewHeartbeatTimeout() <-chan time.Time {
	return time.After(electBase / 10)
}

// Votes

func (n *Node) ResetVotes() {
	n.Votes = 0
}

func (n *Node) VoteFor(candidate data.NODE_ID) {
	n.State.VotedFor = candidate
}

func (n *Node) VoteForSelf() {
	n.VoteFor(n.Id)
	n.IncrementVotes()
}

func (n *Node) IncrementVotes() {
	n.Votes++
}

func (n *Node) VoteQuorum() bool {
	return n.Votes > ((len(n.Neighbors) + 1) / 2)
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

// Terms

func (n *Node) IncrementTerm() {
	n.State.CurrentTerm++
}

// Log Indices

func (n *Node) ResetLeaderIndices() {
	for idx, _ := range n.State.NextIndex {
		n.State.NextIndex[idx] = n.LastLogIndex() // Leader Last Log Index + 1
	}

	for idx, _ := range n.State.MatchIndex {
		n.State.MatchIndex[idx] = 0
	}
}

func (n *Node) LastLogIndex() data.ENTRY_INDEX {
	l := len(n.State.Log)
	if l > 0 {
		return data.ENTRY_INDEX(l - 1)
	} else {
		return 0 // TODO: Validate correct default
	}
}

func (n *Node) LastLogTerm() data.TERM_ID {
	l := len(n.State.Log)
	if l > 0 {
		return n.State.Log[l-1].Term
	} else {
		return n.State.CurrentTerm // TODO: Validate correct default
	}
}

// Log values

func (n *Node) AppendLog(entries ...data.LogEntry) {
	n.State.Log = append(n.State.Log, entries...)
}
