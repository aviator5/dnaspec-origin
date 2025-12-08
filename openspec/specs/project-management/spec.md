# project-management Specification

## Purpose

Manages DNASpec project configuration, enabling users to initialize projects, add/update/remove DNA sources from git repositories or local directories, validate configuration, and synchronize all sources.

## Requirements
### Requirement: Project Configuration Structure

The system SHALL define a project configuration structure (`dnaspec.yaml`) that stores DNA sources, selected guidelines, and agent settings.

#### Scenario: Valid project configuration structure

- **WHEN** project configuration contains version 1, optional agents array, and sources array
- **THEN** configuration parsing succeeds
- **AND** configuration can be saved and loaded

#### Scenario: Configuration includes source metadata

- **WHEN** source entry contains name, type, location details, and selected guidelines
- **THEN** source metadata is properly structured
- **AND** metadata includes commit hash for git sources or path for local sources

### Requirement: Project Initialization Command

The system SHALL provide a `dnaspec init` command that creates a new `dnaspec.yaml` file with empty configuration and helpful examples.

#### Scenario: Create new project configuration successfully

- **WHEN** user runs `dnaspec init` in a directory without existing `dnaspec.yaml`
- **THEN** create `dnaspec.yaml` with version 1, empty sources array, and commented examples
- **AND** display success message with next steps (dnaspec add command)
- **AND** exit with code 0

#### Scenario: Prevent overwriting existing configuration

- **WHEN** user runs `dnaspec init` in a directory with existing `dnaspec.yaml`
- **THEN** exit with error message "dnaspec.yaml already exists"
- **AND** do not modify the existing file
- **AND** exit with non-zero code

### Requirement: Add DNA Source from Git Repository

The system SHALL provide ability to add DNA guidelines from a git repository using `dnaspec add --git-repo <url>` command.

#### Scenario: Add source from git repository successfully

- **WHEN** user runs `dnaspec add --git-repo <url>` with valid repository URL
- **THEN** clone repository to temporary directory
- **AND** parse `dnaspec-manifest.yaml` from repository root
- **AND** derive source name from repository URL
- **AND** display interactive guideline selection
- **AND** copy selected guideline and prompt files to `dnaspec/<source-name>/` directory
- **AND** update `dnaspec.yaml` with source metadata (type, url, ref, commit, guidelines, prompts)
- **AND** cleanup temporary directory
- **AND** display next steps (run dnaspec update-agents)

#### Scenario: Specify git reference

- **WHEN** user runs `dnaspec add --git-repo <url> --git-ref <tag-or-branch>`
- **THEN** clone repository at specified ref (tag or branch)
- **AND** record ref in source metadata

#### Scenario: Custom source name

- **WHEN** user runs `dnaspec add --git-repo <url> --name <custom-name>`
- **THEN** use custom name instead of derived name
- **AND** create files under `dnaspec/<custom-name>/` directory

#### Scenario: Duplicate source name

- **WHEN** user attempts to add source with name that already exists
- **THEN** exit with error message "source with name '<name>' already exists"
- **AND** suggest using --name flag with different name
- **AND** do not modify configuration

#### Scenario: Repository clone failure

- **WHEN** git clone operation fails (network, authentication, invalid URL)
- **THEN** exit with clear error message describing failure
- **AND** cleanup partial temporary directory
- **AND** do not modify configuration

#### Scenario: Missing manifest in repository

- **WHEN** cloned repository does not contain `dnaspec-manifest.yaml`
- **THEN** exit with error message "dnaspec-manifest.yaml not found in repository"
- **AND** cleanup temporary directory
- **AND** do not modify configuration

### Requirement: Add DNA Source from Local Directory

The system SHALL provide ability to add DNA guidelines from a local directory using `dnaspec add <local-path>` command.

#### Scenario: Add source from local directory successfully

