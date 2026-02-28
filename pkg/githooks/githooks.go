package githooks

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrGitDirNotFound indicates that the .git directory could not be found.
var ErrGitDirNotFound = errors.New(".git directory not found")

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

	// Placeholder for actual hook installation logic
	fmt.Printf("Attempting to install hooks in: %s\n", gitDir)
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

	// Placeholder for actual hook uninstallation logic
	fmt.Printf("Attempting to uninstall hooks from: %s\n", gitDir)
	return nil
}
