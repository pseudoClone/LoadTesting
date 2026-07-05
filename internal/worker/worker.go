package worker

import (
	"net/http"
	"sync"

	"httpLoadTester/internal/httpclient"
)

func Run(jobs <-chan string, results chan<- httpclient.ReturnResult,
	wg *sync.WaitGroup, client *http.Client) {
	defer wg.Done()
	for url := range jobs {
		results <- httpclient.DoRequest(url, client)
	}
}
