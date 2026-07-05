package config

import (
	"flag"
	"log"
	"net/url"
)

type Config struct {
	URL         string
	NumWorkers  int
	NumRequests int
}

func Load() *Config {
	serverURL := flag.String("s", "", "Enter server url")
	numWorkers := flag.Int("w", 3, "Enter the number of workers")
	numberOfRequests := flag.Int("n", 1,
		"Enter the number of concurrent clients")
	flag.Parse()

	if *serverURL == "" {
		log.Fatalf("Server URL is required")
	}

	_, err := url.ParseRequestURI(*serverURL)
	if err != nil {
		log.Fatalf("Invalid URL %s", err)
	}

	return &Config{
		URL:         *serverURL,
		NumWorkers:  *numWorkers,
		NumRequests: *numberOfRequests,
	}
}
