package loadtest

import (
	"net/http"
	"sync"
	"time"

	"httpLoadTester/internal/config"
	"httpLoadTester/internal/httpclient"
	"httpLoadTester/internal/reporter"
	"httpLoadTester/internal/worker"

	"github.com/schollz/progressbar/v3"
)

func Run(cfg *config.Config) {
	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxConnsPerHost:     1000,
		MaxIdleConnsPerHost: 1000,
	}
	client := &http.Client{Timeout: 15 * time.Second, Transport: tr}

	jobsCh := make(chan string, cfg.NumRequests)
	resultsCh := make(chan httpclient.ReturnResult, cfg.NumRequests)
	bar := progressbar.Default(int64(cfg.NumRequests), "Fetching")

	var wg sync.WaitGroup
	for i := 1; i <= cfg.NumWorkers; i++ {
		wg.Add(1)
		go worker.Run(jobsCh, resultsCh, &wg, client)
	}

	for i := 0; i < cfg.NumRequests; i++ {
		jobsCh <- cfg.URL
	}
	close(jobsCh)

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	allResults := make([]httpclient.ReturnResult, 0, cfg.NumRequests)
	for rr := range resultsCh {
		bar.Add(1)
		allResults = append(allResults, rr)
	}

	report := reporter.Aggregate(allResults)
	reporter.Print(report)
}
