# project-management Specification Delta

This delta adds requirements for project configuration validation and synchronization commands.

## ADDED Requirements

### Requirement: Project Configuration Validation Command

The system SHALL provide a `dnaspec validate` command that validates the project configuration (`dnaspec.yaml`) without modifying any files.

#### Scenario: Validate project configuration successfully

- **WHEN** user runs `dnaspec validate` with valid `dnaspec.yaml`
- **THEN** validate YAML syntax and schema structure
- **AND** verify config version is supported (currently version 1)
- **AND** check all sources have required fields
- **AND** verify all guideline file references exist in `dnaspec/` directory
- **AND** verify all prompt file references exist in `dnaspec/` directory
- **AND** validate agent IDs are recognized (claude-code, github-copilot)
- **AND** check for duplicate source names
- **AND** display success message "✓ Configuration is valid"
- **AND** show summary of validated sources and files
- **AND** exit with code 0

#### Scenario: Missing configuration file

- **WHEN** user runs `dnaspec validate` without existing `dnaspec.yaml`
- **THEN** display error message "dnaspec.yaml not found"
- **AND** suggest running `dnaspec init` first
- **AND** exit with non-zero code

#### Scenario: Invalid YAML syntax

- **WHEN** `dnaspec.yaml` contains invalid YAML syntax
- **THEN** display error message with parse failure details
- **AND** exit with non-zero code

#### Scenario: Unsupported config version

- **WHEN** `dnaspec.yaml` has unsupported version number
- **THEN** display error message "unsupported config version: \<version\>"
- **AND** indicate supported versions (currently version 1)
- **AND** exit with non-zero code

#### Scenario: Missing source fields

- **WHEN** source entry is missing required fields (name, type, type-specific fields)
- **THEN** display error message indicating missing field
- **AND** include source name in error context
- **AND** continue validation to report all errors
- **AND** exit with non-zero code after all checks

#### Scenario: Missing guideline file

- **WHEN** source references guideline file that doesn't exist
- **THEN** display error message "File not found: dnaspec/\<source\>/\<file\>"
- **AND** continue validation to report all missing files
- **AND** exit with non-zero code after all checks

#### Scenario: Missing prompt file

- **WHEN** source references prompt file that doesn't exist
- **THEN** display error message "File not found: dnaspec/\<source\>/\<file\>"
- **AND** continue validation to report all missing files
- **AND** exit with non-zero code after all checks

#### Scenario: Unknown agent ID

- **WHEN** config contains unrecognized agent ID
- **THEN** display error message "Unknown agent ID: '\<id\>'"
- **AND** list recognized agent IDs (claude-code, github-copilot)
- **AND** continue validation to report all errors
- **AND** exit with non-zero code after all checks

#### Scenario: Duplicate source names

- **WHEN** config contains multiple sources with same name
- **THEN** display error message "Duplicate source name: '\<name\>'"
- **AND** continue validation to report all duplicates
- **AND** exit with non-zero code after all checks

#### Scenario: Symlinked source with missing target path

- **WHEN** source has `symlinked: true` but target path doesn't exist
- **THEN** display warning message "Source '\<name\>' has symlinked=true but path '\<path\>' doesn't exist"
- **AND** continue validation (warning only, not error)
- **AND** if other errors exist, exit with non-zero code
- **AND** if only warnings, exit with code 0

#### Scenario: Multiple validation errors

- **WHEN** configuration has multiple validation errors
- **THEN** collect all errors before reporting
- **AND** display comprehensive error list showing all issues
- **AND** show error count in summary
- **AND** exit with non-zero code

### Requirement: Project Synchronization Command

The system SHALL provide a `dnaspec sync` command that updates all sources and regenerates agent files in a single non-interactive operation.

#### Scenario: Synchronize project successfully

- **WHEN** user runs `dnaspec sync`
- **THEN** display "Syncing all DNA sources..." header
- **AND** update all sources from their origins (equivalent to `dnaspec update --all`)
- **AND** use non-interactive mode (--add-new=none policy)
- **AND** display summary of source updates
- **AND** display "Regenerating agent files..." message
- **AND** regenerate all agent files (equivalent to `dnaspec update-agents --no-ask`)
- **AND** display summary of generated agent files
- **AND** display "✓ Sync complete" message
- **AND** exit with code 0

#### Scenario: Sync with dry-run preview

- **WHEN** user runs `dnaspec sync --dry-run`
- **THEN** preview what would be updated without writing files
- **AND** show which sources would be updated
- **AND** show which agent files would be regenerated
- **AND** do not modify any files
- **AND** exit with code 0

#### Scenario: No sources configured

- **WHEN** user runs `dnaspec sync` with no sources in configuration
- **THEN** display "No sources configured" message
- **AND** exit with code 0

#### Scenario: Source update fails

- **WHEN** updating a source fails during sync
- **THEN** display error message for failed source
- **AND** continue with remaining sources
- **AND** collect all errors
- **AND** skip agent regeneration if any source failed
- **AND** display error summary
- **AND** exit with non-zero code

#### Scenario: Agent regeneration fails

- **WHEN** agent file regeneration fails during sync
- **THEN** display error message with details
- **AND** exit with non-zero code

#### Scenario: Sync updates unchanged sources

- **WHEN** source is already at latest version
- **THEN** display "✓ No changes (already at latest commit)" for that source
- **AND** continue with remaining sources
- **AND** still regenerate agent files
- **AND** exit with code 0

#### Scenario: Non-interactive operation

- **WHEN** user runs `dnaspec sync`
- **THEN** do not prompt for any user input
- **AND** use saved agent configuration from `dnaspec.yaml`
- **AND** do not add new guidelines automatically (--add-new=none policy)
- **AND** suitable for CI/CD pipeline execution

### Requirement: Validation Error Reporting

The system SHALL provide comprehensive and actionable error reporting during validation.

#### Scenario: Display structured error list

- **WHEN** validation errors are found
- **THEN** display error count header
- **AND** list each error with bullet point
- **AND** include field/path context for each error
- **AND** use consistent formatting (ui.ErrorStyle, ui.CodeStyle)

#### Scenario: Display file reference list on success

- **WHEN** validation succeeds
- **THEN** display count of validated sources
- **AND** list all validated file paths
- **AND** show confirmation message
- **AND** use consistent formatting (ui.SuccessStyle)

### Requirement: Sync Operation Workflow

The system SHALL execute sync operations in correct sequence with proper error handling.

#### Scenario: Execute update before agent regeneration

- **WHEN** user runs `dnaspec sync`
- **THEN** update all sources first
- **AND** only regenerate agents after all sources updated successfully
- **AND** preserve atomic operation semantics

#### Scenario: Display consolidated summary

- **WHEN** sync operation completes
- **THEN** show count of updated sources
- **AND** show list of unchanged sources
- **AND** show count of regenerated agent files
- **AND** provide concise summary of all changes

#### Scenario: Preserve existing patterns

- **WHEN** implementing sync command
- **THEN** reuse existing update logic from `internal/cli/project/update.go`
- **AND** reuse existing agent update logic from `internal/cli/project/update_agents.go`
- **AND** maintain consistency with existing command patterns
