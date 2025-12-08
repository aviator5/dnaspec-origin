package config

import (
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

// DeriveSourceName derives a source name from a git URL or local path
func DeriveSourceName(gitURL, localPath string) string {
	var raw string
	if gitURL != "" {
		raw = extractRepoName(gitURL)
	} else {
		raw = filepath.Base(localPath)
	}
	return SanitizeName(raw)
}

// extractRepoName extracts the repository name from a git URL
// Examples:
//   https://github.com/company/dna-guidelines.git -> dna-guidelines
//   git@github.com:company/dna.git -> dna
//   https://github.com/company/dna -> dna
func extractRepoName(gitURL string) string {
	// Handle SSH URLs (git@github.com:company/repo.git)
	if strings.HasPrefix(gitURL, "git@") {
		// Split on : to get the path part
		parts := strings.Split(gitURL, ":")
		if len(parts) >= 2 {
			gitURL = parts[1]
		}
	} else {
		// Try to parse as HTTP/HTTPS URL
		if u, err := url.Parse(gitURL); err == nil {
			gitURL = u.Path
		}
	}

	// Remove leading slashes
	gitURL = strings.TrimPrefix(gitURL, "/")

	// Get the last path component
	parts := strings.Split(gitURL, "/")
	if len(parts) > 0 {
		name := parts[len(parts)-1]
		// Remove .git suffix if present
		name = strings.TrimSuffix(name, ".git")
		return name
	}

	return "unknown"
}

// SanitizeName converts a name to a valid source name
// - Lowercase
// - Replace non-alphanumeric with hyphens
// - Trim and collapse hyphens
func SanitizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace non-alphanumeric with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	name = re.ReplaceAllString(name, "-")

	// Trim hyphens from start and end
	name = strings.Trim(name, "-")

	// Collapse multiple consecutive hyphens
	re = regexp.MustCompile(`-+`)
	name = re.ReplaceAllString(name, "-")

	return name
}
