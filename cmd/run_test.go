package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// This test verifies the 'run' command's behavior when provided with a commit message file
// that passes all the defined checks. It ensures that the command executes successfully
// and does not produce any error output.
func TestRunCommandWithCommitMessageSuccess(t *testing.T) {
	// --- Setup for temporary commit message file ---
	// Create a temporary directory for the test.
	tempDir := t.TempDir()
	// Define the path for the temporary commit message file.
	commitMsgFilePath := filepath.Join(tempDir, "COMMIT_EDITMSG")
	// Content for a valid commit message subject (imperative, starts with uppercase, no period, not too long).
	commitMsgContent := "Feat: Implement user authentication"
	// Write the content to the temporary file.
	err := os.WriteFile(commitMsgFilePath, []byte(commitMsgContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary commit message file: %v", err)
	}

	// --- Output Capture Setup ---
	// Store original stdout and stderr to restore them later.
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// Create a pipe for capturing stdout.
	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = wOut // Redirect stdout to the pipe writer

	// Create a pipe for capturing stderr.
	rErr, wErr, err := os.Pipe()
	if err != nil {
		_ = wOut.Close() // Ensure the stdout writer is closed if stderr pipe creation fails.
		os.Stdout = originalStdout
		t.Fatalf("Failed to create stderr pipe: %v", err)
	}
	os.Stderr = wErr // Redirect stderr to the pipe writer

	// Variable to capture the exit code from the goroutine.
	capturedExitCode := -1 // Default to an invalid exit code.

	// Use a goroutine to execute the command logic to allow capturing output
	// and handling panics from the mock exit function.
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

		// Create a mock exit function that panics with the exit code.
		mockExit := func(code int) {
			panic(code) // Panic with the exit code to stop execution and be recovered.
		}

		// Execute the command's logic directly.
		// Pass a dummy cobra.Command as `cmd` argument as it's not directly used
		// within `runCommandLogic` for this test.
		dummyCmd := &cobra.Command{}
		runCommandLogic(dummyCmd, []string{commitMsgFilePath}, mockExit)
	}()

	// Wait for the command to finish executing.
	<-done

	// Close the pipe writers to signal EOF for readers, now that the goroutine has finished.
	wOut.Close()
	wErr.Close()

	// Restore original stdout and stderr.
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	// Read all captured output.
	capturedOutput, err := io.ReadAll(rOut)
	if err != nil {
		t.Fatalf("Failed to read captured stdout: %v", err)
	}
	capturedError, err := io.ReadAll(rErr)
	if err != nil {
		t.Fatalf("Failed to read captured stderr: %v", err)
	}

	// --- Assertions ---
	// Verify the exit code. For success, it should be 0.
	if capturedExitCode != 0 {
		t.Errorf("Expected os.Exit(0) but got %d. Stderr: %s", capturedExitCode, strings.TrimSpace(string(capturedError)))
	}

	// Verify no stdout or stderr output for a successful run.
	if len(capturedOutput) > 0 {
		t.Errorf("Expected no stdout output, but got: %s", strings.TrimSpace(string(capturedOutput)))
	}
	if len(capturedError) > 0 {
		t.Errorf("Expected no stderr output, but got: %s", strings.TrimSpace(string(capturedError)))
	}
}

// This test verifies the 'run' command's behavior when provided with a commit message file
// that fails the imperative mood check.
func TestRunCommandWithCommitMessageFailure_Imperative(t *testing.T) {
	// --- Setup for temporary commit message file ---
	tempDir := t.TempDir()
	commitMsgFilePath := filepath.Join(tempDir, "COMMIT_EDITMSG")
	// Content for an invalid commit message subject (third-person singular).
	commitMsgContent := "Feat: Implements user authentication"
	err := os.WriteFile(commitMsgFilePath, []byte(commitMsgContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary commit message file: %v", err)
	}

	// --- Output Capture Setup ---
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	rOut, wOut, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = wOut

	rErr, wErr, err := os.Pipe()
	if err != nil {
		_ = wOut.Close()
		os.Stdout = originalStdout
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

		dummyCmd := &cobra.Command{}
		runCommandLogic(dummyCmd, []string{commitMsgFilePath}, mockExit)
	}()

	<-done
	wOut.Close()
	wErr.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	capturedOutput, err := io.ReadAll(rOut)
	if err != nil {
		t.Fatalf("Failed to read captured stdout: %v", err)
	}
	capturedError, err := io.ReadAll(rErr)
	if err != nil {
		t.Fatalf("Failed to read captured stderr: %v", err)
	}

	// --- Assertions ---
	// Verify the exit code. For failure, it should be 1.
	if capturedExitCode != 1 {
		t.Errorf("Expected os.Exit(1) but got %d. Stderr: %s", capturedExitCode, strings.TrimSpace(string(capturedError)))
	}

	// Verify no stdout output, but stderr output for a failed run.
	if len(capturedOutput) > 0 {
		t.Errorf("Expected no stdout output, but got: %s", strings.TrimSpace(string(capturedOutput)))
	}

	expectedStderrPart := "Commit message checks failed:\n- commit message subject should use imperative mood: The first word 'Implements' in your commit message subject appears to be in the third-person singular."
	if !strings.Contains(string(capturedError), expectedStderrPart) {
		t.Errorf("Expected stderr to contain \"%s\", got \"%s\"", expectedStderrPart, strings.TrimSpace(string(capturedError)))
	}
}
