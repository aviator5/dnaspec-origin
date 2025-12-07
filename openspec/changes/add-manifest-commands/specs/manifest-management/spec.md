# Capability: Manifest Management

## ADDED Requirements

### Requirement: Manifest Initialization

The system SHALL provide a `dnaspec manifest init` command that creates a new `dnaspec-manifest.yaml` file with example structure and documentation.

#### Scenario: Create new manifest successfully

- **WHEN** user runs `dnaspec manifest init` in a directory without existing manifest
- **THEN** create `dnaspec-manifest.yaml` with version 1, example guidelines and prompts sections, and helpful comments
- **AND** display success message with next steps (manifest validate command)

#### Scenario: Prevent overwriting existing manifest

- **WHEN** user runs `dnaspec manifest init` in a directory with existing `dnaspec-manifest.yaml`
- **THEN** exit with error message indicating manifest already exists
- **AND** do not modify the existing file

### Requirement: Manifest Validation Command

The system SHALL provide a `dnaspec manifest validate` command that validates the manifest file structure, references, and constraints.

#### Scenario: Validate manifest successfully

- **WHEN** user runs `dnaspec manifest validate` with valid manifest
- **THEN** display success message "✓ Manifest is valid"
- **AND** exit with code 0

#### Scenario: Report validation errors

- **WHEN** user runs `dnaspec manifest validate` with invalid manifest
- **THEN** display all validation errors with clear descriptions
- **AND** display error count summary
- **AND** exit with non-zero code

#### Scenario: Manifest file not found

- **WHEN** user runs `dnaspec manifest validate` without `dnaspec-manifest.yaml` in current directory
- **THEN** exit with error message "dnaspec-manifest.yaml not found"
- **AND** exit with non-zero code

### Requirement: Manifest Schema Structure

The manifest file SHALL follow a defined YAML schema with version, guidelines, and prompts sections.

#### Scenario: Valid manifest structure

- **WHEN** manifest contains version 1, guidelines array, and prompts array
- **THEN** validation passes structural checks

#### Scenario: Missing version field

- **WHEN** manifest does not contain version field
- **THEN** validation fails with error "missing required field: version"

### Requirement: Guideline Validation

Each guideline entry SHALL have required fields (name, file, description, applicable_scenarios) and valid references.

#### Scenario: Valid guideline entry

- **WHEN** guideline has name, file, description, non-empty applicable_scenarios, and valid prompt references
- **THEN** guideline validation passes

#### Scenario: Missing required guideline fields

- **WHEN** guideline is missing name, file, or description
- **THEN** validation fails with error identifying missing field and guideline

#### Scenario: Empty applicable scenarios

- **WHEN** guideline has empty or missing applicable_scenarios array
- **THEN** validation fails with error "guideline '{name}' has empty applicable_scenarios (required for AGENTS.md)"

#### Scenario: Guideline file does not exist

- **WHEN** guideline references a file path that does not exist
- **THEN** validation fails with error "file not found: {path}"

#### Scenario: Invalid prompt reference

- **WHEN** guideline references a prompt name that is not defined in prompts section
- **THEN** validation fails with error "guideline '{name}' references non-existent prompt '{prompt}'"

#### Scenario: Duplicate guideline names

- **WHEN** multiple guidelines have the same name
- **THEN** validation fails with error "duplicate guideline name: {name}"

### Requirement: Prompt Validation

Each prompt entry SHALL have required fields (name, file, description) and the referenced file must exist.

#### Scenario: Valid prompt entry

- **WHEN** prompt has name, file, description, and the file exists
- **THEN** prompt validation passes

#### Scenario: Missing required prompt fields

- **WHEN** prompt is missing name, file, or description
- **THEN** validation fails with error identifying missing field and prompt

#### Scenario: Prompt file does not exist

- **WHEN** prompt references a file path that does not exist
- **THEN** validation fails with error "file not found: {path}"

#### Scenario: Duplicate prompt names

- **WHEN** multiple prompts have the same name
- **THEN** validation fails with error "duplicate prompt name: {name}"

### Requirement: Naming Convention Validation

Guideline and prompt names SHALL use spinal-case format (lowercase letters separated by hyphens).

#### Scenario: Valid spinal-case names

- **WHEN** names use lowercase letters and hyphens only (e.g., "go-style", "rest-api")
- **THEN** naming validation passes

#### Scenario: Invalid naming format

- **WHEN** names contain uppercase letters, underscores, or special characters
- **THEN** validation fails with error indicating invalid naming format and expected format

### Requirement: File Path Security Validation

All file paths in manifest SHALL be relative paths within guidelines/ or prompts/ directories to prevent path traversal attacks.

#### Scenario: Valid relative paths

- **WHEN** guideline file is "guidelines/go-style.md" and prompt file is "prompts/review.md"
- **THEN** path validation passes

#### Scenario: Absolute path rejection

- **WHEN** file path is absolute (e.g., "/etc/passwd")
- **THEN** validation fails with error "absolute paths not allowed: {path}"

#### Scenario: Path traversal rejection

- **WHEN** file path contains ".." (e.g., "../../../etc/passwd")
- **THEN** validation fails with error "path traversal not allowed: {path}"

#### Scenario: Invalid directory prefix

- **WHEN** file path does not start with "guidelines/" or "prompts/"
- **THEN** validation fails with error "path must be within guidelines/ or prompts/: {path}"
