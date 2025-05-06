package main

import "testing"

func TestHelloWorld(t *testing.T) {
	result := Hello()
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", Hello())
	}

	t.Logf("Test passed: %s", Hello())
}
