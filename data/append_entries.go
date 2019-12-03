package data

type AppendCore struct {
	// NOTE: Term ID from TermCore
	// NOTE: Leader ID from MessageCore

	// Index of log entry immediately preceding new ones
	PrevLogIndex ENTRY_INDEX `json:"prevLogIndex"`

	// Term of prevLogIndex entry
	PrevLogTerm TERM_ID `json:"prevLogTerm"`

	// The log entries to store
	Entries []LogEntry `json:"entries"`

	// The Leader’s CommitIndex
	LeaderCommit ENTRY_INDEX `json:"leaderCommit"`
}

type ENTRY_INDEX uint32

type AppendEntries struct {
	*MessageCore
	*TermCore
	*AppendCore
}

type AppendEntriesResponseData struct {
	// Whether the log entry was successfully appended (and applied, if needed)
	Success bool `json:"success"`
}

type AppendEntriesResponse struct {
	*MessageCore
	*TermCore
	*AppendEntriesResponseData
}
