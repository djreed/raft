package data

type ENTRY_TYPE string

const (
	// Append Types
	PUT = ENTRY_TYPE("put")
)

type NODE_STATE int

const (
	FOLLOWER  = NODE_STATE(1)
	CANDIDATE = NODE_STATE(2)
	LEADER    = NODE_STATE(3)
)

type LogEntry struct {
	Key   KEY_TYPE   `json:"key"`
	Value VAL_TYPE   `json:"val"`
	Type  ENTRY_TYPE `json:"type"`
	Term  TERM_ID    `json:"term"`

	Sender NODE_ID    `json:"senderId"`
	MID    MESSAGE_ID `json:"MID"`
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
	// was received by leader (first index is 1) (TODO validate indices checking)
	Log []LogEntry

	// The actual key<>value store, built from the Log above
	Data map[KEY_TYPE]VAL_TYPE

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
		Log:         make([]LogEntry, 0), // Index starts at NOT 1, ITS 0, EAT SHIT AND DIE
		Data:        make(map[KEY_TYPE]VAL_TYPE),
		CommitIndex: -1,                                                // TODO Should be initialized to 0 if 1-indexed
		LastApplied: -1,                                                // TODO should be initialized to 0 if 1-indexed
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

func (s *RaftState) SetTerm(term TERM_ID) {
	s.CurrentTerm = term
}

// Log Indices

func (s *RaftState) ResetLeaderIndices() {
	for idx, _ := range s.NextIndex {
		s.NextIndex[idx] = s.LastLogIndex() + 1 // Leader Last Log Index + 1
	}

	for idx, _ := range s.MatchIndex {
		s.MatchIndex[idx] = -1 // TODO would be 0 if we're 1-indexed
	}
}

func (s *RaftState) LastLogIndex() ENTRY_INDEX {
	return ENTRY_INDEX(len(s.Log) - 1)
}

func (s *RaftState) LastLogTerm() TERM_ID {
	idx := s.LastLogIndex()
	if idx >= 0 {
		return s.Log[idx].Term
	} else {
		return s.CurrentTerm // TODO: Validate correct default
	}
}

// Log values

func (s *RaftState) GetLogEntry(idx ENTRY_INDEX) (entry LogEntry, found bool) {
	if 0 < idx && int(idx) < len(s.Log) {
		entry = s.Log[idx]
		found = true
	}
	return
}

func (s *RaftState) AppendLog(entries ...LogEntry) {
	s.Log = append(s.Log, entries...)
}

func (s *RaftState) CommitAll() {
	// TODO validate initial index assignment
	s.CommitTo(ENTRY_INDEX(len(s.Log) - 1))
}

func (s *RaftState) CommitTo(commitTo ENTRY_INDEX) {
	start := s.LastApplied
	for idx := start + 1; /* TODO validate the `+1` here */ idx <= commitTo; idx++ {
		s.ApplyEntry(s.Log[idx])
		s.LastApplied = ENTRY_INDEX(idx)
	}
}

func (s *RaftState) ApplyEntry(entry LogEntry) {
	s.Data[entry.Key] = entry.Value
}

func (s *RaftState) Get(key KEY_TYPE) VAL_TYPE {
	return s.Data[key]
}

// Neighbor index tracking

// Index of start of log entries to send
func (s *RaftState) IndexToSend(neighborIdx int) ENTRY_INDEX {
	return s.NextIndex[neighborIdx]
}

// Index of most up to date committed (replicated) log entry
func (s *RaftState) IndexReplicated(neighborIdx int) ENTRY_INDEX {
	return s.MatchIndex[neighborIdx]
}
