# Project Management Spec Delta: Interactive Guideline Selection

## REMOVED Requirements

### ~~Requirement: Update All Sources~~

~~The system SHALL support updating all configured sources in a single operation.~~

**Rationale**: Removed to enforce one-source-at-a-time updates with fine-grained control. Users can iterate over sources in scripts if batch updates are needed.

#### ~~Scenario: Update all sources with --all flag~~

~~**WHEN** user runs `dnaspec update --all`~~
~~**THEN** iterate through all configured sources~~
~~**AND** fetch latest from each source origin~~
~~**AND** update each source following single-source update logic~~

#### ~~Scenario: Display progress for each source~~

~~**WHEN** updating all sources~~
~~**THEN** display section header for each source~~
~~**AND** show individual progress and results~~
~~**AND** continue even if one source fails~~

#### ~~Scenario: Report failures at end~~

~~**WHEN** updating all sources~~
~~**AND** one or more sources fail~~
~~**THEN** display error for each failed source~~
~~**AND** exit with error status~~
~~**AND** report total number of failures~~

### ~~Requirement: Mutual Exclusion of Arguments~~

~~The system SHALL enforce mutual exclusivity of `source-name` argument and `--all` flag.~~

**Rationale**: Removed because `--all` flag no longer exists. Source name is now always required.

#### ~~Scenario: Require source name or --all flag~~

~~**WHEN** user runs `dnaspec update` without arguments or flags~~
~~**THEN** exit with error message "must specify either source name or --all flag"~~

#### ~~Scenario: Reject both source name and --all flag~~

~~**WHEN** user runs `dnaspec update <source-name> --all`~~
~~**THEN** exit with error message "cannot specify both source name and --all flag"~~

### ~~Requirement: New Guidelines Handling~~

~~The system SHALL detect and handle new guidelines that appear in updated manifests.~~

**Rationale**: Replaced by interactive guideline selection. Users now explicitly select which guidelines to add, update, or keep through a multi-select UI.

#### ~~Scenario: Interactive prompt for new guidelines (default)~~

~~**WHEN** updated manifest contains new guidelines not in current configuration~~
~~**AND** user has not specified `--add-new` flag~~
~~**THEN** display list of new guidelines with descriptions~~
~~**AND** prompt user "Add new guidelines? [y/N]"~~
~~**AND** if user confirms, add all new guidelines~~
~~**AND** if user declines, skip new guidelines~~

#### ~~Scenario: Non-interactive add all new guidelines~~

~~**WHEN** user runs `dnaspec update <source-name> --add-new=all`~~
~~**AND** updated manifest contains new guidelines~~
~~**THEN** automatically add all new guidelines without prompting~~
~~**AND** display confirmation for each added guideline~~

#### ~~Scenario: Non-interactive skip new guidelines~~

~~**WHEN** user runs `dnaspec update <source-name> --add-new=none`~~
~~**AND** updated manifest contains new guidelines~~
~~**THEN** skip new guidelines without prompting~~
~~**AND** display count of skipped guidelines~~

#### ~~Scenario: No new guidelines available~~

~~**WHEN** updated manifest contains no new guidelines~~
~~**THEN** skip new guidelines handling~~
~~**AND** do not prompt user~~

### ~~Requirement: Add-New Policy Validation~~

~~The system SHALL validate the `--add-new` flag value if provided.~~

**Rationale**: Removed because `--add-new` flag no longer exists. Interactive selection replaces binary policy.

#### ~~Scenario: Accept valid add-new values~~

~~**WHEN** user specifies `--add-new=all` or `--add-new=none`~~
~~**THEN** accept and apply the specified policy~~

#### ~~Scenario: Reject invalid add-new values~~

~~**WHEN** user specifies `--add-new=invalid`~~
~~**THEN** exit with error message "--add-new must be 'all' or 'none'"~~

## MODIFIED Requirements

### Requirement: Update Single Source Command

The system SHALL update a single source from its origin with interactive guideline selection.

**Changes**: Added interactive selection requiring user to explicitly choose guidelines.

#### Scenario: Update specific source by name

- **WHEN** user runs `dnaspec update <source-name>`
- **THEN** validate source exists in current configuration
- **AND** fetch latest from source origin
- **AND** parse source manifest
- **AND** present interactive guideline selection
- **AND** update configuration with selected guidelines
- **AND** copy selected guideline and prompt files
- **AND** update commit hash (for git sources)

#### Scenario: Source not found

- **WHEN** user runs `dnaspec update <source-name>`
- **AND** source name does not exist in configuration
- **THEN** display error "Source not found: <source-name>"
- **AND** list all available source names
- **AND** exit with error status

#### Scenario: Update requires source name

- **WHEN** user runs `dnaspec update` without arguments
- **THEN** exit with error message "must specify a source name"

