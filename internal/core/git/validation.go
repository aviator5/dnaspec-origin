package git

import (
	"errors"
	"strings"
)

// ValidateGitURL validates that a git URL is secure and supported
// Only allows https:// and git@ (SSH) URLs
// Rejects insecure git:// protocol
func ValidateGitURL(url string) error {
	if url == "" {
		return errors.New("git URL cannot be empty")
	}

	// Only allow https://, git@ (SSH), and file:// (local)
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "git@") && !strings.HasPrefix(url, "file://") {
		return errors.New("only HTTPS, SSH (git@), and local (file://) URLs are supported")
	}

	// Reject insecure git:// protocol
	if strings.HasPrefix(url, "git://") {
		return errors.New("git:// protocol is not allowed (insecure)")
	}

	return nil
}
