package githooks

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// setupTestRepo creates a temporary directory with a simulated .git directory
// and returns its path.
func setupTestRepo(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir()
	gitDir := filepath.Join(tempDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	return tempDir
}

func TestFindGitDir(t *testing.T) {
	// Test case 1: .git in current directory
	repoPath := setupTestRepo(t)
	expectedGitDir := filepath.Join(repoPath, ".git")
	
	gitDir, err := findGitDir(repoPath)
	if err != nil {
		t.Fatalf("findGitDir failed: %v", err)
	}
	if gitDir != expectedGitDir {
		t.Errorf("Expected git dir %s, got %s", expectedGitDir, gitDir)
	}

	// Test case 2: .git in parent directory
	subDir := filepath.Join(repoPath, "sub", "dir")
	err = os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}
	
	gitDir, err = findGitDir(subDir)
	if err != nil {
		t.Fatalf("findGitDir failed from subdirectory: %v", err)
	}
	if gitDir != expectedGitDir {
		t.Errorf("Expected git dir %s, got %s when searching from %s", expectedGitDir, gitDir, subDir)
	}

	// Test case 3: .git not found
	nonRepoDir := t.TempDir()
	_, err = findGitDir(nonRepoDir)
	if !errors.Is(err, ErrGitDirNotFound) {
		t.Errorf("Expected ErrGitDirNotFound, got %v", err)
	}
}

func TestInstall(t *testing.T) {
	repoPath := setupTestRepo(t)

	err := Install(repoPath)
	if err != nil {
		t.Errorf("Install failed: %v", err)
	}

	// Future: Check if hook file exists and has correct content
}

func TestUninstall(t *testing.T) {
	repoPath := setupTestRepo(t)

	// First install a dummy hook to ensure uninstall has something to remove
	// This is a placeholder for now, actual hook creation will be added later.
	hooksDir := filepath.Join(repoPath, ".git", "hooks")
	os.MkdirAll(hooksDir, 0755)
	dummyHookPath := filepath.Join(hooksDir, "pre-commit")
	os.WriteFile(dummyHookPath, []byte("#!/bin/sh\necho 'dummy hook'"), 0755)

	err := Uninstall(repoPath)
	if err != nil {
		t.Errorf("Uninstall failed: %v", err)
	}

	// Future: Check if hook file is removed or content is reverted
}

func TestInstall_NoGitDir(t *testing.T) {
	nonRepoDir := t.TempDir()
	err := Install(nonRepoDir)
	if !errors.Is(err, ErrGitDirNotFound) {
		t.Errorf("Expected Install to fail with ErrGitDirNotFound, got %v", err)
	}
}

func TestUninstall_NoGitDir(t *testing.T) {
	nonRepoDir := t.TempDir()
	err := Uninstall(nonRepoDir)
	if !errors.Is(err, ErrGitDirNotFound) {
		t.Errorf("Expected Uninstall to fail with ErrGitDirNotFound, got %v", err)
	}
}
