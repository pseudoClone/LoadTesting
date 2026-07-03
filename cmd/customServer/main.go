package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Downloadable body string")
		// log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
	})
	fmt.Println("Server running on :8000")
	http.ListenAndServe(":8000", nil)
}
