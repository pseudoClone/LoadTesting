package httpclient

import (
	"io"
	"net/http"
	"time"
)

type Result struct {
	StatusCode int
	Bytes      int
	Duration   time.Duration
}

type ReturnResult struct {
	Result
	Err error
}

func DoRequest(url string, client *http.Client) ReturnResult {
	var res Result
	start := time.Now()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ReturnResult{Err: err}
	}
	resp, err := client.Do(req)
	if err != nil {
		return ReturnResult{Err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ReturnResult{Err: err}
	}

	res.StatusCode = resp.StatusCode
	res.Duration = time.Since(start)
	res.Bytes = len(body)

	return ReturnResult{Result: res, Err: nil}
}
