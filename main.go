package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type TestEvent struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Output  string  // Captures log messages, errors, or any output
	Elapsed float64 // Elapsed time in seconds
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	// Buffer to hold test outputs, keyed by test name
	outputBuffer := make(map[string]string)

	for scanner.Scan() {
		var event TestEvent
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			fmt.Fprintf(os.Stderr, "error parsing event: %v", err)
			continue
		}

		// Skip empty test names
		if event.Test == "" {
			continue
		}

		// Normalize the test name
		testName := strings.ReplaceAll(event.Test, "_", " ")

		// Accumulate output for each test
		if event.Action == "output" {
			outputBuffer[testName] += event.Output
			continue
		}

		if event.Action == "pass" {
			fmt.Printf("\033[32m✓ %s\033[0m (%.3fs)\n", testName, event.Elapsed)
		} else if event.Action == "fail" {
			// Print the test name, elapsed time, and buffered output for failed tests
			fmt.Printf("\033[31m✗ %s\033[0m (%.3fs)\n", testName, event.Elapsed)
			fmt.Printf("%s\n", outputBuffer[testName])
		}
	}
}
