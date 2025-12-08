# project-management Specification

## Purpose
TBD - created by archiving change add-project-commands. Update Purpose after archive.
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

