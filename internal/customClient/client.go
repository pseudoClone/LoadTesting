package customclient

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Result struct {
	StatusCode int
	Bytes      int
	Duration   time.Duration
}

func ClientRunner(URLvar string, client *http.Client) (*Result, error) {
	var res Result
	start := time.Now()
	req, _ := http.NewRequest("GET", URLvar, nil)
	resp, err := client.Do(req)
	if err != nil {
		return &res, err
	}
	defer resp.Body.Close()
	res.StatusCode = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return &Result{}, err
	}
	duration := time.Since(start)
	res.Duration = duration

	res.Bytes = len(body)
	// fmt.Println("Total Bytes Read: ", len(body))
	// fmt.Printf("%s\n", body)

	return &res, nil
}
