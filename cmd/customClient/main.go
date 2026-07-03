package main

import (
	"flag"
	"fmt"
	customclient "httpLoadTester/internal/customClient"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	resultsCh := make(chan customclient.ReturnResult)
	serverURL := flag.String("s", "", "Enter server url")
	numberOfConnections := flag.Int("n", 1,
		"Enter the number of concurrent clients")
	flag.Parse()
	parsedUrl, err := url.ParseRequestURI(*serverURL)
	client := http.Client{Timeout: 15 * time.Second}
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
	for i := range *numberOfConnections {
		fmt.Println("Running Connection number", i+1)
		go customclient.ClientRunner(parsedUrl.String(), &client, resultsCh)
		// if err != nil {
		// 	log.Println("Connection failed: ", i+1, "\n", err)
		// 	continue
		// }
		// results = append(results, *res)
	}
	for i := 0; i < *numberOfConnections; i++ {
		rr := <-resultsCh
		if rr.Err != nil {
			log.Println(rr.Err)
			continue
		}
		results = append(results, rr.Result)
	}
	totalTime := 0.0
	totalBytes := 0
	for _, res := range results {
		totalTime += float64(res.Duration.Milliseconds())
		totalBytes += int(res.Bytes)
	}
	averageTime := totalTime / float64(len(results))
	fmt.Println("Average Time in milliseconds: ", averageTime, "ms")
	fmt.Println("Total bytes downloaded: ", totalBytes, " bytes")
}
