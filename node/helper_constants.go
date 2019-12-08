package node

import "time"

const (
	electRange    = 200
	electBase     = time.Duration(300 * time.Millisecond)
	heartbeatBase = electBase / 5
)

func Min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
