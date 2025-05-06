package main

import "testing"

func TestHelloWorld(t *testing.T) {
	result := Hello()
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", result)
	}

	t.Logf("Test passed: %s", result)
}
