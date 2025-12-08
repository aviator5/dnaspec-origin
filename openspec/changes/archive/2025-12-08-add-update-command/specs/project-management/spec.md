# Spec Delta: Project Management - Update Command

## ADDED Requirements

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

### Requirement: Dry Run Mode

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

### Requirement: Configuration Updates

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

### Requirement: Progress Reporting and User Feedback

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

### Requirement: File Operations

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

### Requirement: Temporary Directory Cleanup

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

## MODIFIED Requirements

None. This change adds new functionality without modifying existing requirements.

## REMOVED Requirements

None. This change does not remove any existing requirements.
