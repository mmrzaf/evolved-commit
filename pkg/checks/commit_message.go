package checks

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// ErrCommitMessageSubjectEmpty indicates that the commit message subject line is empty.
var ErrCommitMessageSubjectEmpty = errors.New("commit message subject cannot be empty")

// ErrCommitMessageSubjectTooLong indicates that the commit message subject line exceeds the maximum allowed length.
var ErrCommitMessageSubjectTooLong = errors.New("commit message subject is too long")

// ErrCommitMessageSubjectTrailingPeriod indicates that the commit message subject line ends with a period.
var ErrCommitMessageSubjectTrailingPeriod = errors.New("commit message subject should not end with a period")

// ErrCommitMessageSubjectStartsWithLowercase indicates that the commit message subject line starts with a lowercase letter.
var ErrCommitMessageSubjectStartsWithLowercase = errors.New("commit message subject should start with an uppercase letter")

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

// CheckCommitMessageSubjectNoTrailingPeriod validates that the commit message subject line does not end with a period.
// It returns an error if the subject ends with a period, providing a clear explanation and fix guidance.
func CheckCommitMessageSubjectNoTrailingPeriod(subjectLine string) error {
	trimmedSubject := strings.TrimSpace(subjectLine)
	if strings.HasSuffix(trimmedSubject, ".") {
		return fmt.Errorf("%w: The first line of your commit message ends with a period.\nFix: Remove the trailing period from the subject line. Example: 'feat: Add user authentication'", ErrCommitMessageSubjectTrailingPeriod)
	}
	return nil
}

// CheckCommitMessageSubjectStartsWithUppercase validates that the commit message subject line starts with an uppercase letter.
// It returns an error if the subject starts with a lowercase letter, providing a clear explanation and fix guidance.
func CheckCommitMessageSubjectStartsWithUppercase(subjectLine string) error {
	trimmedSubject := strings.TrimSpace(subjectLine)
	if len(trimmedSubject) == 0 {
		return nil // This case is covered by CheckCommitMessageSubjectNotEmpty
	}

	firstChar := rune(trimmedSubject[0])
	if unicode.IsLower(firstChar) {
		return fmt.Errorf("%w: The first line of your commit message starts with a lowercase letter. \nFix: Start your commit message subject line with an uppercase letter. Example: 'Feat: Add user authentication'", ErrCommitMessageSubjectStartsWithLowercase)
	}
	return nil
}
