package worker

import (
	"net/http"
	"sync"

	"httpLoadTester/internal/config"
	"httpLoadTester/internal/httpClient"
)

func Run(
	jobs <-chan struct{},
	results chan<- httpClient.ReturnResult,
	wg *sync.WaitGroup,
	client *http.Client,
	req *config.RequestConfig,
) {
	defer wg.Done()
	for range jobs {
		results <- httpClient.DoRequest(req, client)
	}
}
