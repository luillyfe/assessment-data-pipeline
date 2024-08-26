package main

import (
	"os"
	"testing"
)

// TestReadFile tests the readFile function.
func TestReadFile(t *testing.T) {
	// Temporary file for testing
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %s", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Write some content to the temporary file
	content := "Hello, World!"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %s", err)
	}

	// Close the file so it can be read
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %s", err)
	}

	// Test case: successful file read
	t.Run("successful file read", func(t *testing.T) {
		got, err := readFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("Expected no error, got %s", err)
		}
		if got != content {
			t.Errorf("Expected content %s, got %s", content, got)
		}
	})

	// Test case: file does not exist
	t.Run("file does not exist", func(t *testing.T) {
		_, err := readFile("nonexistentfile.txt")
		if err == nil {
			t.Fatal("Expected an error, got none")
		}
	})
}
