package files

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDetectManagedBlock(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectFound   bool
		expectStart   int
		expectEnd     int
	}{
		{
			name: "content with valid managed block",
			content: `Header content
<!-- DNASPEC:START -->
Managed content
<!-- DNASPEC:END -->
Footer content`,
			expectFound: true,
			expectStart: 15, // Index of start marker
			expectEnd:   74, // Index after END marker
		},
		{
			name:        "content without markers",
			content:     "Just some regular content",
			expectFound: false,
			expectStart: -1,
			expectEnd:   -1,
		},
		{
			name: "content with only start marker",
			content: `Header
<!-- DNASPEC:START -->
Content`,
			expectFound: false,
			expectStart: -1,
			expectEnd:   -1,
		},
		{
			name: "content with only end marker",
			content: `Content
<!-- DNASPEC:END -->
Footer`,
			expectFound: false,
			expectStart: -1,
			expectEnd:   -1,
		},
		{
			name: "markers in wrong order",
			content: `<!-- DNASPEC:END -->
Content
<!-- DNASPEC:START -->`,
			expectFound: false,
			expectStart: -1,
			expectEnd:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasBlock, startIdx, endIdx := DetectManagedBlock(tt.content)

			assert.Equal(t, tt.expectFound, hasBlock)
			assert.Equal(t, tt.expectStart, startIdx)
			assert.Equal(t, tt.expectEnd, endIdx)
		})
	}
}

func TestReplaceManagedBlock(t *testing.T) {
	tests := []struct {
		name        string
		original    string
		newBlock    string
		expectation func(t *testing.T, result string)
	}{
		{
			name: "replace existing managed block",
			original: `User content before
<!-- DNASPEC:START -->
Old managed content
<!-- DNASPEC:END -->
User content after`,
			newBlock: "New managed content",
			expectation: func(t *testing.T, result string) {
				assert.Contains(t, result, "User content before")
				assert.Contains(t, result, "New managed content")
				assert.Contains(t, result, "User content after")
				assert.NotContains(t, result, "Old managed content")
				assert.Contains(t, result, ManagedBlockStart)
				assert.Contains(t, result, ManagedBlockEnd)
			},
		},
		{
			name:     "append when no managed block exists",
			original: "Existing user content",
			newBlock: "New managed content",
			expectation: func(t *testing.T, result string) {
				assert.Contains(t, result, "Existing user content")
				assert.Contains(t, result, "New managed content")
				assert.Contains(t, result, ManagedBlockStart)
				assert.Contains(t, result, ManagedBlockEnd)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceManagedBlock(tt.original, tt.newBlock)
			tt.expectation(t, result)
		})
	}
}

func TestAppendManagedBlock(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		newBlock string
		expected string
	}{
		{
			name:     "append to empty content",
			content:  "",
			newBlock: "Managed content",
			expected: `<!-- DNASPEC:START -->
Managed content
<!-- DNASPEC:END -->
`,
		},
		{
			name:     "append to existing content",
			content:  "Existing content",
			newBlock: "Managed content",
			expected: `Existing content

<!-- DNASPEC:START -->
Managed content
<!-- DNASPEC:END -->
`,
		},
		{
			name:     "append to content with trailing newline",
			content:  "Existing content\n",
			newBlock: "Managed content",
			expected: `Existing content

<!-- DNASPEC:START -->
Managed content
<!-- DNASPEC:END -->
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendManagedBlock(tt.content, tt.newBlock)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateFileWithManagedBlock(t *testing.T) {
	newBlock := "Test managed content"
	result := CreateFileWithManagedBlock(newBlock)

	assert.Contains(t, result, "# DNASpec Agent Instructions")
	assert.Contains(t, result, "Test managed content")
	assert.Contains(t, result, ManagedBlockStart)
	assert.Contains(t, result, ManagedBlockEnd)

	// Should have proper structure
	lines := strings.Split(result, "\n")
	assert.True(t, len(lines) > 5, "Should have header + managed block")
}

func TestFormatManagedBlock(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:    "content without trailing newline",
			content: "Test content",
			expected: `<!-- DNASPEC:START -->
Test content
<!-- DNASPEC:END -->
`,
		},
		{
			name:    "content with trailing newline",
			content: "Test content\n",
			expected: `<!-- DNASPEC:START -->
Test content
<!-- DNASPEC:END -->
`,
		},
		{
			name:    "multiline content",
			content: "Line 1\nLine 2\nLine 3",
			expected: `<!-- DNASPEC:START -->
Line 1
Line 2
Line 3
<!-- DNASPEC:END -->
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatManagedBlock(tt.content)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestManagedBlockPreservation(t *testing.T) {
	// Test that user content is fully preserved
	original := `# My Custom Header

Some important user notes that should be preserved.

<!-- DNASPEC:START -->
Old generated content that will be replaced
<!-- DNASPEC:END -->

## Additional Notes

More user content at the bottom.`

	newBlock := "Updated generated content"
	result := ReplaceManagedBlock(original, newBlock)

	// Check preservation
	assert.Contains(t, result, "# My Custom Header")
	assert.Contains(t, result, "Some important user notes that should be preserved.")
	assert.Contains(t, result, "## Additional Notes")
	assert.Contains(t, result, "More user content at the bottom.")

	// Check replacement
	assert.Contains(t, result, "Updated generated content")
	assert.NotContains(t, result, "Old generated content that will be replaced")
}
