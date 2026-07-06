package config

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type HeaderFlags []string

func (h *HeaderFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *HeaderFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

type Config struct {
	NumWorkers        int
	NumRequests       int
	RequestsPerSecond int
	Request           RequestConfig
	/* I embedded this because so many fucking instance of
	like 1000 fucking structs
	*/
}

type RequestConfig struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

var headers HeaderFlags

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
