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
				// Check for specific length information in the error message
				expectedLengthInfo := fmt.Sprintf("is %d characters long, but it should not exceed %d characters", len(strings.TrimSpace(tt.subject)), MaxSubjectLength)
				if !strings.Contains(err.Error(), expectedLengthInfo) {
					t.Errorf("CheckCommitMessageSubjectLength() error message \"%s\" did not contain expected length info \"%s\"", err.Error(), expectedLengthInfo)
				}
			}
		})
	}
}
