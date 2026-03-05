package checks

import (
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
