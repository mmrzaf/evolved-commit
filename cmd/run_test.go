package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	// Store original stdout and stderr
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Create pipes to capture stdout and stderr
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = wOut // Redirect stdout to the pipe writer

	rErr, wErr, err := os.Pipe()
	if err != nil {
		_ = wOut.Close() // Close the stdout writer as it won't be used now.
		os.Stdout = originalStdout // Restore stdout before failing
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}
	os.Stderr = wErr // Redirect stderr to the pipe writer

	// Variable to capture the exit code from the goroutine
	capturedExitCode := -1 // Default to an invalid exit code

	// Use a goroutine to execute the command to allow capturing output
	done := make(chan struct{})
	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				if code, ok := r.(int); ok {
					capturedExitCode = code
				} else {
					t.Errorf("Recovered from unexpected panic: %v", r)
				}
			}
		}()

		// Create a mock exit function that panics with the exit code
		mockExit := func(code int) {
			panic(code) // Panic with the exit code to stop execution and be recovered
		}

		// Execute the command's logic directly with the mock exit function
		// rootCmd is passed as the first argument as `cmd *cobra.Command`
		runCommandLogic(rootCmd, []string{}, mockExit)
	}()

	// Wait for the command to finish executing
	<-done

	// Close the pipe writers to signal EOF for readers, now that the goroutine has finished.
	wOut.Close()
	wErr.Close()

	// Restore original stdout and stderr immediately after the goroutine finishes
	// and pipes are closed, but before reading the captured output.
	// This ensures that any subsequent t.Fatalf or other logging uses the actual stdout/stderr.
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	// Read all captured output
	capturedOutput, err := io.ReadAll(rOut)
	if err != nil {
		t.Fatalf("Failed to read captured stdout: %v", err)
	}
	capturedError, err := io.ReadAll(rErr)
	if err != nil {
		t.Fatalf("Failed to read captured stderr: %v", err)
	}

	// Verify the exit code.
	if capturedExitCode != 0 {
		t.Errorf("Expected os.Exit(0) but got %d", capturedExitCode)
	}

	expectedStdout := "evolved-commit run: No commit message file provided. Running general checks (not yet implemented)."
	if !strings.Contains(string(capturedOutput), expectedStdout) {
		t.Errorf("Expected stdout to contain \"%s\", got \"%s\"", expectedStdout, strings.TrimSpace(string(capturedOutput)))
	}

	if len(capturedError) > 0 {
		t.Errorf("Expected no stderr output, but got: %s", strings.TrimSpace(string(capturedError)))
	}
}
