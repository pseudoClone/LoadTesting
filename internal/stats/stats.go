package stats

import (
	"math"
	"time"
)

func Percentile(p int, times []time.Duration) time.Duration {
	if len(times) == 0 {
		return 0
	}

	rank := int(math.Round((float64(p) / 100.0) * float64(len(times)-1)))

	if rank < 0 {
		rank = 0
	}
	if rank >= len(times) {
		rank = len(times) - 1
	}

	return times[rank]
}
