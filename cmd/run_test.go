package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunCommand_CommitMessageFile_SubjectNotEmptyAndLengthValid(t *testing.T) {
	// Create a temporary commit message file for the test
	tempFile, err := os.CreateTemp("", "commit-msg-test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	commitMessageContent := "feat: Add new feature\n\nThis is a test commit body."
	if _, err := tempFile.WriteString(commitMessageContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	_ = tempFile.Close()

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
		// Pass a dummy cobra.Command as the `cmd` argument as it's not directly used
		// within the `runCommandLogic` for this specific case.
		dummyCmd := rootCmd // Using rootCmd for consistency, though it's not directly used here.
		runCommandLogic(dummyCmd, []string{tempFile.Name()}, mockExit)
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

	// The run command for commit message files should not print success messages to stdout
	// as per the requirement to remove verbose output.
	expectedStdout := ""
	if string(capturedOutput) != expectedStdout {
		t.Errorf("Expected stdout to indicate success, got: %q", string(capturedOutput))
	}

	if len(capturedError) > 0 {
		t.Errorf("Expected no stderr output, but got: %s", strings.TrimSpace(string(capturedError)))
	}
}

func TestRunCommand_CommitMessageFile_SubjectEmpty(t *testing.T) {
	// Create a temporary commit message file for the test
	tempFile, err := os.CreateTemp("", "commit-msg-empty-subject-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	commitMessageContent := "\nThis is a test commit body with an empty subject."
	if _, err := tempFile.WriteString(commitMessageContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	_ = tempFile.Close()

	originalStderr := os.Stderr
	rErr, wErr, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}
	os.Stderr = wErr

	capturedExitCode := -1
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

		mockExit := func(code int) {
			panic(code)
		}
		dummyCmd := rootCmd
		runCommandLogic(dummyCmd, []string{tempFile.Name()}, mockExit)
	}()

	<-done
	wErr.Close()
	os.Stderr = originalStderr

	capturedError, err := io.ReadAll(rErr)
	if err != nil {
		t.Fatalf("Failed to read captured stderr: %v", err)
	}

	if capturedExitCode != 1 {
		t.Errorf("Expected os.Exit(1) but got %d", capturedExitCode)
	}

	expectedStderrPart := "Commit message check failed:\ncommit message subject cannot be empty"
	if !strings.Contains(string(capturedError), expectedStderrPart) {
		t.Errorf("Expected stderr to contain \"%s\", got \"%s\"", expectedStderrPart, strings.TrimSpace(string(capturedError)))
	}
}

func TestRunCommand_NoArgs(t *testing.T) {
	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w

	capturedExitCode := -1
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
		mockExit := func(code int) {
			panic(code)
		}
		dummyCmd := rootCmd
		runCommandLogic(dummyCmd, []string{}, mockExit)
	}()

	<-done
	w.Close()
	os.Stdout = originalStdout

	capturedOutput, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("Failed to read captured output: %v", err)
	}

	if capturedExitCode != 0 {
		t.Errorf("Expected os.Exit(0) but got %d", capturedExitCode)
	}

	expectedStdout := "evolved-commit run: No commit message file provided. Running general checks (not yet implemented).\n"
	if string(capturedOutput) != expectedStdout {
		t.Errorf("Expected stdout to contain \"%s\", got \"%s\"", expectedStdout, strings.TrimSpace(string(capturedOutput)))
	}
}
