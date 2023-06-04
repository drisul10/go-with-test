package concurrency

import (
	"reflect"
	"testing"
)

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
