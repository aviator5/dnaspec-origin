# Implementation Tasks: Add Manifest Commands

## 1. Project Setup

- [ ] 1.1 Initialize Go module (go.mod) with Go 1.25
- [ ] 1.2 Add dependencies: cobra, yaml.v3, charmbracelet (huh, lipgloss)
- [ ] 1.3 Create project directory structure following golang-standards/project-layout (cmd/dnaspec/, internal/)
- [ ] 1.4 Set up cmd/dnaspec/main.go entry point

## 2. Core Data Structures

- [ ] 2.1 Implement manifest schema types (Manifest, ManifestGuideline, ManifestPrompt) in internal/core/config/manifest.go
- [ ] 2.2 Implement YAML marshaling/unmarshaling with yaml.v3 tags
- [ ] 2.3 Add manifest template generation function in internal/core/config/templates.go
- [ ] 2.4 Write unit tests for schema types

## 3. Validation Engine

- [ ] 3.1 Create validation error types in internal/core/validate/errors.go
- [ ] 3.2 Implement manifest structure validation (version, required fields) in internal/core/validate/validator.go
- [ ] 3.3 Implement guideline validation (required fields, applicable_scenarios, duplicates)
- [ ] 3.4 Implement prompt validation (required fields, duplicates)
- [ ] 3.5 Implement cross-reference validation (guideline -> prompt references)
- [ ] 3.6 Implement file existence validation with filepath checking
- [ ] 3.7 Implement path security validation (relative paths, no traversal, directory prefix)
- [ ] 3.8 Implement naming convention validation (spinal-case)
- [ ] 3.9 Write comprehensive unit tests for all validation rules

## 4. CLI Commands

- [ ] 4.1 Set up Cobra root command in internal/cli/root.go
- [ ] 4.2 Create `manifest` command group in internal/cli/manifest/manifest.go
- [ ] 4.3 Implement `manifest init` subcommand in internal/cli/manifest/init.go
- [ ] 4.4 Add existence check to prevent overwriting in manifest init
- [ ] 4.5 Add success message and next steps output with lipgloss formatting (internal/ui/styles.go)
- [ ] 4.6 Implement `manifest validate` subcommand in internal/cli/manifest/validate.go
- [ ] 4.7 Add error formatting and reporting for manifest validate using lipgloss
- [ ] 4.8 Implement proper exit codes (0 for success, non-zero for errors)
- [ ] 4.9 Add command help text and examples for both subcommands

## 5. Testing

- [ ] 5.1 Create test fixtures (valid and invalid manifest examples)
- [ ] 5.2 Write integration tests for manifest init command
- [ ] 5.3 Write integration tests for manifest validate with various error cases
- [ ] 5.4 Test path traversal security scenarios
- [ ] 5.5 Test naming convention validation
- [ ] 5.6 Test all error messages for clarity and accuracy

## 6. Documentation

- [ ] 6.1 Add command documentation to README.md
- [ ] 6.2 Create example manifest files in examples/ directory
- [ ] 6.3 Document validation rules and common errors
- [ ] 6.4 Add troubleshooting section for manifest validation errors

## Notes

- CI/CD integration postponed as requested
- Focus on correctness and clear error messages
- Use lipgloss for terminal output formatting
- Ensure atomic file writes for manifest init (write to temp, then rename)
- Commands organized as subcommands under `manifest` group for better CLI organization
