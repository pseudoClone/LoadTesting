package config

import (
	"flag"
	"log"
	"net/http"
	"net/url"
)

type Config struct {
	URL               string
	NumWorkers        int
	NumRequests       int
	RequestsPerSecond int
}

func Load() (*Config, *http.Transport) {
	serverURL := flag.String("s", "", "Enter server url")
	numWorkers := flag.Int("w", 3, "Enter the number of workers")
	numberOfRequests := flag.Int("n", 1,
		"Enter the number of concurrent clients")
	rps := flag.Int("rps", 0, "Enter requests per second")
	flag.Parse()

	if *serverURL == "" {
		log.Fatalf("Server URL is required")
	}

	_, err := url.ParseRequestURI(*serverURL)
	if err != nil {
		log.Fatalf("Invalid URL %s", err)
	}
	tr := &http.Transport{
		MaxIdleConns:        1000,
		MaxConnsPerHost:     1000,
		MaxIdleConnsPerHost: 1000,
	}

	return &Config{
		URL:               *serverURL,
		NumWorkers:        *numWorkers,
		NumRequests:       *numberOfRequests,
		RequestsPerSecond: *rps,
	}, tr
}