- **WHEN** user runs `dnaspec add <local-path>` with valid directory path
- **THEN** read `dnaspec-manifest.yaml` from directory
- **AND** derive source name from directory name
- **AND** display interactive guideline selection
- **AND** copy selected guideline and prompt files to `dnaspec/<source-name>/` directory
- **AND** update `dnaspec.yaml` with source metadata (type, path, guidelines, prompts)
- **AND** display next steps

#### Scenario: Local directory does not exist

- **WHEN** user runs `dnaspec add <local-path>` with non-existent path
- **THEN** exit with error message "directory not found: <path>"
- **AND** do not modify configuration

#### Scenario: Missing manifest in local directory

- **WHEN** local directory does not contain `dnaspec-manifest.yaml`
- **THEN** exit with error message "dnaspec-manifest.yaml not found in <path>"
- **AND** do not modify configuration

### Requirement: Source Name Derivation

The system SHALL automatically derive source names from git repository URLs or local directory paths.

#### Scenario: Derive name from git URL

- **WHEN** git URL is "https://github.com/company/dna-guidelines.git"
- **THEN** derived source name is "dna-guidelines"

#### Scenario: Derive name from local path

- **WHEN** local path is "/Users/me/my-dna-patterns"
- **THEN** derived source name is "my-dna-patterns"

#### Scenario: Sanitize derived name

- **WHEN** derived name contains uppercase letters, spaces, or special characters
- **THEN** convert to lowercase
- **AND** replace non-alphanumeric characters with hyphens
- **AND** collapse consecutive hyphens
- **AND** trim leading and trailing hyphens

### Requirement: Interactive Guideline Selection

The system SHALL provide interactive selection of guidelines when adding a source without flags.

#### Scenario: Display guideline checklist

- **WHEN** adding source without --all or --guideline flags
- **THEN** display checklist of all available guidelines from manifest
- **AND** show guideline name, description, and applicable scenarios
- **AND** allow user to select multiple guidelines

#### Scenario: User selects specific guidelines

- **WHEN** user selects subset of available guidelines
- **THEN** copy only selected guidelines and their associated prompts
- **AND** record only selected guidelines in configuration

#### Scenario: User cancels selection

- **WHEN** user cancels guideline selection
- **THEN** abort operation without modifying configuration
- **AND** display cancellation message

### Requirement: Non-Interactive Guideline Selection

The system SHALL provide non-interactive modes for adding guidelines using flags.

#### Scenario: Add all guidelines

- **WHEN** user runs `dnaspec add --git-repo <url> --all`
- **THEN** automatically add all guidelines from source without prompting
- **AND** skip interactive selection

#### Scenario: Add specific guidelines by name

- **WHEN** user runs `dnaspec add --git-repo <url> --guideline go-style --guideline rest-api`
- **THEN** add only specified guidelines
- **AND** skip interactive selection

#### Scenario: Specified guideline does not exist

- **WHEN** user specifies guideline name not in manifest
- **THEN** exit with error message "guideline '<name>' not found in source"
- **AND** list available guideline names
- **AND** do not modify configuration

### Requirement: Dry Run Mode

The system SHALL provide a --dry-run flag that previews changes without modifying files.

#### Scenario: Dry run shows preview

- **WHEN** user runs `dnaspec add --git-repo <url> --dry-run`
- **THEN** perform all operations except writing files
- **AND** display what would be created/modified
- **AND** show selected guidelines and file paths
- **AND** do not create or modify any files

### Requirement: File Copying to Project Directory

The system SHALL copy selected guideline and prompt files to the project's `dnaspec/<source-name>/` directory.

#### Scenario: Copy files preserving structure

- **WHEN** guideline file is "guidelines/go-style.md" in source
- **THEN** copy to "dnaspec/<source-name>/guidelines/go-style.md" in project
- **AND** preserve relative directory structure from manifest

#### Scenario: Copy associated prompts

- **WHEN** selected guideline references prompts
- **THEN** copy all referenced prompt files to `dnaspec/<source-name>/prompts/`
- **AND** include prompts in source metadata

#### Scenario: Handle file copy errors

