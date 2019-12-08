package node

import "time"

const (
	electRange    = 250
	electBase     = time.Duration(500 * time.Millisecond)
	heartbeatBase = electBase / 10
)

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
