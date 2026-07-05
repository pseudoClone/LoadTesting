package reporter

import (
	"fmt"
	"slices"
	"time"

	"httpLoadTester/internal/httpClient"
	"httpLoadTester/internal/stats"
)

type Report struct {
	Results      []httpClient.Result
	StatusCounts map[int]int
	ErrorCounts  map[string]int
	TotalBytes   int
	TotalTimeMs  float64
}

func Aggregate(results []httpClient.ReturnResult) *Report {
	rep := &Report{
		Results:      make([]httpClient.Result, 0, len(results)),
		StatusCounts: make(map[int]int),
		ErrorCounts:  make(map[string]int),
	}

	for _, rr := range results {
		if rr.Err != nil {
			rep.ErrorCounts[rr.Err.Error()]++
			continue
		}
		rep.StatusCounts[rr.StatusCode]++
		rep.Results = append(rep.Results, rr.Result)
		rep.TotalBytes += rr.Bytes
		rep.TotalTimeMs += float64(rr.Duration.Milliseconds())
	}

	return rep
}

func Print(rep *Report) {
	durations := make([]time.Duration, 0, len(rep.Results))
	for _, r := range rep.Results {
		durations = append(durations, r.Duration)
	}
	slices.Sort(durations)

	if len(rep.Results) > 0 {
		avgTime := rep.TotalTimeMs / float64(len(rep.Results))
		fmt.Printf("Average Time in milliseconds: %.2f ms\n", avgTime)
	}

	for code, count := range rep.StatusCounts {
		fmt.Printf("==========================================="+
			"\nStatus Code[%d]:\n%d responses\n==="+
			"========================================\n",
			code, count)
	}

	if len(rep.ErrorCounts) > 0 {
		for err, count := range rep.ErrorCounts {
			fmt.Printf("\n %d errors of type \n"+
				"=====================================\n%s\n"+
				"===========================\n", count, err)
		}
	}

	if len(durations) > 0 {
		fmt.Println("Maximum Time Request in milliseconds: ",
			slices.Max(durations))
		fmt.Println("Minimum Time Request in milliseconds: ",
			slices.Min(durations))
		fmt.Println("p90: ", stats.Percentile(90, durations))
		fmt.Println("p99: ", stats.Percentile(99, durations))
	}
	fmt.Println("Total bytes downloaded: ", rep.TotalBytes, " bytes")
}