- **WHEN** file copy operation fails (permissions, disk space)
- **THEN** exit with clear error message
- **AND** do not modify configuration
- **AND** attempt cleanup of partial copies

### Requirement: Configuration Update

The system SHALL update `dnaspec.yaml` with source metadata after successfully adding a source.

#### Scenario: Add git source to configuration

- **WHEN** git source is successfully added
- **THEN** append source entry with type "git-repo", url, ref, commit hash, guidelines, and prompts
- **AND** write updated configuration atomically

#### Scenario: Add local source to configuration

- **WHEN** local source is successfully added
- **THEN** append source entry with type "local-dir", path, guidelines, and prompts
- **AND** write updated configuration atomically

#### Scenario: Preserve existing configuration

- **WHEN** updating configuration with new source
- **THEN** preserve all existing sources and settings
- **AND** append new source to sources array

### Requirement: Git Repository Security

The system SHALL validate git repository URLs and implement security measures for cloning.

#### Scenario: Reject insecure protocols

- **WHEN** git URL uses git:// protocol
- **THEN** exit with error message "only HTTPS and SSH URLs supported"
- **AND** do not attempt clone

#### Scenario: Accept secure protocols

- **WHEN** git URL uses https:// or git@ (SSH) protocol
- **THEN** proceed with clone operation

#### Scenario: Clone timeout

- **WHEN** git clone operation exceeds timeout threshold (5 minutes)
- **THEN** abort clone operation
- **AND** cleanup temporary directory
- **AND** exit with timeout error message

#### Scenario: Shallow clone for efficiency

- **WHEN** cloning git repository
- **THEN** use --depth=1 for shallow clone
- **AND** reduce bandwidth and storage requirements

### Requirement: Temporary Directory Management

The system SHALL manage temporary directories for git clone operations with proper cleanup.

#### Scenario: Create unique temporary directory

- **WHEN** cloning git repository
- **THEN** create unique temp directory using process ID and random identifier
- **AND** use system temp directory location (cross-platform)

#### Scenario: Cleanup on success

- **WHEN** add operation completes successfully
- **THEN** remove temporary directory and all contents

#### Scenario: Cleanup on failure

- **WHEN** add operation fails at any step
- **THEN** remove temporary directory and all contents
- **AND** do not leave orphaned temp directories

### Requirement: Manifest Validation During Add

The system SHALL validate fetched manifest before processing guidelines during add operation.

#### Scenario: Validate manifest structure

- **WHEN** manifest is fetched from source
- **THEN** validate using existing manifest validation logic
- **AND** exit with validation errors if invalid

#### Scenario: Path security validation

- **WHEN** processing manifest file paths
- **THEN** validate paths using existing security rules (no absolute paths, no traversal, within expected directories)
- **AND** reject manifest with invalid paths

### Requirement: Error Messages and User Feedback

The system SHALL provide clear, actionable error messages and progress feedback.

#### Scenario: Display operation progress

- **WHEN** performing long operations (cloning, copying files)
- **THEN** display progress indicators
- **AND** show current operation (e.g., "Cloning repository...", "Copying files...")

#### Scenario: Clear error messages

- **WHEN** operation fails
- **THEN** display error message with specific cause
- **AND** suggest corrective action when applicable
- **AND** include relevant context (URLs, paths, names)

#### Scenario: Success confirmation

- **WHEN** operation completes successfully
- **THEN** display success message with summary
- **AND** show what was created/modified
- **AND** display next steps

### Requirement: List Project Configuration Command

The system SHALL provide a `dnaspec list` command that displays all configured DNA sources, guidelines, prompts, and AI agents from the project configuration.

#### Scenario: Display full configuration successfully

- **WHEN** user runs `dnaspec list` in a project with valid `dnaspec.yaml`
- **THEN** display "Configured Agents" section showing all agent IDs or "None configured" if empty
- **AND** display "Sources:" header
- **AND** for each source, display source name with type in parentheses
- **AND** for git-repo sources, display URL, Ref, and Commit fields
- **AND** for local-dir sources, display Path field
- **AND** for each source, display "Guidelines:" section with indented list of guidelines (name: description)
- **AND** for each source, display "Prompts:" section with indented list of prompts (name: description)
- **AND** exit with code 0

