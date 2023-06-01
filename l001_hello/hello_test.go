package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Andri")
	want := "Hello, Andri"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
