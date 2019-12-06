package data

type ENTRY_INDEX uint32

type AppendEntries struct {
	*MessageCore
	*TermCore

	// Index of log entry immediately preceding new ones
	PrevLogIndex ENTRY_INDEX `json:"prevLogIndex"`

	// Term of prevLogIndex entry
	PrevLogTerm TERM_ID `json:"prevLogTerm"`

	// The log entries to store
	Entries []LogEntry `json:"entries"`

	// The Leader’s CommitIndex
	LeaderCommit ENTRY_INDEX `json:"leaderCommit"`
}

type AppendEntriesResponse struct {
	*MessageCore
	*TermCore

	// Whether all of the entries here were appended/committed successfully
	Success bool `json:"success"`
}
