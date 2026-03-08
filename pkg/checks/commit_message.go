package checks

import (
	"errors"
	"fmt"
	"strings"
)

// ErrCommitMessageSubjectEmpty indicates that the commit message subject line is empty.
var ErrCommitMessageSubjectEmpty = errors.New("commit message subject cannot be empty")

// ErrCommitMessageSubjectTooLong indicates that the commit message subject line exceeds the maximum allowed length.
var ErrCommitMessageSubjectTooLong = errors.New("commit message subject is too long")

const MaxSubjectLength = 50

// CheckCommitMessageSubjectNotEmpty validates that the commit message subject line is not empty.
// It returns an error if the subject is empty, providing a clear explanation and fix guidance.
func CheckCommitMessageSubjectNotEmpty(subjectLine string) error {
	trimmedSubject := strings.TrimSpace(subjectLine)
	if trimmedSubject == "" {
		return fmt.Errorf("%w: The first line of your commit message is empty.\nFix: Please provide a concise subject line for your commit.", ErrCommitMessageSubjectEmpty)
	}
	return nil
}

// CheckCommitMessageSubjectLength validates that the commit message subject line does not exceed MaxSubjectLength characters.
// It returns an error if the subject is too long, providing a clear explanation and fix guidance.
func CheckCommitMessageSubjectLength(subjectLine string) error {
	trimmedSubject := strings.TrimSpace(subjectLine)
	if len(trimmedSubject) > MaxSubjectLength {
		return fmt.Errorf("%w: The first line of your commit message is %d characters long, but it should not exceed %d characters.\nFix: Shorten your commit message subject line to be %d characters or less.", ErrCommitMessageSubjectTooLong, len(trimmedSubject), MaxSubjectLength, MaxSubjectLength)
	}
	return nil
}
