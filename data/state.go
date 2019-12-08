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
	// was received by leader (first index is 1)
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
		VotedFor:    "", // Stand-in for `null`
		Log:         make([]LogEntry, 0),
		Data:        make(map[KEY_TYPE]VAL_TYPE),
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

func (s *RaftState) SetTerm(term TERM_ID) {
	s.CurrentTerm = term
}

// Log Indices

func (s *RaftState) ResetLeaderIndices() {
	lastLogIdx := s.LastLogIndex()
	for idx, _ := range s.NextIndex {
		s.NextIndex[idx] = lastLogIdx + 1 // Leader Last Log Index + 1
	}

	for idx, _ := range s.MatchIndex {
		s.MatchIndex[idx] = 0
	}
}

func (s *RaftState) LastLogIndex() ENTRY_INDEX {
	return ENTRY_INDEX(len(s.Log))
}

func (s *RaftState) LastLogTerm() TERM_ID {
	idx := s.LastLogIndex()
	if idx > 0 {
		return s.Log[idx-1].Term
	} else {
		return 0
	}
}

// Log values

func (s *RaftState) GetLogEntry(idx ENTRY_INDEX) (entry LogEntry, found bool) {
	if 0 < idx && int(idx) <= len(s.Log) {
		entry = s.Log[idx-1]
		found = true
	}
	return
}

func (s *RaftState) AppendLog(entries ...LogEntry) {
	s.Log = append(s.Log, entries...)
}

func (s *RaftState) CommitAll() {
	lastIndex := s.LastLogIndex()
	s.CommitTo(lastIndex)
}

func (s *RaftState) CommitTo(commitTo ENTRY_INDEX) {
	s.CommitIndex = commitTo
}

func (s *RaftState) ApplyAll() {
	lastIndex := s.LastLogIndex()
	s.ApplyTo(lastIndex)
}

func (s *RaftState) ApplyTo(applyTo ENTRY_INDEX) {
	for idx := s.LastApplied + 1; idx <= applyTo; idx++ {
		s.LastApplied = ENTRY_INDEX(idx)
		s.ApplyEntry(s.Log[idx-1])
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
func (s *RaftState) IndexToSend(nodeIdx int) ENTRY_INDEX {
	return s.NextIndex[nodeIdx]
}

// Index of most up to date committed (replicated) log entry
func (s *RaftState) IndexReplicated(nodeIdx int) ENTRY_INDEX {
	return s.MatchIndex[nodeIdx]
}

// Voting

func (s *RaftState) Voted() bool {
	return s.VotedFor != ""
}

func (s *RaftState) VoteCandidate() NODE_ID {
	return s.VotedFor
}

func (s *RaftState) ResetVotedFor() {
	s.VotedFor = ""
}
