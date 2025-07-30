package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"
)

// Creates an HTTP server that never responds to a client.
func newHangServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Minute) //Simulates a long hang for the server
		},
	))
}

// Creates an HTTP server that sends an HTML page with an unstable DOM to a client.
func newDOMUnstableServer() *httptest.Server {
	//Path to the HTML file to serve
	filePath := filepath.Join("data", "unstable.html")

	//Create a handler for the HTML file
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := os.ReadFile(filePath)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "text/html")
			w.Write(data)
		},
	))
}
