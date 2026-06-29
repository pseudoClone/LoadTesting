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
	elapsedTimes := make([]time.Duration, 0, *numberOfConnections)
	/* slice of length 0 and capacity is the numberOfConnection */
	/* Couldn't find exact function signature,
	had to look up in stackoverflow
	https://stackoverflow.com/questions/36349045/how-can-the-make-function-take-three-parameters */
	fmt.Println(parsedUrl)
	for i := range *numberOfConnections {
		fmt.Println("Running Connection number", i+1)
		start := time.Now()
		fmt.Println(start)
		if err := clientRunner(parsedUrl, &client); err != nil {
			log.Fatal(err)
		}
		// end := time.Now()
		// fmt.Print("Elapsed: ", end.Sub(start), "\n")
		elapsedTime := time.Since(start)
		elapsedTimes = append(elapsedTimes, elapsedTime)
	}
	totalTime := 0.0
	for _, valTime := range elapsedTimes {
		totalTime += float64(valTime.Milliseconds())
	}
	averageTime := totalTime / float64(*numberOfConnections)
	fmt.Println("Average Time in milliseconds: ", averageTime, "ms")
}

func clientRunner(URLvar *url.URL, client *http.Client) error {
	resp, err := client.Get(URLvar.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Total Bytes Read: ", len(body))
		fmt.Printf("%s\n", body)
	}
	// fmt.Println(resp, err)
	/* if err != nil {
		fmt.Fprintf(os.Stderr, "\nCould not get the requested URL\n")
		return
	}

	// print(resp.Request.Host)
	if resp.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading HTTP response body")
			return
		}
		bodyString := string(bodyBytes)
		log.Print(bodyString)
	} */
	return nil
}

/*


type Result struct {
	StatusCode int
	Bytes int
	Duration time.Duration
}

*/
