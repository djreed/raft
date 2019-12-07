package data

type ENTRY_INDEX int32

type AppendEntries struct {
	*MessageCore
	*TermCore
	PrevLogIndex ENTRY_INDEX `json:"prevLogIndex"` // Index of log entry immediately preceding new ones
	PrevLogTerm  TERM_ID     `json:"prevLogTerm"`  // Term of prevLogIndex entry
	Entries      []LogEntry  `json:"entries"`      // The log entries to store
	LeaderCommit ENTRY_INDEX `json:"leaderCommit"` // The Leader’s CommitIndex
}

type AppendEntriesResponse struct {
	*MessageCore
	*TermCore

	// Whether all of the entries here were appended/committed successfully
	Success bool `json:"success"`
}
