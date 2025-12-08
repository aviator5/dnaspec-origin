package config

import (
	"testing"
)

func TestDeriveSourceName(t *testing.T) {
	tests := []struct {
		name      string
		gitURL    string
		localPath string
		want      string
	}{
		{
			name:   "https url with .git",
			gitURL: "https://github.com/company/dna-guidelines.git",
			want:   "dna-guidelines",
		},
		{
			name:   "https url without .git",
			gitURL: "https://github.com/company/dna",
			want:   "dna",
		},
		{
			name:   "ssh url",
			gitURL: "git@github.com:company/dna.git",
			want:   "dna",
		},
		{
			name:      "local path",
			localPath: "/Users/me/my-patterns",
			want:      "my-patterns",
		},
		{
			name:      "local path with spaces",
			localPath: "/Users/me/My DNA Patterns",
			want:      "my-dna-patterns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveSourceName(tt.gitURL, tt.localPath)
			if got != tt.want {
				t.Errorf("DeriveSourceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "lowercase conversion",
			input: "MySource",
			want:  "mysource",
		},
		{
			name:  "replace spaces with hyphens",
			input: "my source name",
			want:  "my-source-name",
		},
		{
			name:  "replace special chars with hyphens",
			input: "my@source#name!",
			want:  "my-source-name",
		},
		{
			name:  "collapse multiple hyphens",
			input: "my---source--name",
			want:  "my-source-name",
		},
		{
			name:  "trim hyphens",
			input: "-my-source-",
			want:  "my-source",
		},
		{
			name:  "complex case",
			input: "My DNA@Guidelines 2024!",
			want:  "my-dna-guidelines-2024",
		},
		{
			name:  "already clean",
			input: "my-source",
			want:  "my-source",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeName(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractRepoName(t *testing.T) {
	tests := []struct {
		name   string
		gitURL string
		want   string
	}{
		{
			name:   "github https",
			gitURL: "https://github.com/company/repo.git",
			want:   "repo",
		},
		{
			name:   "github https without .git",
			gitURL: "https://github.com/company/repo",
			want:   "repo",
		},
		{
			name:   "github ssh",
			gitURL: "git@github.com:company/repo.git",
			want:   "repo",
		},
		{
			name:   "gitlab https",
			gitURL: "https://gitlab.com/company/project/repo.git",
			want:   "repo",
		},
		{
			name:   "with hyphens",
			gitURL: "https://github.com/company/my-awesome-repo.git",
			want:   "my-awesome-repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractRepoName(tt.gitURL)
			if got != tt.want {
				t.Errorf("extractRepoName() = %v, want %v", got, tt.want)
			}
		})
	}
}
