package stats

import (
	"math"
	"time"
)

type Stats struct {
	Duration    time.Duration
	Max         time.Duration
	Min         time.Duration
	P90         time.Duration
	P99         time.Duration
	ErrorCounts map[string]int
	StatusCount map[int]int
}

func Percentile(p int, times []time.Duration) time.Duration {
	rank := math.Round((float64(p) / 100.0) * float64(len(times)-1))
	return times[int(rank)]
}
func getAverageTime(totalTime, len float64) float64 {
	return totalTime / len
}
