package concurrency

import (
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

func RacerWebsites(webA, webB string) (winner string) {
	startA := time.Now()
	http.Get(webA)
	durationA := time.Since(startA)

	startB := time.Now()
	http.Get(webB)
	durationB := time.Since(startB)

	if durationA < durationB {
		return webA
	}

	return webB
}
