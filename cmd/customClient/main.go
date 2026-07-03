package main

import (
	"flag"
	"fmt"
	customclient "httpLoadTester/internal/customClient"
	"log"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {

	/* https://pkg.go.dev/net/http#hdr-Clients_and_Transports

	For control over proxies, TLS configuration, keep-alives, compression,
	and other settings, create a Transport:

	*/
	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxConnsPerHost:     1000,
		MaxIdleConnsPerHost: 1000,
	}

	resultsCh := make(chan customclient.ReturnResult)
	serverURL := flag.String("s", "", "Enter server url")
	numberOfConnections := flag.Int("n", 1,
		"Enter the number of concurrent clients")
	flag.Parse()
	parsedUrl, err := url.ParseRequestURI(*serverURL)
	client := http.Client{Timeout: 15 * time.Second, Transport: tr}
	/* ParseRequestURI validates URLs. Checked with:
	go run .\main.go -s kjddksjnfds
	go run .\main.go -s what
	Both of which return invalid URL.
	This mean, I don't have to use regex or validation myself*/
	if err != nil {
		log.Fatalf("Invalid URL %s", err)
	}
	results := make([]customclient.Result, 0, *numberOfConnections)
	/* slice of length 0 and capacity is the numberOfConnection */
	/* Couldn't find exact function signature,
	had to look up in stackoverflow
	https://stackoverflow.com/questions/36349045/how-can-the-make-function-take-three-parameters */
	fmt.Println(parsedUrl)
	bar := progressbar.Default(int64(*numberOfConnections), "Fetching")
	for range *numberOfConnections {
		// fmt.Println("Running Connection number", i+1)
		go customclient.ClientRunner(parsedUrl.String(), &client, resultsCh)
	}
	durationSlice := make([]time.Duration, 0)
	for i := 0; i < *numberOfConnections; i++ {
		rr := <-resultsCh
		bar.Add(1)
		if rr.Err != nil {
			log.Println(rr.Err)
			continue
		}
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
	// for _, x := range durationSlice {
	// 	fmt.Println(x)
	// }
	slices.Sort(durationSlice)
	fmt.Println("Maximum Time Request in milliseconds: ", slices.Max(durationSlice))
	fmt.Println("Minimum Time Request in milliseconds: ", slices.Min(durationSlice))
	fmt.Println("p90: ", customclient.Percentile(90, durationSlice))
	fmt.Println("p99: ", customclient.Percentile(99, durationSlice))
	fmt.Println("Total bytes downloaded: ", totalBytes, " bytes")
}
