package test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-rod/rod"

	rutil "github.com/jgilman1337/rod_util/pkg"
)

// Creates an HTTP server that never responds to a client.
func newHangServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Minute) //Simulates a long hang for the server
		},
	))
}

// Creates an HTTP server that sends an HTML page to a client.
func newServer(file string) *httptest.Server {
	//Path to the HTML file to serve
	filePath := filepath.Join("data", file)

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

// Runs common tests related to static HTML pages.
func htmlRunner(t *testing.T, htmlPath string, debug bool, timeout int, runner func(b *rod.Browser, p *rod.Page)) {
	//Force-kills the httptest.server
	srv := newServer(htmlPath)
	defer func() {
		srv.Config.Close()
		srv.Listener.Close()
	}()

	//Launch Rod
	opts := rutil.DefaultBrowserOpts()
	if debug {
		opts = rutil.DefaultBrowserOptsDbg()
	}
	browser, launcher, err := rutil.BuildSandboxless(opts)
	if err != nil {
		t.Fatalf("Failed to launch browser: %s", err)
	}
	defer rutil.RodFree(browser, launcher)

	//Setup a page, with timeout; applies to the entire process, not just `page.Navigate`
	page := rutil.BlankPage(browser).Timeout(time.Duration(timeout) * time.Second)
	defer page.Close()
	if err := page.Navigate(srv.URL); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.Log("Successfully stopped page load early")
		} else {
			t.Fatalf("Failed to create a webpage: %s", err)
		}
	}

	//Wait for the page to load
	t.Log("Done creating page; waiting on DOM")
	if err := page.WaitDOMStable(time.Second, 0); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.Log("Successfully stopped DOM wait early")
		} else {
			t.Fatalf("WaitDOMStable failed: %s", err)
		}
	}

	//Run the tests
	runner(browser, page)
}
