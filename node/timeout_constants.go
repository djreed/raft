package node

import "time"

const (
	electRange     = 200
	electBase      = time.Duration(300 * time.Millisecond)
	heartbeatScale = 20
)
