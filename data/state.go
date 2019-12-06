package data

type ENTRY_TYPE string

const (
	// Append Types
	PROMISE = ENTRY_TYPE("promise")
	COMMIT  = ENTRY_TYPE("commit")
)

type NODE_STATE int

const (
	FOLLOWER  = NODE_STATE(1)
	CANDIDATE = NODE_STATE(2)
	LEADER    = NODE_STATE(3)
)

type LogEntry struct {
	// EntryId ENTRY_INDEX `json:"id"`
	Key   KEY_TYPE   `json:"key"`
	Value VAL_TYPE   `json:"val"`
	Type  ENTRY_TYPE `json:"type"`
	Term  TERM_ID    `json:"term"` // Ignore in JSON
}

// RaftState is the state used for Consensus and log replication
type RaftState struct {
	// currentTerm latest term server has seen (initialized to 0
	// on first boot, increases monotonically)
	CurrentTerm TERM_ID

	// votedFor candidateId that received vote in current
	// term (or null if none)
	VotedFor NODE_ID

	// log[] log entries; each entry contains command
	// for state machine, and term when entry
	// was received by leader (first index is 1)
	Log []LogEntry

	// commitIndex index of highest log entry known to be
	// committed (initialized to 0, increases
	// monotonically)
	CommitIndex ENTRY_INDEX

	// lastApplied index of highest log entry applied to state
	// machine (initialized to 0, increases
	// monotonically)
	LastApplied ENTRY_INDEX

	// nextIndex[] for each server, index of the next log entry
	// to send to that server (initialized to leader
	// last log index + 1)
	// NOTE -- Reinitialized after election
	NextIndex []ENTRY_INDEX // LEADER ONLY STATE

	// matchIndex[] for each server, index of highest log entry
	// known to be replicated on server
	// (initialized to 0, increases monotonically
	// NOTE -- Reinitialized after election
	MatchIndex []ENTRY_INDEX // LEADER ONLY STATE
}

func NewRaftState(neighborCount int) RaftState {
	initialState := RaftState{
		CurrentTerm: 0,
		VotedFor:    "",                  // Stand-in for `null`
		Log:         make([]LogEntry, 1), // Index starts at `1`
		CommitIndex: 0,
		LastApplied: 0,
		NextIndex:   make([]ENTRY_INDEX, neighborCount, neighborCount), // Re-initialized on leader election
		MatchIndex:  make([]ENTRY_INDEX, neighborCount, neighborCount), // Re-initialized on leader election
	}
	return initialState
}

// Vote tracking

func (s *RaftState) SetVotedFor(candidate NODE_ID) {
	s.VotedFor = candidate
}

// Terms

func (s *RaftState) IncrementTerm() {
	s.CurrentTerm++
}

// Log Indices

func (s *RaftState) ResetLeaderIndices() {
	for idx, _ := range s.NextIndex {
		s.NextIndex[idx] = s.LastLogIndex() // Leader Last Log Index + 1
	}

	for idx, _ := range s.MatchIndex {
		s.MatchIndex[idx] = 0
	}
}

func (s *RaftState) LastLogIndex() ENTRY_INDEX {
	l := len(s.Log)
	if l > 0 {
		return ENTRY_INDEX(l - 1)
	} else {
		return 0 // TODO: Validate correct default
	}
}

func (s *RaftState) LastLogTerm() TERM_ID {
	l := len(s.Log)
	if l > 0 {
		return s.Log[l-1].Term
	} else {
		return s.CurrentTerm // TODO: Validate correct default
	}
}

// Log values

func (s *RaftState) AppendLog(entries ...LogEntry) {
	s.Log = append(s.Log, entries...)
}
