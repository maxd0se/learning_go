// To crawl a site run the following command: go run gophish.go URL_TO_CRAWL. Include the protocol (http or https)

package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// function to pull the href attributes
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over token attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	// "bare" return will return the variables (ok, href) as
	// defined in the function definition
	return
}

// extracts all http** links from provided url
func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)
	defer func() {
		chFinished <- true
	}()
	if err != nil {
		fmt.Println("ERROR: Failed to crawl:", url)
		return
	}
	b := resp.Body
	defer b.Close() 
	z := html.NewTokenizer(b)
	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			return
		case tt == html.StartTagToken:
			t := z.Token()
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			ok, url := getHref(t)
			if !ok {
				continue
			}
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}
	}
}
func main() {
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]
	// Channels
	chUrls := make(chan string)
	chFinished := make(chan bool)
	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}

	for c := 0; c < len(seedUrls); {
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")
	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}
	close(chUrls)
}
