# project-management Specification Delta

## ADDED Requirements

### Requirement: Remove DNA Source Command

The system SHALL provide a `dnaspec remove <source-name>` command that safely removes a DNA source from the project configuration.

#### Scenario: Remove source successfully with confirmation

- **WHEN** user runs `dnaspec remove <source-name>` with valid source name
- **THEN** display impact showing what will be deleted:
  - Source entry in dnaspec.yaml
  - Files in dnaspec/<source-name>/ directory with count
  - Generated Claude command files matching pattern
  - Generated Copilot prompt files matching pattern
- **AND** prompt user for confirmation "This cannot be undone. Continue? [y/N]"
- **AND** if user confirms with 'y':
  - Delete generated agent files (Claude commands and Copilot prompts)
  - Delete source directory dnaspec/<source-name>/
  - Remove source entry from dnaspec.yaml using atomic write
  - Display success message
  - Suggest running 'dnaspec update-agents' to regenerate AGENTS.md
- **AND** exit with code 0

#### Scenario: Remove source with force flag

- **WHEN** user runs `dnaspec remove <source-name> --force`
- **THEN** skip confirmation prompt
- **AND** proceed directly to deletion
- **AND** perform all deletion steps
- **AND** display success message
- **AND** exit with code 0

#### Scenario: User cancels removal

- **WHEN** user runs `dnaspec remove <source-name>` without --force flag
- **AND** user responds 'n' or 'N' to confirmation prompt
- **THEN** display cancellation message
- **AND** do not delete any files
- **AND** do not modify configuration
- **AND** exit with code 0

#### Scenario: Source not found

- **WHEN** user runs `dnaspec remove <source-name>` with non-existent source name
- **THEN** display error message "Source '<source-name>' not found"
- **AND** list available source names from configuration
- **AND** do not modify any files
- **AND** exit with non-zero code

#### Scenario: Project configuration not found

- **WHEN** user runs `dnaspec remove` in directory without dnaspec.yaml
- **THEN** display error message "Project configuration not found: dnaspec.yaml"
- **AND** suggest running 'dnaspec init' to create configuration
- **AND** exit with non-zero code

#### Scenario: Clean up generated agent files

- **WHEN** removing source with prompts that generated agent files
- **THEN** discover and delete all generated files matching patterns:
  - `.claude/commands/dnaspec/<source-name>-*.md`
  - `.github/prompts/dnaspec-<source-name>-*.prompt.md`
- **AND** handle missing directories gracefully (directories may not exist if agents not configured)
- **AND** count deleted files and include in success message

#### Scenario: Source directory does not exist

- **WHEN** removing source and dnaspec/<source-name>/ directory does not exist
- **THEN** proceed with removal (idempotent operation)
- **AND** remove source entry from configuration
- **AND** clean up any generated agent files that exist
- **AND** display success message noting directory was already removed

#### Scenario: File deletion failure

- **WHEN** file deletion fails due to permissions or other I/O error
- **THEN** display error message with specific file path that failed
- **AND** suggest manual cleanup
- **AND** exit with non-zero code
- **AND** leave configuration unchanged to maintain consistency

### Requirement: Remove Command Arguments and Flags

The system SHALL accept a source name as required positional argument and optional --force flag.

#### Scenario: Missing source name argument

- **WHEN** user runs `dnaspec remove` without source name
- **THEN** display usage help with command syntax
- **AND** show example: "dnaspec remove <source-name>"
- **AND** exit with non-zero code

#### Scenario: Force flag skips confirmation

- **WHEN** user provides --force or -f flag
- **THEN** skip confirmation prompt
- **AND** proceed directly to deletion