#### Scenario: Display configuration with no agents

- **WHEN** user runs `dnaspec list` and configuration has empty or missing agents array
- **THEN** display "Configured Agents" section with message indicating no agents configured
- **AND** continue to display sources normally
- **AND** exit with code 0

#### Scenario: Display configuration with no sources

- **WHEN** user runs `dnaspec list` and configuration has empty sources array
- **THEN** display configured agents if any
- **AND** display "Sources:" header with message indicating no sources configured
- **AND** exit with code 0

#### Scenario: Display source with no guidelines

- **WHEN** source has empty guidelines array
- **THEN** display source metadata normally
- **AND** display "Guidelines:" section with message indicating no guidelines or empty list
- **AND** continue to display prompts if any

#### Scenario: Display source with no prompts

- **WHEN** source has empty prompts array
- **THEN** display source metadata and guidelines normally
- **AND** display "Prompts:" section with message indicating no prompts or empty list

#### Scenario: Configuration file not found

- **WHEN** user runs `dnaspec list` in directory without `dnaspec.yaml`
- **THEN** display error message indicating configuration file not found
- **AND** suggest running `dnaspec init` to create configuration
- **AND** exit with non-zero code

#### Scenario: Malformed configuration file

- **WHEN** user runs `dnaspec list` and `dnaspec.yaml` contains invalid YAML syntax
- **THEN** display error message with YAML parsing error details
- **AND** exit with non-zero code

#### Scenario: Display multiple sources of different types

- **WHEN** configuration contains both git-repo and local-dir sources
- **THEN** display each source with appropriate type-specific fields
- **AND** maintain consistent formatting across all sources
- **AND** display sources in order they appear in configuration

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

### Requirement: Update Single Source

The system SHALL provide a command to update a specific DNA source from its origin.

#### Scenario: Update git source with changes

- **WHEN** user runs `dnaspec update <source-name>` for a git-based source
- **AND** the source has new commits at the configured ref
- **THEN** clone repository at configured ref
- **AND** fetch latest manifest
- **AND** update selected guidelines with latest content and metadata
- **AND** copy updated files to `dnaspec/<source-name>/`
- **AND** update commit hash in `dnaspec.yaml`
- **AND** display summary of changes

#### Scenario: Update git source with no changes

- **WHEN** user runs `dnaspec update <source-name>` for a git-based source
- **AND** the source commit hash matches latest commit at configured ref
- **THEN** display "Source is up to date"
- **AND** do not modify any files or configuration

#### Scenario: Update local directory source

- **WHEN** user runs `dnaspec update <source-name>` for a local directory source
- **THEN** read manifest from configured path
- **AND** update selected guidelines with latest content and metadata
- **AND** copy updated files to `dnaspec/<source-name>/`
- **AND** do not check commit hash (not applicable for local sources)

#### Scenario: Source not found

- **WHEN** user runs `dnaspec update <source-name>` with non-existent source name
- **THEN** exit with error message "source '<name>' not found"
- **AND** list available source names
- **AND** do not modify any files or configuration

#### Scenario: Project not initialized

- **WHEN** user runs `dnaspec update` without `dnaspec.yaml` in current directory
- **THEN** exit with error message "No project configuration found"
- **AND** suggest running `dnaspec init` first

### Requirement: Update All Sources

The system SHALL provide a flag to update all configured DNA sources.

#### Scenario: Update all sources successfully

- **WHEN** user runs `dnaspec update --all`
- **THEN** iterate through all sources in configuration
- **AND** update each source using single-source update logic
- **AND** display progress for each source
- **AND** display final summary

#### Scenario: Update all sources with no sources configured

- **WHEN** user runs `dnaspec update --all`
- **AND** no sources are configured in `dnaspec.yaml`
- **THEN** display "No sources configured"
- **AND** do not modify configuration

