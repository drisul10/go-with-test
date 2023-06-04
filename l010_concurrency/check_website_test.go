package concurrency

import (
	"net/http"
	"net/http/httptest"
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

func TestRacerWebsites(t *testing.T) {
	t.Run("should return fastest server", func(t *testing.T) {
		slowServer := makeDelayedServer(1 * time.Nanosecond)
		fastServer := makeDelayedServer(0 * time.Nanosecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, err := RacerWebsites(slowURL, fastURL)
		if err != nil {
			t.Errorf("did not expect an error but got one: %v", err.Error())
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error if a server does not respond within 10s", func(t *testing.T) {
		server := makeDelayedServer(11 * time.Millisecond)

		defer server.Close()

		tenMilisecondTimeout := 10 * time.Millisecond

		_, err := ConfigurableRacerWebsites(server.URL, server.URL, tenMilisecondTimeout)

		if err == nil {
			t.Errorf("expected an error but didn't get one")
		}
	})
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
