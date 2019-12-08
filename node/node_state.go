package node

import (
	"math/rand"
	"time"

	"github.com/djreed/raft/data"
)

func (n *Node) BecomeFollower(leader data.NODE_ID) {
	n.SetRole(data.FOLLOWER)
	n.SetLeader(leader)

	n.ResetQuorum()

	n.ResetElectionTimeout()

	n.UnsetHeartbeatTimeout()

	ERR.Printf("(%v) BECAME FOLLOWER", n.Id)
}

func (n *Node) BecomeCandidate() {
	n.SetRole(data.CANDIDATE)
	n.SetLeader(data.UNKNOWN_LEADER)

	n.ResetQuorum()

	n.State.IncrementTerm()
	n.VoteForSelf()

	n.ResetElectionTimeout()

	n.UnsetHeartbeatTimeout()
	ERR.Printf("(%v) BECAME CANDIDATE", n.Id)
}

func (n *Node) BecomeLeader() {
	n.SetRole(data.LEADER)
	n.SetLeader(n.Id)

	n.ResetQuorum()

	n.State.CommitAll()

	n.UnsetElectionTimeout()

	n.ResetHeartbeatTimeout()

	n.State.ResetLeaderIndices()
	ERR.Printf("(%v) BECAME LEADER", n.Id)
}

func (n *Node) BeginCommit() {
	n.ResetQuorum()

}

func (n *Node) EndCommit() {
	n.ResetQuorum()
}

func (n *Node) ResetQuorum() {
	n.ResetVotes()
	n.ResetReplications()
	n.State.ResetVotedFor()
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
	n.Votes = 1
}

func (n *Node) VoteForSelf() {
	n.State.SetVotedFor(n.Id)
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

func (n *Node) RecordAppendMessage(msg data.AppendEntries) {
	n.AppendMessages[msg.MessageId] = msg
}

func (n *Node) DeleteAppendMessage(mid data.MESSAGE_ID) {
	delete(n.AppendMessages, mid)
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
	return (lastLogIndex >= n.State.LastLogIndex()) && (lastLogTerm >= n.State.LastLogTerm())
}

func (n *Node) PendingCommits() bool {
	return n.State.CommitIndex > n.State.LastApplied
}

// Timeout Resets

func (n *Node) ResetElectionTimeout() {
	if n.IsLeader() {
		ERR.Printf("(%v)\n\n\n\n\n\n\n\n\n\n\n_STOP_\n\n\n\n\n\n\n\n\n\n\n", n.Id)
		return
	}
	// ERR.Printf("(%v) RESET ELECTION TIMEOUT", n.Id)
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
	for neighborIdx, nid := range n.Neighbors {
		if nid == id {
			return neighborIdx
		}
	}
	return -69
}

func (n *Node) LastReplicatedIdx(nid data.NODE_ID) data.ENTRY_INDEX {
	// Which index is known to be replicated
	nodeIdx := n.NeighborIndex(nid)
	return n.State.IndexReplicated(nodeIdx)
}

func (n *Node) SendStartIdx(nid data.NODE_ID) data.ENTRY_INDEX {
	// Which index did we start sending at
	nodeIdx := n.NeighborIndex(nid)
	return n.State.IndexToSend(nodeIdx)
}
