package pkg

import (
	"context"
	"io"
	"log"
	"os"
)

// Contains options for building the sandboxless browser.
type BrowserOpts struct {
	Headless bool            //Whether to launch the browser without a GUI.
	Leakless bool            //Whether to automatically close the browser when the process exits.
	DevTools bool            //Whether to show the dev tools console; requires the browser to not be headless.
	Logger   io.Writer       //A generic writer to which logs should be sent.
	Ctx      context.Context //A context for the browser.
	GPU      bool            //Whether the GPU should be used.
}

// Returns the default options for the sandboxless browser.
func DefaultBrowserOpts() BrowserOpts {
	return BrowserOpts{
		Headless: true,
		Leakless: true,
		DevTools: false,
		Logger:   os.Stdout,
		Ctx:      context.Background(),
		GPU:      false,
	}
}

/*
Returns the default options for the sandboxless browser, with a logger. This logger is
request independent and should be used for global debugging, not on a per-request basis.
*/
func DefaultBrowserOptsWLogger(l *log.Logger) BrowserOpts {
	defaults := DefaultBrowserOpts()
	defaults.Logger = l.Writer()
	return defaults
}

// Returns the default options for the sandboxless browser when running manual tests.
func DefaultBrowserOptsDbg() BrowserOpts {
	defaults := DefaultBrowserOpts()
	defaults.Headless = false
	defaults.Leakless = false //Keeps the browser open for post-run examination
	defaults.DevTools = true
	return defaults
}
