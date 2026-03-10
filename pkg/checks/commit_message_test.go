package checks

import (
	"fmt"
	"strings"
	"testing"
)

func TestCheckCommitMessageSubjectNotEmpty(t *testing.T) {
	tests := []struct {
		name        string
		subject     string
		expectedErr bool
		errMsgPart  string
	}{
		{
			name:        "Non-empty subject",
			subject:     "feat: Add new feature",
			expectedErr: false,
		},
		{
			name:        "Subject with only spaces",
			subject:     "   ",
			expectedErr: true,
			errMsgPart:  "commit message subject cannot be empty",
		},
		{
			name:        "Empty subject",
			subject:     "",
			expectedErr: true,
			errMsgPart:  "commit message subject cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCommitMessageSubjectNotEmpty(tt.subject)
			if (err != nil) != tt.expectedErr {
				t.Errorf("CheckCommitMessageSubjectNotEmpty() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if err != nil && tt.expectedErr {
				if !strings.Contains(err.Error(), tt.errMsgPart) {
					t.Errorf("CheckCommitMessageSubjectNotEmpty() error message \"%s\" did not contain \"%s\"", err.Error(), tt.errMsgPart)
				}
			}
		})
	}
}

func TestCheckCommitMessageSubjectLength(t *testing.T) {
	tests := []struct {
		name        string
		subject     string
		expectedErr bool
		errMsgPart  string
	}{
		{
			name:        "Subject within length limit",
			subject:     "feat: Optimize data fetching and reduce API calls", // 49 chars
			expectedErr: false,
		},
		{
			name:        "Subject exactly at length limit",
			subject:     "fix: Correct an issue where an entity was not save", // 50 chars
			expectedErr: false,
		},
		{
			name:        "Subject slightly over length limit",
			subject:     "refactor: Improve performance of data retrieval process X", // 55 chars
			expectedErr: true,
			errMsgPart:  "commit message subject is too long",
		},
		{
			name:        "Very long subject",
			subject:     "This is a very very very very very very very long commit message subject that definitely exceeds the fifty character limit",
			expectedErr: true,
			errMsgPart:  "commit message subject is too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCommitMessageSubjectLength(tt.subject)
			if (err != nil) != tt.expectedErr {
				t.Errorf("CheckCommitMessageSubjectLength() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if err != nil && tt.expectedErr {
				if !strings.Contains(err.Error(), tt.errMsgPart) {
					t.Errorf("CheckCommitMessageSubjectLength() error message \"%s\" did not contain \"%s\"", err.Error(), tt.errMsgPart)
				}
				expectedLengthInfo := fmt.Sprintf("is %d characters long, but it should not exceed %d characters", len(strings.TrimSpace(tt.subject)), MaxSubjectLength)
				if !strings.Contains(err.Error(), expectedLengthInfo) {
					t.Errorf("CheckCommitMessageSubjectLength() error message \"%s\" did not contain expected length info \"%s\"", err.Error(), expectedLengthInfo)
				}
			}
		})
	}
}

func TestCheckCommitMessageSubjectNoTrailingPeriod(t *testing.T) {
	tests := []struct {
		name        string
		subject     string
		expectedErr bool
		errMsgPart  string
	}{
		{
			name:        "Subject without trailing period",
			subject:     "feat: Add user authentication",
			expectedErr: false,
		},
		{
			name:        "Subject with trailing period",
			subject:     "fix: Resolve bug with login page.",
			expectedErr: true,
			errMsgPart:  "commit message subject should not end with a period",
		},
		{
			name:        "Subject ending with space and period",
			subject:     "docs: Update README .", // trailing space will be trimmed
			expectedErr: true,
			errMsgPart:  "commit message subject should not end with a period",
		},
		{
			name:        "Subject with period in middle",
			subject:     "feat: Version 1.0 release",
			expectedErr: false,
		},
		{
			name:        "Empty subject (handled by other check)",
			subject:     "",
			expectedErr: false, // This check specifically looks for a trailing period, not emptiness.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCommitMessageSubjectNoTrailingPeriod(tt.subject)
			if (err != nil) != tt.expectedErr {
				t.Errorf("CheckCommitMessageSubjectNoTrailingPeriod() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if err != nil && tt.expectedErr {
				if !strings.Contains(err.Error(), tt.errMsgPart) {
					t.Errorf("CheckCommitMessageSubjectNoTrailingPeriod() error message \"%s\" did not contain \"%s\"", err.Error(), tt.errMsgPart)
				}
			}
		})
	}
}

func TestCheckCommitMessageSubjectStartsWithUppercase(t *testing.T) {
	tests := []struct {
		name        string
		subject     string
		expectedErr bool
		errMsgPart  string
	}{
		{
			name:        "Starts with uppercase",
			subject:     "Feat: Add user authentication",
			expectedErr: false,
		},
		{
			name:        "Starts with lowercase",
			subject:     "feat: Add user authentication",
			expectedErr: true,
			errMsgPart:  "commit message subject should start with an uppercase letter",
		},
		{
			name:        "Starts with number",
			subject:     "1.0.0: Release new version",
			expectedErr: false,
		},
		{
			name:        "Starts with symbol",
			subject:     "#: Hotfix for critical bug",
			expectedErr: false,
		},
		{
			name:        "Empty subject (skipped)",	
			subject:     "",
			expectedErr: false, // This check returns nil if subject is empty, as it's handled by NotEmpty check
		},
		{
			name:        "Subject with leading space",
			subject:     " feat: leading space",
			expectedErr: true,
			errMsgPart:  "commit message subject should start with an uppercase letter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckCommitMessageSubjectStartsWithUppercase(tt.subject)
			if (err != nil) != tt.expectedErr {
				t.Errorf("CheckCommitMessageSubjectStartsWithUppercase() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if err != nil && tt.expectedErr {
				if !strings.Contains(err.Error(), tt.errMsgPart) {
					t.Errorf("CheckCommitMessageSubjectStartsWithUppercase() error message \"%s\" did not contain \"%s\"", err.Error(), tt.errMsgPart)
				}
			}
		})
	}
}
