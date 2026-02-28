package cmd

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestInstallCommand(t *testing.T) {
	// Store original stdout
	oldStdout := os.Stdout

	// Create a pipe to capture stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create pipe: %v", err)
	}
	os.Stdout = w // Redirect stdout to the pipe writer

	// Set command arguments for 'install'
	rootCmd.SetArgs([]string{"install"})

	// Execute the command
	err = rootCmd.Execute()
	if err != nil {
		// Restore stdout before failing
		w.Close()
		os.Stdout = oldStdout
		t.Fatalf("install command failed: %v", err)
	}

	// Close the pipe writer and read all captured output
	w.Close()
	capturedOutput, err := io.ReadAll(r)
	if err != nil {
		// Restore stdout before failing
		os.Stdout = oldStdout
		t.Fatalf("Failed to read captured output: %v", err)
	}

	// Restore original stdout
	os.Stdout = oldStdout

	expected := "Git hooks installed successfully."
	if !strings.Contains(string(capturedOutput), expected) {
		t.Errorf("Expected output to contain \"%s\", got \"%s\"", expected, strings.TrimSpace(string(capturedOutput)))
	}
}
