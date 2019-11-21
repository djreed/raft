package data

const (
  	// Consensus Algorithm Messages
  	APPEND = MSG_TYPE("append")
  	VOTE   = MSG_TYPE("vote")

  	// Append Types
  	PROMISE = MSG_TYPE("promise")
  	COMMIT  = MSG_TYPE("commit")
)

///
// Current Term
///

// Term the message corresponds to
type TERM_ID uint32

type TermCore struct {
	TermId TERM_ID `json:"term"`
}

///
// Log Append
///

/*
Invoked by leader to replicate log entries (§5.3); also used as
heartbeat (§5.2).

Arguments:
  term leader’s term
  leaderId so follower can redirect clients
  prevLogIndex index of log entry immediately preceding
        new ones
  prevLogTerm term of prevLogIndex entry
  entries[] log entries to store (empty for heartbeat;
        may send more than one for efficiency)
  leaderCommit leader’s commitIndex

Results:
  term currentTerm, for leader to update itself
  success true if follower contained entry matching
        prevLogIndex and prevLogTerm

Receiver implementation:
1. Reply false if term < currentTerm (§5.1)
2. Reply false if log doesn’t contain an entry at prevLogIndex
    whose term matches prevLogTerm (§5.3)
3. If an existing entry conflicts with a new one (same index
    but different terms), delete the existing entry and all that
    follow it (§5.3)
4. Append any new entries not already in the log
5. If leaderCommit > commitIndex, set commitIndex =
    min(leaderCommit, index of last new entry)
*/

type APPEND_ID uint32

type UpdateValue struct {
	AppendId APPEND_ID `json:"id"`
	Key      KEY_TYPE  `json:"key"`
	Value    VAL_TYPE  `json:"val"`
	Type     MSG_TYPE  `json:"type"`
}

type UpdatePayload struct {
	Updates []UpdateValue `json:"updates"`
}

type AppendMessage struct {
	*MessageCore
	*MessageIdBase
	*TermCore
	*UpdatePayload
}

///
// Leader Election Votes
///

/*
Invoked by candidates to gather votes (§5.2).

Arguments:
  term candidate’s term
  candidateId candidate requesting vote
  lastLogIndex index of candidate’s last log entry (§5.4)
  lastLogTerm term of candidate’s last log entry (§5.4)

Results:
  term currentTerm, for candidate to update itself
  voteGranted true means candidate received vote

Receiver implementation:
1. Reply false if term < currentTerm (§5.1)
2. If votedFor is null or candidateId, and candidate’s log is at
    least as up-to-date as receiver’s log, grant vote (§5.2, §5.4)

*/

type VoteMessageBody struct {
	VoteId NODE_ID `json:"vote"`
}

type VoteMessage struct {
	*MessageCore
	*MessageIdBase
	*TermCore
	*VoteMessageBody
}
