package cmd

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Helper function to create a temporary commit message file
func createTempCommitMsgFile(t *testing.T, content string) *os.File {
	tmpFile, err := os.CreateTemp("", "COMMIT_EDITMSG_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	_, err = tmpFile.WriteString(content)
	if err != nil {
		os.Remove(tmpFile.Name()) // Clean up on write error
		tmpFile.Close()
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close() // Close the file after writing to ensure content is flushed
	return tmpFile
}

// TestRunCommand_CommitMessageFile_SubjectNotEmptyAndLengthValid tests the run command
// with a commit message file that should pass all subject line checks.
func TestRunCommand_CommitMessageFile_SubjectNotEmptyAndLengthValid(t *testing.T) {
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

	// Create a temporary commit message file that should pass all checks
	// Subject: "Feat: Add user authentication" - Non-empty, length < 50, no trailing period, starts with uppercase.
	commitMsgContent := "Feat: Add user authentication\n\nThis is the body of the commit message."
	tmpCommitMsgFile := createTempCommitMsgFile(t, commitMsgContent)
	defer os.Remove(tmpCommitMsgFile.Name()) // Clean up the file

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
		runCommandLogic(dummyCmd, []string{tmpCommitMsgFile.Name()}, mockExit)
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

	if capturedExitCode != 0 {
		t.Errorf("Expected os.Exit(0) but got %d. Stderr: %s", capturedExitCode, strings.TrimSpace(string(capturedError)))
	}

	if len(capturedError) > 0 {
		t.Errorf("Expected no stderr output, but got: %s", strings.TrimSpace(string(capturedError)))
	}

	if len(capturedOutput) > 0 {
		t.Errorf("Expected no stdout output, but got: %s", strings.TrimSpace(string(capturedOutput)))
	}
}

// TestRunCommand_NoArgs tests the run command when no arguments are provided.
func TestRunCommand_NoArgs(t *testing.T) {
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
		runCommandLogic(dummyCmd, []string{}, mockExit) // No arguments
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

	if capturedExitCode != 0 {
		t.Errorf("Expected os.Exit(0) but got %d. Stderr: %s", capturedExitCode, strings.TrimSpace(string(capturedError)))
	}

	expectedStdoutPart := "evolved-commit run: No commit message file provided. Running general checks (not yet implemented)."
	if !strings.Contains(string(capturedOutput), expectedStdoutPart) {
		t.Errorf("Expected stdout to contain \"%s\", got \"%s\"", expectedStdoutPart, strings.TrimSpace(string(capturedOutput)))
	}

	if len(capturedError) > 0 {
		t.Errorf("Expected no stderr output, but got: %s", strings.TrimSpace(string(capturedError)))
	}
}
