package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/aviator5/dnaspec/internal/core/config"
)

// spinalCaseRegex matches valid spinal-case names (lowercase letters and hyphens)
var spinalCaseRegex = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

// ValidateManifest validates a manifest and returns all validation errors
func ValidateManifest(manifest *config.Manifest, baseDir string) ValidationErrors {
	var errors ValidationErrors

	// Validate version
	if manifest.Version == 0 {
		errors.Add("version", "missing required field: version")
	}

	// Validate guidelines
	guidelineNames := make(map[string]bool)
	for i, guideline := range manifest.Guidelines {
		prefix := fmt.Sprintf("guidelines[%d]", i)
		errors = append(errors, validateGuideline(guideline, prefix, baseDir, guidelineNames)...)
	}

	// Validate prompts
	promptNames := make(map[string]bool)
	for i, prompt := range manifest.Prompts {
		prefix := fmt.Sprintf("prompts[%d]", i)
		errors = append(errors, validatePrompt(prompt, prefix, baseDir, promptNames)...)
	}

	// Validate cross-references (guideline prompts must exist)
	for i, guideline := range manifest.Guidelines {
		for _, promptName := range guideline.Prompts {
			if !promptNames[promptName] {
				errors.Add(
					fmt.Sprintf("guidelines[%d].prompts", i),
					fmt.Sprintf("guideline '%s' references non-existent prompt '%s'", guideline.Name, promptName),
				)
			}
		}
	}

	return errors
}

// validateGuideline validates a single guideline entry
func validateGuideline(g config.ManifestGuideline, prefix string, baseDir string, seenNames map[string]bool) ValidationErrors {
	var errors ValidationErrors

	// Check required fields
	if g.Name == "" {
		errors.Add(prefix+".name", "missing required field: name")
	} else {
		// Check for duplicates
		if seenNames[g.Name] {
			errors.Add(prefix+".name", fmt.Sprintf("duplicate guideline name: %s", g.Name))
		}
		seenNames[g.Name] = true

		// Validate naming convention
		if !spinalCaseRegex.MatchString(g.Name) {
			errors.Add(
				prefix+".name",
				fmt.Sprintf("invalid naming format: '%s' (expected spinal-case: lowercase letters and hyphens only)", g.Name),
			)
		}
	}

	if g.File == "" {
		errors.Add(prefix+".file", "missing required field: file")
	} else {
		// Validate file path security and existence
		errors = append(errors, validateFilePath(g.File, prefix+".file", baseDir, "guidelines/")...)
	}

	if g.Description == "" {
		errors.Add(prefix+".description", "missing required field: description")
	}

	if len(g.ApplicableScenarios) == 0 {
		errors.Add(
			prefix+".applicable_scenarios",
			fmt.Sprintf("guideline '%s' has empty applicable_scenarios (required for AGENTS.md)", g.Name),
		)
	}

	return errors
}

// validatePrompt validates a single prompt entry
func validatePrompt(p config.ManifestPrompt, prefix string, baseDir string, seenNames map[string]bool) ValidationErrors {
	var errors ValidationErrors

	// Check required fields
	if p.Name == "" {
		errors.Add(prefix+".name", "missing required field: name")
	} else {
		// Check for duplicates
		if seenNames[p.Name] {
			errors.Add(prefix+".name", fmt.Sprintf("duplicate prompt name: %s", p.Name))
		}
		seenNames[p.Name] = true

		// Validate naming convention
		if !spinalCaseRegex.MatchString(p.Name) {
			errors.Add(
				prefix+".name",
				fmt.Sprintf("invalid naming format: '%s' (expected spinal-case: lowercase letters and hyphens only)", p.Name),
			)
		}
	}

	if p.File == "" {
		errors.Add(prefix+".file", "missing required field: file")
	} else {
		// Validate file path security and existence
		errors = append(errors, validateFilePath(p.File, prefix+".file", baseDir, "prompts/")...)
	}

	if p.Description == "" {
		errors.Add(prefix+".description", "missing required field: description")
	}

	return errors
}

// validateFilePath validates a file path for security and existence
func validateFilePath(path, field, baseDir, expectedPrefix string) ValidationErrors {
	var errors ValidationErrors

	// Check for absolute paths
	if filepath.IsAbs(path) {
		errors.Add(field, fmt.Sprintf("absolute paths not allowed: %s", path))
		return errors
	}

	// Check for path traversal
	if strings.Contains(path, "..") {
		errors.Add(field, fmt.Sprintf("path traversal not allowed: %s", path))
		return errors
	}

	// Check directory prefix
	if !strings.HasPrefix(path, expectedPrefix) {
		errors.Add(field, fmt.Sprintf("path must be within %s: %s", expectedPrefix, path))
		return errors
	}

	// Check if file exists
	fullPath := filepath.Join(baseDir, path)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		errors.Add(field, fmt.Sprintf("file not found: %s", path))
	}

	return errors
}
