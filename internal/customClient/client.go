package customclient

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

func ClientRunner(URLvar string, client *http.Client, ch chan ReturnResult) {
	// fmt.Println("starting request")
	var res Result
	start := time.Now()
	req, err := http.NewRequest("GET", URLvar, nil)
	if err != nil {
		ch <- ReturnResult{Result{}, err}
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		ch <- ReturnResult{res, err}
		return
	}
	defer resp.Body.Close()
	res.StatusCode = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- ReturnResult{Result{}, err}
		return
	}
	duration := time.Since(start)
	res.Duration = duration

	res.Bytes = len(body)
	// fmt.Println("Total Bytes Read: ", len(body))
	// fmt.Printf("%s\n", body)
	// fmt.Println("sending result")
	ch <- ReturnResult{res, nil}
}
