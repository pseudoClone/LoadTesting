package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Result struct {
	StatusCode int
	Bytes      int
	Duration   time.Duration
}

func main() {
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
	results := make([]Result, 0, *numberOfConnections)
	/* slice of length 0 and capacity is the numberOfConnection */
	/* Couldn't find exact function signature,
	had to look up in stackoverflow
	https://stackoverflow.com/questions/36349045/how-can-the-make-function-take-three-parameters */
	fmt.Println(parsedUrl)
	for i := range *numberOfConnections {
		fmt.Println("Running Connection number", i+1)
		res, err := clientRunner(parsedUrl, &client)
		if err != nil {
			log.Println("Connection failed: ", i+1, "\n", err)
			continue
		}
		results = append(results, *res)
	}
	totalTime := 0.0
	totalBytes := 0
	for _, res := range results {
		totalTime += float64(res.Duration.Milliseconds())
		totalBytes += int(res.Bytes)
	}
	averageTime := totalTime / float64(*numberOfConnections)
	fmt.Println("Average Time in milliseconds: ", averageTime, "ms")
	fmt.Println("Total bytes downloaded: ", totalBytes, " bytes")
}

func clientRunner(URLvar *url.URL, client *http.Client) (*Result, error) {
	var res Result
	start := time.Now()
	resp, err := client.Get(URLvar.String())
	if err != nil {
		return &res, err
	}
	defer resp.Body.Close()
	res.StatusCode = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	duration := time.Since(start)
	res.Duration = duration

	res.Bytes = len(body)
	// fmt.Println("Total Bytes Read: ", len(body))
	// fmt.Printf("%s\n", body)

	return &res, nil
}
