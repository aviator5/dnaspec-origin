# Implementation Tasks: Add Manifest Commands

## 1. Project Setup

- [x] 1.1 Initialize Go module (go.mod) with Go 1.25
- [x] 1.2 Add dependencies: cobra, yaml.v3, charmbracelet (huh, lipgloss)
- [x] 1.3 Create project directory structure following golang-standards/project-layout (cmd/dnaspec/, internal/)
- [x] 1.4 Set up cmd/dnaspec/main.go entry point

## 2. Core Data Structures

- [x] 2.1 Implement manifest schema types (Manifest, ManifestGuideline, ManifestPrompt) in internal/core/config/manifest.go
- [x] 2.2 Implement YAML marshaling/unmarshaling with yaml.v3 tags
- [x] 2.3 Add manifest template generation function in internal/core/config/templates.go
- [x] 2.4 Write unit tests for schema types

## 3. Validation Engine

- [x] 3.1 Create validation error types in internal/core/validate/errors.go
- [x] 3.2 Implement manifest structure validation (version, required fields) in internal/core/validate/validator.go
- [x] 3.3 Implement guideline validation (required fields, applicable_scenarios, duplicates)
- [x] 3.4 Implement prompt validation (required fields, duplicates)
- [x] 3.5 Implement cross-reference validation (guideline -> prompt references)
- [x] 3.6 Implement file existence validation with filepath checking
- [x] 3.7 Implement path security validation (relative paths, no traversal, directory prefix)
- [x] 3.8 Implement naming convention validation (spinal-case)
- [x] 3.9 Write comprehensive unit tests for all validation rules

## 4. CLI Commands

- [x] 4.1 Set up Cobra root command in internal/cli/root.go
- [x] 4.2 Create `manifest` command group in internal/cli/manifest/manifest.go
- [x] 4.3 Implement `manifest init` subcommand in internal/cli/manifest/init.go
- [x] 4.4 Add existence check to prevent overwriting in manifest init
- [x] 4.5 Add success message and next steps output with lipgloss formatting (internal/ui/styles.go)
- [x] 4.6 Implement `manifest validate` subcommand in internal/cli/manifest/validate.go
- [x] 4.7 Add error formatting and reporting for manifest validate using lipgloss
- [x] 4.8 Implement proper exit codes (0 for success, non-zero for errors)
- [x] 4.9 Add command help text and examples for both subcommands

## 5. Testing

- [x] 5.1 Create test fixtures (valid and invalid manifest examples)
- [x] 5.2 Write integration tests for manifest init command
- [x] 5.3 Write integration tests for manifest validate with various error cases
- [x] 5.4 Test path traversal security scenarios
- [x] 5.5 Test naming convention validation
- [x] 5.6 Test all error messages for clarity and accuracy

## 6. Documentation

- [x] 6.1 Add command documentation to README.md
- [x] 6.2 Create example manifest files in examples/ directory
- [x] 6.3 Document validation rules and common errors
- [x] 6.4 Add troubleshooting section for manifest validation errors

## Notes

- CI/CD integration postponed as requested
- Focus on correctness and clear error messages
- Use lipgloss for terminal output formatting
- Ensure atomic file writes for manifest init (write to temp, then rename)
- Commands organized as subcommands under `manifest` group for better CLI organization
