package main

import (
	"flag"
	"fmt"
	customclient "httpLoadTester/internal/customClient"
	"log"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	var wg sync.WaitGroup

	/* https://pkg.go.dev/net/http#hdr-Clients_and_Transports

	For control over proxies, TLS configuration, keep-alives, compression,
	and other settings, create a Transport:

	*/
	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxConnsPerHost:     1000,
		MaxIdleConnsPerHost: 1000,
	}

	serverURL := flag.String("s", "", "Enter server url")
	numWorkers := flag.Int("w", 3, "Enter the number of workers")
	numberOfRequests := flag.Int("n", 1,
		"Enter the number of concurrent clients")
	flag.Parse()
	parsedUrl, err := url.ParseRequestURI(*serverURL)
	client := http.Client{Timeout: 15 * time.Second, Transport: tr}

	if err != nil {
		log.Fatalf("Invalid URL %s", err)
	}

	fmt.Println(parsedUrl)
	jobsCh := make(chan string, *numberOfRequests)
	resultsCh := make(chan customclient.ReturnResult, *numberOfRequests)
	bar := progressbar.Default(int64(*numberOfRequests), "Fetching")

	for i := 1; i <= *numWorkers; i++ {
		wg.Add(1)
		go Worker(jobsCh, resultsCh, &wg, &client)
	}

	for i := 0; i < *numberOfRequests; i++ {
		jobsCh <- parsedUrl.String()
	}
	close(jobsCh)

	results := make([]customclient.Result, 0, *numberOfRequests)

	durationSlice := make([]time.Duration, 0)

	statusCounts := make(map[int]int)
	errorCounts := make(map[string]int)

	go func() {
		wg.Wait()
		close(resultsCh)
	}() // Define and run a background routine to close result channel

	for rr := range resultsCh {
		bar.Add(1)
		if rr.Err != nil {
			errorCounts[rr.Err.Error()]++
			continue
		}
		statusCounts[rr.StatusCode]++
		results = append(results, rr.Result)
		durationSlice = append(durationSlice, rr.Duration)
	}
	totalTime := 0.0
	totalBytes := 0
	for _, res := range results {
		totalTime += float64(res.Duration.Milliseconds())
		totalBytes += int(res.Bytes)
	}
	averageTime := totalTime / float64(len(results))
	fmt.Println("Average Time in milliseconds: ", averageTime, "ms")

	slices.Sort(durationSlice)
	for code, count := range statusCounts {
		fmt.Printf("===========================================\n"+
			"Status Code[%d]:\n%d responses\n"+
			"======================================"+
			"\n", code, count)
	}

	if len(errorCounts) > 0 {
		for err, count := range errorCounts {
			fmt.Printf("\n %d errors of type \n"+
				"====================================="+
				"\n%s\n ===========================\n",
				count, err)
		}
	}
	fmt.Println("Maximum Time Request in milliseconds: ",
		slices.Max(durationSlice))
	fmt.Println("Minimum Time Request in milliseconds: ",
		slices.Min(durationSlice))
	fmt.Println("p90: ", customclient.Percentile(90, durationSlice))
	fmt.Println("p99: ", customclient.Percentile(99, durationSlice))
	fmt.Println("Total bytes downloaded: ", totalBytes, " bytes")
}

func Worker(jobs <-chan string, results chan customclient.ReturnResult,
	wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()

	for url := range jobs {
		customclient.ClientRunner(url, client, results)
	}
}
