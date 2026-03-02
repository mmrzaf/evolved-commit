package githooks

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
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
	hooksDir := filepath.Join(repoPath, ".git", "hooks")
	hookPath := filepath.Join(hooksDir, PreCommitHook)

	// Ensure hook file does not exist initially
	if _, err := os.Stat(hookPath); err == nil {
		t.Fatalf("Hook file %s should not exist before install", hookPath)
	}

	err := Install(repoPath)
	if err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	// Check if hook file exists
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Errorf("Expected %s hook file to be created, but it does not exist", PreCommitHook)
	}

	// Check content of the hook file
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read hook file: %v", err)
	}
	if string(content) != HookScriptContent {
		t.Errorf("Expected hook content\n%q\n got\n%q", HookScriptContent, string(content))
	}

	// Test installing again (idempotency check)
	err = Install(repoPath)
	if err != nil {
		t.Fatalf("Second Install failed: %v", err)
	}
	// Content should still be the same
	content, err = os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read hook file after second install: %v", err)
	}
	if string(content) != HookScriptContent {
		t.Errorf("Expected hook content to be idempotent, got\n%q", string(content))
	}
}

func TestUninstall(t *testing.T) {
	repoPath := setupTestRepo(t)
	hooksDir := filepath.Join(repoPath, ".git", "hooks")
	hookPath := filepath.Join(hooksDir, PreCommitHook)

	// Install the hook first
	err := Install(repoPath)
	if err != nil {
		t.Fatalf("Failed to install hook for uninstall test: %v", err)
	}

	// Ensure hook file exists before uninstall
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Fatalf("Hook file %s should exist before uninstall", hookPath)
	}

	err = Uninstall(repoPath)
	if err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	// Check if hook file is removed
	if _, err := os.Stat(hookPath); !os.IsNotExist(err) {
		t.Errorf("Expected %s hook file to be removed, but it still exists", PreCommitHook)
	}

	// Test uninstalling again (idempotency check - no error if already gone)
	err = Uninstall(repoPath)
	if err != nil {
		t.Fatalf("Second Uninstall failed: %v", err)
	}

	// Test uninstalling a non-evolved-commit hook
	otherHookPath := filepath.Join(hooksDir, "post-commit")
	otherContent := "#!/bin/sh\necho 'another hook'\n"
	os.WriteFile(otherHookPath, []byte(otherContent), 0755)
	err = Uninstall(repoPath) // Uninstall still tries to remove PreCommitHook, not otherHookPath
	if err != nil {
		// Should not error if the target hook for uninstall doesn't exist, as the method is for pre-commit
		// However, this test specifically checks for the *pre-commit* hook. The check below is more relevant.
	}
	// Verify the other hook is still there
	if _, err := os.Stat(otherHookPath); os.IsNotExist(err) {
		t.Errorf("Expected other hook %s to remain after uninstall, but it was removed", otherHookPath)
	}

	// Test uninstalling an evolved-commit hook that was modified
	err = Install(repoPath) // Re-install our hook
	if err != nil {
		t.Fatalf("Failed to reinstall hook for modified test: %v", err)
	}
	modifiedContent := "#!/bin/sh\necho 'user modified hook'\n"
	os.WriteFile(hookPath, []byte(modifiedContent), 0755) // Modify it

	err = Uninstall(repoPath)
	if err == nil || !strings.Contains(err.Error(), "appears modified") {
		t.Errorf("Expected uninstall to fail for modified hook with 'appears modified' error, got: %v", err)
	}
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Errorf("Expected modified hook to NOT be removed, but it was.")
	}
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
