package httpClient

import (
	"bytes"
	"httpLoadTester/internal/config"
	"io"
	"net/http"
	"time"
)

type Result struct {
	StatusCode int
	Bytes      int64
	Duration   time.Duration
}

type ReturnResult struct {
	Result
	Err error
}

func DoRequest(cfg *config.RequestConfig, client *http.Client) ReturnResult {
	var res Result
	start := time.Now()
	req, err := http.NewRequest(
		cfg.Method,
		cfg.URL,
		bytes.NewReader(cfg.Body),
	)
	if err != nil {
		return ReturnResult{Err: err}
	}

	for k, v := range cfg.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return ReturnResult{Err: err}
	}
	defer resp.Body.Close()

	bytesRead, err := io.Copy(io.Discard, resp.Body)
	if err != nil {
		return ReturnResult{Err: err}
	}

	res.StatusCode = resp.StatusCode
	res.Duration = time.Since(start)
	res.Bytes = bytesRead

	return ReturnResult{Result: res, Err: nil}
}
