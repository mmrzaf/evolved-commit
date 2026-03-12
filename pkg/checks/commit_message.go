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

// ErrCommitMessageSubjectNotImperative indicates that the commit message subject line does not use imperative mood.
var ErrCommitMessageSubjectNotImperative = errors.New("commit message subject should use imperative mood")

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

// CheckCommitMessageSubjectImperative validates that the commit message subject line uses imperative mood.
// It returns an error if the subject appears to be in the past tense, third person singular, or present participle,
// providing clear explanation and fix guidance.
func CheckCommitMessageSubjectImperative(subjectLine string) error {
	trimmedSubject := strings.TrimSpace(subjectLine)
	if len(trimmedSubject) == 0 {
		return nil // Covered by CheckCommitMessageSubjectNotEmpty
	}

	// Find the first colon. If it's there and followed by a space, consider the rest as the effective subject.
	// This handles conventional commit style prefixes like "feat: " or "fix(scope): ".
	effectiveSubject := trimmedSubject
	colonIdx := strings.Index(trimmedSubject, ":")
	if colonIdx != -1 && colonIdx+1 < len(trimmedSubject) && trimmedSubject[colonIdx+1] == ' ' {
		effectiveSubject = strings.TrimSpace(trimmedSubject[colonIdx+2:])
	}

	if len(effectiveSubject) == 0 {
		return nil // After removing prefix, subject is empty, covered by not empty check
	}

	// Get the first word of the effective subject
	firstWord := effectiveSubject
	spaceIdx := strings.Index(effectiveSubject, " ")
	if spaceIdx != -1 {
		firstWord = effectiveSubject[:spaceIdx]
	}

	// Basic heuristic: check for common non-imperative verb endings
	lowerFirstWord := strings.ToLower(firstWord)

	if strings.HasSuffix(lowerFirstWord, "ed") && len(lowerFirstWord) > 3 { // Avoid short words like "red", "bed"
		return fmt.Errorf("%w: The first word '%s' in your commit message subject appears to be in the past tense. Fix: Use the imperative mood (command form) for the subject line. Example: 'Fix: Correct critical bug', not 'Fix: Corrected critical bug'.", ErrCommitMessageSubjectNotImperative, firstWord)
	}
	if strings.HasSuffix(lowerFirstWord, "s") && len(lowerFirstWord) > 2 { // Avoid short words like "is", "us"
		return fmt.Errorf("%w: The first word '%s' in your commit message subject appears to be in the third-person singular. Fix: Use the imperative mood (command form) for the subject line. Example: 'Feat: Add user registration', not 'Feat: Adds user registration'.", ErrCommitMessageSubjectNotImperative, firstWord)
	}
	if strings.HasSuffix(lowerFirstWord, "ing") && len(lowerFirstWord) > 4 { // Avoid short words like "king", "sing"
		return fmt.Errorf("%w: The first word '%s' in your commit message subject appears to be in the present participle (gerund) form. Fix: Use the imperative mood (command form) for the subject line. Example: 'Feat: Add user registration', not 'Feat: Adding user registration'.", ErrCommitMessageSubjectNotImperative, firstWord)
	}

	return nil
}
