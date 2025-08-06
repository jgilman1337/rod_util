package test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/go-rod/rod"

	rutil "github.com/jgilman1337/rod_util/pkg"
)

func TestBasic(t *testing.T) {
	//Set test vars
	url := "https://example.com"
	selector := "body > div:nth-child(1) > h1:nth-child(1)"
	expected := "Example Domain"

	//Launch Rod
	browser, launcher, err := rutil.BuildSandboxless(rutil.DefaultBrowserOpts())
	if err != nil {
		t.Fatalf("Failed to launch browser: %s", err)
	}
	defer rutil.RodFree(browser, launcher)

	//Setup a page
	page, err := rutil.Page(browser, url)
	if err != nil {
		t.Fatalf("Failed to create a webpage: %s", err)
	}
	defer page.Close()

	//Wait for the page to load
	if err := page.WaitDOMStable(time.Second, 0); err != nil {
		t.Fatalf("WaitDOMStable failed: %s", err)
	}

	//Get "actual" content
	elem, err := page.Element(selector)
	if err != nil || elem == nil {
		t.Fatalf("Failed to find target: %s", err)
	}
	actual, err := elem.Text()
	if err != nil {
		t.Fatalf("Failed to grab text from element: %s", err)
	}

	if expected != actual {
		t.Fatalf("Unexpected result\n  Actual:   %s\n  Expected: %s", actual, expected)
	}
}

func TestTimeoutLoad(t *testing.T) {
	//Set test vars
	srv := newHangServer()
	//Force-kills the httptest.server
	defer func() {
		srv.Config.Close()
		srv.Listener.Close()
	}()

	//Launch Rod
	browser, launcher, err := rutil.BuildSandboxless(rutil.DefaultBrowserOpts())
	if err != nil {
		t.Fatalf("Failed to launch browser: %s", err)
	}
	defer rutil.RodFree(browser, launcher)

	//Setup a page, with timeout
	page := rutil.BlankPage(browser).Timeout(5 * time.Second)
	defer page.Close()
	if err := page.Navigate(srv.URL); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			t.Log("Successfully stopped page load early")
		} else {
			t.Fatalf("Failed to create a webpage: %s", err)
		}
	}
}

func TestTimeoutDOM(t *testing.T) {
	test := func(_ *rod.Browser, _ *rod.Page) {
		//Nada
	}
	htmlRunner(t, "unstable.html", false, 5, test)
}

func TestHasElement(t *testing.T) {
	selector := "#e"
	expected := "blue"

	test := func(b *rod.Browser, p *rod.Page) {
		//Get an element
		elem, err := rutil.SafeSelect(p, selector)
		if err != nil {
			t.Fatalf("Has() failed: %s", err)
		}

		//Get text from the element
		actual := "<no match>"
		if elem != nil {
			if err := rod.Try(func() {
				actual = elem.MustText()
			}); err != nil {
				t.Fatalf("MustText() failed: %s", err)
			}
		}

		//Test
		t.Logf("First DOM item: %s", actual)
		if expected != actual {
			t.Logf("%s != %s", expected, actual)
		} else {
			t.Logf("%s = %s", expected, actual)
		}
	}
	htmlRunner(t, "dom_cycle.html", true, 5, test)
}

func TestGetParseJSON(t *testing.T) {
	//Sample struct
	type album struct {
		UserId int    `json:"userId"`
		ID     int    `json:"id"`
		Title  string `json:"title"`
	}

	//Set test vars
	url := "https://jsonplaceholder.typicode.com/albums"
	itemIdx := 69
	expected := album{
		UserId: 7,
		ID:     70,
		Title:  "et deleniti unde",
	}

	//Launch Rod
	browser, launcher, err := rutil.BuildSandboxless(rutil.DefaultBrowserOpts())
	if err != nil {
		t.Fatalf("Failed to launch browser: %s", err)
	}
	defer rutil.RodFree(browser, launcher)

	//Setup a page
	page, err := rutil.Page(browser, url)
	if err != nil {
		t.Fatalf("Failed to create a webpage: %s", err)
	}
	defer page.Close()

	//Get JSON and unmarshal it
	jstr, err := rutil.PageJSON(page)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %s", err)
	}
	albums := make([]album, 0)
	if err := json.Unmarshal([]byte(jstr), &albums); err != nil {
		t.Fatalf("Failed to unmarshal returned JSON: %s", err)
	}

	//Test for equality (shallow)
	actual := albums[itemIdx]
	if expected != actual {
		t.Fatalf("Unexpected result\n  Actual:   %+v\n  Expected: %+v", actual, expected)
	}
}
