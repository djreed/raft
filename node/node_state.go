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

	n.ResetVotes()
	n.ResetReplications()

	n.State.IncrementTerm()
	n.VoteForSelf()

	n.ResetElectionTimeout()

	n.UnsetHeartbeatTimeout()
}

func (n *Node) BecomeLeader() {
	n.SetRole(data.LEADER)
	n.SetLeader(n.Id)

	n.ResetVotes()
	n.ResetReplications()

	n.State.CommitAll()

	n.UnsetElectionTimeout()

	n.ResetHeartbeatTimeout()

	n.State.ResetLeaderIndices()
}

func (n *Node) BeginCommit() {
	n.ResetVotes()
	n.ResetReplications()
}

func (n *Node) EndCommit() {
	n.ResetVotes()
	n.ResetReplications()
}

// Timeout reset

func NewElectionTimeout() <-chan time.Time {
	randomScale := time.Duration(time.Duration(rand.Int31n(electRange)) * time.Millisecond)
	return time.After(electBase + randomScale)
}

func NewHeartbeatTimeout() <-chan time.Time {
	return time.After(heartbeatBase)
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

// Commit promises

func (n *Node) ResetReplications() {
	n.Replications = 1
	n.ReplicatedNodes = make(map[data.NODE_ID]bool)
	n.ReplicationMessages = make(map[data.MESSAGE_ID]bool)
}

func (n *Node) IncrementReplications() {
	n.Replications += 1
}

func (n *Node) ReplicationQuorum() bool {
	return n.Replications > (len(n.Neighbors)+1)/2
}

func (n *Node) AlreadyReplicated(nid data.NODE_ID) bool {
	return n.ReplicatedNodes[nid]
}

func (n *Node) SetReplicated(nid data.NODE_ID) {
	n.ReplicatedNodes[nid] = true
}

func (n *Node) AddReplicationMid(mid data.MESSAGE_ID) {
	n.ReplicationMessages[mid] = true
}

func (n *Node) IsReplicationMessage(mid data.MESSAGE_ID) bool {
	isReplicationMessage, present := n.ReplicationMessages[mid]
	return isReplicationMessage && present
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

// Terms

func (n *Node) HandleTermUpdate(newTerm data.TERM_ID, leader data.NODE_ID) bool {
	shouldUpdate := n.State.CurrentTerm < newTerm
	if shouldUpdate {
		n.BecomeFollower(leader)
		n.State.SetTerm(newTerm)
	}
	return shouldUpdate
}

func (n *Node) TargetUpToDate(lastLogIndex data.ENTRY_INDEX, lastLogTerm data.TERM_ID) bool {
	// Not possible from state machine to have same index with differing term
	// TODO verify
	return (lastLogIndex >= n.State.LastLogIndex()) && (lastLogTerm >= n.State.LastLogTerm())
}

func (n *Node) PendingCommits() bool {
	return n.State.CommitIndex > n.State.LastApplied
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

func (n *Node) NeighborIndex(id data.NODE_ID) int {
	for idx, nid := range n.Neighbors {
		if nid == id {
			return idx
		}
	}
	return -69
}
