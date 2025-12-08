package files

import (
	"strings"
)

const (
	// ManagedBlockStart marks the beginning of DNASpec-managed content
	ManagedBlockStart = "<!-- DNASPEC:START -->"
	// ManagedBlockEnd marks the end of DNASpec-managed content
	ManagedBlockEnd = "<!-- DNASPEC:END -->"
)

// DetectManagedBlock checks if content contains managed block markers
// Returns hasBlock, startIdx, endIdx
func DetectManagedBlock(content string) (bool, int, int) {
	startIdx := strings.Index(content, ManagedBlockStart)
	if startIdx == -1 {
		return false, -1, -1
	}

	endIdx := strings.Index(content, ManagedBlockEnd)
	if endIdx == -1 {
		return false, -1, -1
	}

	// endIdx should be after startIdx
	if endIdx <= startIdx {
		return false, -1, -1
	}

	return true, startIdx, endIdx + len(ManagedBlockEnd)
}

// ReplaceManagedBlock replaces content between markers, preserving outside content
func ReplaceManagedBlock(content, newBlock string) string {
	hasBlock, startIdx, endIdx := DetectManagedBlock(content)
	if !hasBlock {
		// No managed block found, append at end
		return AppendManagedBlock(content, newBlock)
	}

	// Replace content between markers
	before := content[:startIdx]
	after := content[endIdx:]

	return before + formatManagedBlock(newBlock) + after
}

// AppendManagedBlock appends a managed block to existing content
func AppendManagedBlock(content, newBlock string) string {
	// Ensure content ends with newline
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// Add separator if content exists
	if content != "" {
		content += "\n"
	}

	return content + formatManagedBlock(newBlock)
}

// CreateFileWithManagedBlock creates new file content with header and managed block
func CreateFileWithManagedBlock(newBlock string) string {
	header := "# DNASpec Agent Instructions\n\n"
	header += "This file contains DNA (Development Norms & Architecture) guidelines for AI assistants.\n\n"

	return header + formatManagedBlock(newBlock)
}

// formatManagedBlock wraps content with managed block markers
func formatManagedBlock(content string) string {
	var sb strings.Builder

	sb.WriteString(ManagedBlockStart)
	sb.WriteString("\n")
	sb.WriteString(content)

	// Ensure content ends with newline before end marker
	if !strings.HasSuffix(content, "\n") {
		sb.WriteString("\n")
	}

	sb.WriteString(ManagedBlockEnd)
	sb.WriteString("\n")

	return sb.String()
}
