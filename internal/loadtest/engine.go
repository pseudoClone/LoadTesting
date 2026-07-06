package loadtest

import (
	"net/http"
	"sync"
	"time"

	"httpLoadTester/internal/config"
	"httpLoadTester/internal/httpClient"
	"httpLoadTester/internal/reporter"
	"httpLoadTester/internal/worker"

	"github.com/schollz/progressbar/v3"
)

func Run(cfg *config.Config, tr *http.Transport) {
	client := &http.Client{Timeout: 15 * time.Second, Transport: tr}

	jobsCh := make(chan struct{}, cfg.NumWorkers)
	resultsCh := make(chan httpClient.ReturnResult, cfg.NumWorkers)
	bar := progressbar.Default(int64(cfg.NumRequests), "Fetching")

	var wg sync.WaitGroup
	for i := 1; i <= cfg.NumWorkers; i++ {
		wg.Add(1)
		go worker.Run(jobsCh, resultsCh, &wg, client, &cfg.Request)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	go func() {
		if cfg.RequestsPerSecond > 0 {
			interval := time.Second /
				time.Duration(cfg.RequestsPerSecond)
			ticker := time.NewTicker(interval)
			defer ticker.Stop()

			for i := 0; i < cfg.NumRequests; i++ {
				<-ticker.C
				jobsCh <- struct{}{}
			}
		} else {
			for i := 0; i < cfg.NumRequests; i++ {
				jobsCh <- struct{}{}
			}
		}
		close(jobsCh)
	}()

	allResults := make([]httpClient.ReturnResult, 0, cfg.NumRequests)
	for rr := range resultsCh {
		bar.Add(1)
		allResults = append(allResults, rr)
	}

	report := reporter.Aggregate(allResults)
	reporter.Print(report)
}
