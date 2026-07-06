package config

import (
	"flag"
	"fmt"
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

	flag.Var(&headers, "H",
		"Enter HTTP Header as -H \"Content-Type: application-json \" ")

	body := flag.String("d", "", "Enter request body")
	bodyFile := flag.String("body", "",
		"Enter the filename of file containing request body")

	method := flag.String("m", http.MethodGet, "Enter http method")
	flag.Parse()
	m := strings.ToUpper(*method)

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

	headerMap, err := parseHeaders(headers)
	if err != nil {
		log.Fatal(err)
	}
	var bodyBytes []byte

	switch {
	case *body != "" && *bodyFile != "":
		log.Fatal("Cannot use both -d and -body flag")
	case *body != "":
		bodyBytes = []byte(*body)
	case *bodyFile != "":
		data, err := os.ReadFile(*bodyFile)
		if err != nil {
			log.Fatal(err)
		}
		bodyBytes = data
	}
	validMethods := map[string]struct{}{
		http.MethodGet:     {},
		http.MethodPost:    {},
		http.MethodPut:     {},
		http.MethodPatch:   {},
		http.MethodDelete:  {},
		http.MethodHead:    {},
		http.MethodOptions: {},
	}

	if _, ok := validMethods[m]; !ok {
		log.Fatalf("unsupported method %q", m)
	}

	return &Config{
		NumWorkers:        *numWorkers,
		NumRequests:       *numberOfRequests,
		RequestsPerSecond: *rps,
		Request: RequestConfig{
			Method:  m,
			URL:     *serverURL,
			Headers: headerMap,
			Body:    bodyBytes,
		},
	}, tr
}

func parseHeaders(values []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, h := range values {
		parts := strings.SplitN(h, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("Invalid Header %q", h)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}
	return headers, nil
}
