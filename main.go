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
	flag.Parse()
	parsedUrl, err := url.ParseRequestURI(*serverURL)
	/* ParseRequestURI validates URLs. Checked with:
	go run .\main.go -s kjddksjnfds
	go run .\main.go -s what
	Both of which return invalid URL.
	This mean, I don't have to use regex or validation myself*/
	if err != nil {
		log.Fatalf("Invalid URL %s", err)
	}
	fmt.Println(parsedUrl)
	start := time.Now()
	fmt.Println(start)
	if err := clientRunner(parsedUrl); err != nil {
		log.Fatal(err)
	}
	end := time.Now()
	fmt.Print("Elapsed: ", end.Sub(start), "\n")
	fmt.Println(time.Since(start))
}

func clientRunner(URLvar *url.URL) error {
	client := http.Client{Timeout: 15 * time.Second}
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
