package git

import (
	"testing"
)

func TestValidateGitURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid https url",
			url:     "https://github.com/company/repo.git",
			wantErr: false,
		},
		{
			name:    "valid ssh url",
			url:     "git@github.com:company/repo.git",
			wantErr: false,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
		{
			name:    "insecure git protocol",
			url:     "git://github.com/company/repo.git",
			wantErr: true,
		},
		{
			name:    "http protocol",
			url:     "http://github.com/company/repo.git",
			wantErr: true,
		},
		{
			name:    "local path",
			url:     "/path/to/repo",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateGitURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateGitURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
