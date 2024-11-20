package main

import "testing"

func TestHelloWorld(t *testing.T) {
	expected := "Hello, World!"
	got := "Hello, World!"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
