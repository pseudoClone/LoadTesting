package worker

import (
	"net/http"
	"sync"

	"httpLoadTester/internal/httpClient"
)

func Run(jobs <-chan string, results chan<- httpClient.ReturnResult,
	wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()
	for url := range jobs {
		results <- httpClient.DoRequest(url, client)
	}
}
