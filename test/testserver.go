package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><head><title>Test Server</title></head><body>")
		fmt.Fprintf(w, "<h1>Welcome to Test Server</h1>")
		fmt.Fprintf(w, `<script src="/static/app.js"></script>`)
		fmt.Fprintf(w, "</body></html>")
	})

	http.HandleFunc("/static/app.js", func(w http.ResponseWriter, r *http.Request) {
		js := `
			const apiURL = "/api/v1/data";
			const secretKey = "demo_secret";
		`
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, js)
	})

	http.HandleFunc("/api/v1/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"success","data":{"id":1,"value":"test"}}`)
	})

	fmt.Println("Starting test HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
