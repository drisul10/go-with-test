package concurrency

import (
	"reflect"
	"testing"
	"time"
)

func slowStubWebsiteChecker(_ string) bool {
	time.Sleep(20 * time.Millisecond)
	return true
}

func BenchmarkCheckWebsites(b *testing.B) {
	urls := make([]string, 100)
	for i := 0; i < len(urls); i++ {
		urls[i] = "a.url.com"
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CheckWebsites(slowStubWebsiteChecker, urls)
	}
}

func mockWebsiteChecker(url string) bool {
	return url != "waat://de.hell"
}

func TestCheckWebsites(t *testing.T) {
	websites := []string{
		"https://google.com",
		"https://chat.openai.com",
		"waat://de.hell",
	}

	want := map[string]bool{
		"https://google.com":      true,
		"https://chat.openai.com": true,
		"waat://de.hell":          false,
	}

	got := CheckWebsites(mockWebsiteChecker, websites)

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got %v, want %v", got, want)
	}
}
