package githooks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrGitDirNotFound indicates that the .git directory could not be found.
var ErrGitDirNotFound = errors.New(".git directory not found")

const (
	// PreCommitHook is the name of the pre-commit Git hook.
	PreCommitHook = "pre-commit"
	// HookScriptContent is the content that will be written to the hook files.
	// It executes the evolved-commit run command.
	HookScriptContent = "#!/bin/sh\n\n# evolved-commit hook - do not modify directly\nevolved-commit run\n"
)

// findGitDir attempts to locate the .git directory by checking the current directory
// and its parents.
func findGitDir(startPath string) (string, error) {
	currentPath, err := filepath.Abs(startPath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	for {
		gitPath := filepath.Join(currentPath, ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			return gitPath, nil
		}

		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath { // Reached root directory
			break
		}
		currentPath = parentPath
	}
	return "", ErrGitDirNotFound
}

// Install sets up the necessary Git hooks for evolved-commit in the given repository path.
// If repoPath is empty, it attempts to find the .git directory from the current working directory.
func Install(repoPath string) error {
	if repoPath == "" {
		repoPath = "." // Start search from current directory
	}

	gitDir, err := findGitDir(repoPath)
	if err != nil {
		return fmt.Errorf("could not install hooks: %w", err)
	}

	hooksDir := filepath.Join(gitDir, "hooks")
	// Ensure the hooks directory exists
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory %s: %w", hooksDir, err)
	}

	hookPath := filepath.Join(hooksDir, PreCommitHook)
	// Write the hook script content to the file
	if err := os.WriteFile(hookPath, []byte(HookScriptContent), 0755); err != nil {
		return fmt.Errorf("failed to write %s hook: %w", PreCommitHook, err)
	}

	return nil
}

// Uninstall removes Git hooks installed by evolved-commit from the given repository path.
// If repoPath is empty, it attempts to find the .git directory from the current working directory.
func Uninstall(repoPath string) error {
	if repoPath == "" {
		repoPath = "." // Start search from current directory
	}

	gitDir, err := findGitDir(repoPath)
	if err != nil {
		return fmt.Errorf("could not uninstall hooks: %w", err)
	}

	hooksDir := filepath.Join(gitDir, "hooks")
	hookPath := filepath.Join(hooksDir, PreCommitHook)

	// Check if the hook file exists and remove it.
	// Only remove if it contains the evolved-commit specific content.
	if _, err := os.Stat(hookPath); err == nil {
		content, readErr := os.ReadFile(hookPath)
		if readErr != nil {
			return fmt.Errorf("failed to read %s hook for verification: %w", PreCommitHook, readErr)
		}

		if string(content) == HookScriptContent {
			if err := os.Remove(hookPath); err != nil {
				return fmt.Errorf("failed to remove %s hook: %w", PreCommitHook, err)
			}
		} else {
			// If the hook exists but was modified by a user or another tool,
			// we should not blindly delete it.
			return fmt.Errorf("pre-commit hook exists but appears modified. Not removed to prevent data loss. Please remove manually if desired: %s", hookPath)
		}
	}
	// If file doesn't exist, it's already uninstalled, so no error.
	return nil
}