#### Scenario: Update all sources with partial failures

- **WHEN** user runs `dnaspec update --all`
- **AND** some sources fail to update
- **THEN** continue updating remaining sources
- **AND** display error for each failed source
- **AND** exit with error code indicating partial failure

### Requirement: Guideline Comparison and Updates

The system SHALL compare current guidelines with latest manifest and categorize changes.

#### Scenario: Detect updated guidelines

- **WHEN** comparing current configuration with latest manifest
- **AND** guideline exists in both but has changed metadata
- **THEN** categorize as "updated"
- **AND** update description, applicable_scenarios, and prompts in configuration
- **AND** copy latest file content

#### Scenario: Detect unchanged guidelines

- **WHEN** comparing current configuration with latest manifest
- **AND** guideline exists in both with identical metadata
- **THEN** categorize as "unchanged"
- **AND** still copy file content (may have changed without metadata change)

#### Scenario: Metadata changes include description

- **WHEN** guideline description changes between current and latest
- **THEN** categorize as updated
- **AND** update description in configuration

#### Scenario: Metadata changes include applicable scenarios

- **WHEN** guideline applicable_scenarios changes between current and latest
- **THEN** categorize as updated
- **AND** update applicable_scenarios in configuration

#### Scenario: Metadata changes include prompts

- **WHEN** guideline prompts references change between current and latest
- **THEN** categorize as updated
- **AND** update prompts in configuration

### Requirement: New Guidelines Handling

The system SHALL detect and handle new guidelines that appear in updated manifests.

#### Scenario: Interactive prompt for new guidelines (default)

- **WHEN** updated manifest contains new guidelines not in current configuration
- **AND** user has not specified `--add-new` flag
- **THEN** display list of new guidelines with descriptions
- **AND** prompt user "Add new guidelines? [y/N]"
- **AND** if user confirms, add all new guidelines
- **AND** if user declines, skip new guidelines

#### Scenario: Non-interactive add all new guidelines

- **WHEN** user runs `dnaspec update <source-name> --add-new=all`
- **AND** updated manifest contains new guidelines
- **THEN** automatically add all new guidelines without prompting
- **AND** display confirmation for each added guideline

#### Scenario: Non-interactive skip new guidelines

- **WHEN** user runs `dnaspec update <source-name> --add-new=none`
- **AND** updated manifest contains new guidelines
- **THEN** skip new guidelines without prompting
- **AND** display count of skipped guidelines

#### Scenario: No new guidelines available

- **WHEN** updated manifest contains no new guidelines
- **THEN** skip new guidelines handling
- **AND** do not prompt user

### Requirement: Removed Guidelines Handling

The system SHALL detect and report guidelines removed from source manifest.

#### Scenario: Report removed guidelines

- **WHEN** current configuration contains guideline not in updated manifest
- **THEN** display warning "guideline '<name>' removed from source"
- **AND** keep guideline in configuration (do not auto-delete)
- **AND** keep files in `dnaspec/<source-name>/` directory

#### Scenario: Continue update despite removed guidelines

- **WHEN** some guidelines are removed from source
- **THEN** continue updating remaining guidelines
- **AND** do not fail entire update operation

### Requirement: Update Dry Run Mode

The system SHALL provide a dry-run mode that previews changes without writing files.

#### Scenario: Dry run shows preview

- **WHEN** user runs `dnaspec update <source-name> --dry-run`
- **THEN** fetch source and parse manifest
- **AND** compare with current configuration
- **AND** display summary of changes (updated, new, removed counts)
- **AND** do not copy files
- **AND** do not update configuration
- **AND** display "No changes made (dry run)"

#### Scenario: Dry run with --all flag

- **WHEN** user runs `dnaspec update --all --dry-run`
- **THEN** preview changes for each source
- **AND** do not modify any files or configuration

### Requirement: Update Configuration Updates

The system SHALL update project configuration with latest source metadata.

#### Scenario: Update git source commit hash

