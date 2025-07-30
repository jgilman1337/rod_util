package pkg

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
)

// Builds a browser without GPU acceleration or sandboxing. See: https://go-rod.github.io/#/custom-launch
func BuildSandboxless(opts BrowserOpts) (*rod.Browser, *launcher.Launcher, error) {
	//Build the launcher
	lc := launcher.New().
		Headless(opts.Headless).
		Devtools(opts.DevTools).
		Set("disable-gpu", strconv.FormatBool(!opts.GPU)).
		NoSandbox(true).
		Logger(opts.Logger).
		Leakless(opts.Leakless).
		Context(opts.Ctx)

	//Launch the browser
	u, err := lc.Launch() //This launches the browser and returns the WebSocket URL
	if err != nil {
		lc.Kill()
		lc.Cleanup()
		return nil, nil, err
	}
	browser := rod.New().ControlURL(u)
	if err := browser.Connect(); err != nil {
		RodFreeManual(browser, lc)
		return nil, nil, err
	}

	return browser, lc, nil
}

// Creates a webpage in a browser using the same interface as `browser.MustPage()`, but without using `panic()`. Instead, an error is returned.
func Page(b *rod.Browser, url ...string) (*rod.Page, error) {
	return b.Page(proto.TargetCreateTarget{URL: strings.Join(url, "/")})
}

// Creates a blank webpage in a web browser.
func BlankPage(b *rod.Browser, url ...string) *rod.Page {
	p, _ := b.Page(proto.TargetCreateTarget{URL: strings.Join(url, "/")})
	return p
}

// Gets a JSON string from a page returning `application/json`.
func PageJSON(p *rod.Page) (string, error) {
	jsonz, err := p.Eval("() => document.documentElement.innerText")
	if err != nil {
		return "", err
	}
	return jsonz.Value.String(), nil
}

// Cleanup function for Rod browser and launcher.
func RodFree(b *rod.Browser, l *launcher.Launcher) {
	b.MustClose()
	l.Kill()
	l.Cleanup()
}

// Cleanup function for Rod browser and launcher with manual cleanup of the data directory.
func RodFreeManual(b *rod.Browser, l *launcher.Launcher) {
	b.MustClose()
	l.Kill()
	//l.Cleanup()

	//Taken from `Launcher.Cleanup()`, sans the waiting channel since we are running this semi-detached
	dir := l.Get(flags.UserDataDir)
	_ = os.RemoveAll(dir)
}
