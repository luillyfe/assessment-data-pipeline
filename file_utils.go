package main

import (
	"fmt"
	"os"
)

// readFile reads the content of a file and returns it as a string.
// It takes the filename as a string argument.
// It returns the content of the file as a string and an error if any.
func readFile(filename string) (string, error) {
	// Read the entire content of the file with the given filename.
	content, err := os.ReadFile(filename)
	if err != nil {
		// If there is an error reading the file, return an empty string and the error.
		return "", fmt.Errorf("error reading file %s: %w", filename, err)
	}
	// If the file is read successfully, return the content of the file as a string and nil error.
	return string(content), nil
}