- **WHEN** git source is updated successfully
- **THEN** update `commit` field with latest commit hash
- **AND** preserve all other source fields (name, type, url, ref)

#### Scenario: Update guidelines metadata

- **WHEN** guidelines are updated
- **THEN** update description, applicable_scenarios, and prompts for each guideline
- **AND** preserve guideline names and file paths

#### Scenario: Preserve other sources

- **WHEN** updating specific source
- **THEN** preserve all other sources in configuration unchanged

#### Scenario: Atomic configuration write

- **WHEN** writing updated configuration
- **THEN** use atomic write operation (write to temp file, then rename)
- **AND** prevent corrupted configuration on interrupted write

### Requirement: Update Progress Reporting and User Feedback

The system SHALL provide clear progress indicators and change summaries.

#### Scenario: Display operation progress

- **WHEN** performing update operation
- **THEN** display "Fetching latest from <url>..." or "Refreshing from local directory..."
- **AND** show current and latest commit hashes for git sources

#### Scenario: Display change summary

- **WHEN** update operation completes
- **THEN** display summary sections:
  - "Updated guidelines:" with list of updated guideline names
  - "New guidelines available:" with list of new guideline names and descriptions
  - "Removed from source:" with list of removed guideline names

#### Scenario: Display next steps

- **WHEN** update operation completes successfully
- **THEN** display "Run 'dnaspec update-agents' to regenerate agent files"

#### Scenario: Clear error messages

- **WHEN** operation fails
- **THEN** display error message with specific cause
- **AND** suggest corrective action when applicable

### Requirement: Mutual Exclusion of Arguments

The system SHALL enforce mutual exclusivity of `source-name` argument and `--all` flag.

#### Scenario: Require source name or --all flag

- **WHEN** user runs `dnaspec update` without arguments or flags
- **THEN** exit with error message "must specify either source name or --all flag"

#### Scenario: Reject both source name and --all flag

- **WHEN** user runs `dnaspec update <source-name> --all`
- **THEN** exit with error message "cannot specify both source name and --all flag"

### Requirement: Add-New Policy Validation

The system SHALL validate the `--add-new` flag value if provided.

#### Scenario: Accept valid add-new values

- **WHEN** user specifies `--add-new=all` or `--add-new=none`
- **THEN** accept and apply the specified policy

#### Scenario: Reject invalid add-new values

- **WHEN** user specifies `--add-new=invalid`
- **THEN** exit with error message "--add-new must be 'all' or 'none'"

### Requirement: Update File Operations

The system SHALL copy updated guideline and prompt files to project directory.

#### Scenario: Copy updated files

- **WHEN** updating guidelines
- **THEN** copy guideline markdown files to `dnaspec/<source-name>/guidelines/`
- **AND** copy associated prompt files to `dnaspec/<source-name>/prompts/`
- **AND** overwrite existing files

#### Scenario: Preserve directory structure

- **WHEN** copying files
- **THEN** preserve relative paths from manifest
- **AND** create subdirectories as needed

### Requirement: Update Temporary Directory Cleanup

The system SHALL properly clean up temporary directories for git sources.

#### Scenario: Cleanup on success

- **WHEN** git source update completes successfully
- **THEN** remove temporary clone directory and all contents

#### Scenario: Cleanup on failure

- **WHEN** git source update fails at any step
- **THEN** remove temporary clone directory and all contents
- **AND** do not leave orphaned temp directories

### Requirement: Prompt File Updates

The system SHALL update associated prompt files when guidelines are updated.

#### Scenario: Update prompts referenced by updated guidelines

- **WHEN** guideline is updated and references prompts
- **THEN** extract referenced prompts from latest manifest
- **AND** copy prompt files to `dnaspec/<source-name>/prompts/`
- **AND** update prompts list in source configuration

#### Scenario: Handle changed prompt references

- **WHEN** guideline adds or removes prompt references
- **THEN** update prompts list in configuration
- **AND** copy any newly referenced prompts

