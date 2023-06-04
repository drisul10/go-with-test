package concurrency

import (
	"fmt"
	"net/http"
	"time"
)

type WebsiteChecker func(string) bool

type ResultWebsiteChecker struct {
	string
	bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
	results := make(map[string]bool)
	resultsChan := make(chan ResultWebsiteChecker)

	for _, url := range urls {
		go func(u string) {
			resultsChan <- ResultWebsiteChecker{u, wc(u)}
		}(url)
	}

	for i := 0; i < len(urls); i++ {
		result := <-resultsChan
		results[result.string] = result.bool
	}

	return results
}

func RacerWebsites(webA, webB string) (winner string, err error) {
	select {
	case <-ping(webA):
		return webA, nil
	case <-ping(webB):
		return webB, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timed out waiting for %s and %s", webA, webB)
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()

	return ch
}