### Requirement: Update Progress Reporting and User Feedback

The system SHALL provide clear progress indicators and change summaries.

**Changes**: Removed reference to "Updated guidelines" section since selection is now interactive, not automatic.

#### Scenario: Display operation progress

- **WHEN** performing update operation
- **THEN** display "Fetching latest from <url>..." or "Refreshing from local directory..."
- **AND** show current and latest commit hashes for git sources

#### Scenario: Display selection results

- **WHEN** update operation completes
- **THEN** display summary of selected guidelines
- **AND** show count of added, updated, and removed guidelines

#### Scenario: Display next steps

- **WHEN** update operation completes successfully
- **THEN** display "Run 'dnaspec update-agents' to regenerate agent files"

#### Scenario: Clear error messages

- **WHEN** operation fails
- **THEN** display error message with specific cause
- **AND** suggest corrective action when applicable

## ADDED Requirements

### Requirement: Interactive Guideline Selection

The system SHALL present an interactive multi-select UI for choosing guidelines to add, update, or keep.

#### Scenario: Display available guidelines

- **WHEN** source is fetched successfully
- **AND** source manifest is parsed
- **THEN** display all guidelines from source manifest in multi-select UI
- **AND** show guideline name and description for each option
- **AND** allow space to select/deselect, enter to confirm

#### Scenario: Pre-select existing guidelines

- **WHEN** displaying guideline selection
- **AND** guideline exists in current project configuration
- **THEN** pre-check that guideline in the selection UI
- **AND** allow user to deselect if they want to remove it

#### Scenario: Display orphaned guidelines

- **WHEN** displaying guideline selection
- **AND** project configuration contains guideline not in source manifest
- **THEN** include orphaned guideline in selection list
- **AND** append warning icon (⚠️) to guideline label
- **AND** pre-check the orphaned guideline
- **AND** allow user to deselect to remove from config

#### Scenario: Apply guideline selection

- **WHEN** user confirms selection
- **THEN** update project configuration with only selected guidelines
- **AND** remove any previously configured guidelines that were not selected
- **AND** copy files for selected guidelines
- **AND** update associated prompts for selected guidelines

#### Scenario: Empty selection

- **WHEN** user confirms empty selection (no guidelines selected)
- **THEN** remove all guidelines from source in configuration
- **AND** preserve source entry with empty guidelines list
- **AND** display confirmation "All guidelines removed from <source-name>"

#### Scenario: Selection canceled

- **WHEN** user cancels the selection UI (Ctrl+C or similar)
- **THEN** exit without making changes
- **AND** display "Update canceled"
- **AND** exit with non-zero status

### Requirement: Guideline Categorization

The system SHALL categorize guidelines into available, existing, and orphaned states for selection.

#### Scenario: Categorize available guidelines

- **WHEN** comparing source manifest with current configuration
- **AND** guideline exists in manifest but not in configuration
- **THEN** categorize as "available"
- **AND** display normally in selection UI
- **AND** do not pre-check

#### Scenario: Categorize existing guidelines

- **WHEN** comparing source manifest with current configuration
- **AND** guideline exists in both manifest and configuration
- **THEN** categorize as "existing"
- **AND** display normally in selection UI
- **AND** pre-check by default

#### Scenario: Categorize orphaned guidelines

- **WHEN** comparing source manifest with current configuration
- **AND** guideline exists in configuration but not in manifest
- **THEN** categorize as "orphaned"
- **AND** display with warning icon (⚠️) in selection UI
- **AND** pre-check by default
- **AND** append to end of selection list

### Requirement: Dry Run with Interactive Flow

The system SHALL provide dry-run preview that shows guideline categorization without interactive selection.

#### Scenario: Dry run shows available guidelines

- **WHEN** user runs `dnaspec update <source-name> --dry-run`
- **THEN** fetch source and parse manifest
- **AND** categorize guidelines (available, existing, orphaned)
- **AND** display "Available guidelines:" with list of all manifest guidelines
- **AND** display "Already in config:" with checked guidelines
- **AND** display "Orphaned (in config but not in source):" with warning icon
- **AND** display "=== Dry Run - Preview ==="
- **AND** do not show interactive selection
- **AND** do not copy files or update configuration

## UNCHANGED Requirements

The following requirements from the existing project-management spec remain unchanged:

- **Requirement: Update Git Source Fetching** (all scenarios)
- **Requirement: Update Local Source Fetching** (all scenarios)
- **Requirement: Update Source Comparison** (all scenarios)
- **Requirement: Removed Guidelines Handling** (all scenarios)
- **Requirement: Update Configuration Updates** (all scenarios)
- **Requirement: Update File Operations** (all scenarios)
- **Requirement: Update Temporary Directory Cleanup** (all scenarios)
- **Requirement: Prompt File Updates** (all scenarios)
