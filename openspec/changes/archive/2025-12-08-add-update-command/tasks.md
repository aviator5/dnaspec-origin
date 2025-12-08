# Implementation Tasks: Add Update Command

## 1. Core Comparison and Update Logic

- [x] 1.1 Create internal/core/config/compare.go with GuidelineComparison struct (Updated, New, Removed, Unchanged fields)
- [x] 1.2 Implement CompareGuidelines(current, manifest) function to categorize guidelines
- [x] 1.3 Implement hasChanges() helper to detect metadata changes (description, scenarios, prompts)
- [x] 1.4 Add FindSourceByName(config, name) function in internal/core/config/config.go
- [x] 1.5 Add UpdateSourceInConfig(config, sourceName, updatedSource) function in internal/core/config/update.go
- [x] 1.6 Write unit tests for comparison logic (all categorization scenarios)
- [x] 1.7 Write unit tests for source lookup and update functions

## 2. Update Command Structure

- [x] 2.1 Create internal/cli/project/update.go with NewUpdateCmd() cobra command
- [x] 2.2 Define updateFlags struct with all, dryRun, addNew fields
- [x] 2.3 Add command flags: --all, --dry-run, --add-new (with validation for "all" or "none")
- [x] 2.4 Implement flag validation (require source name OR --all, not both)
- [x] 2.5 Add command help text with usage examples matching docs/design.md
- [x] 2.6 Register update command in cmd/dnaspec/main.go

## 3. Single Source Update Logic

- [x] 3.1 Implement runUpdate() function with argument and flag validation
- [x] 3.2 Load project configuration and check project is initialized
- [x] 3.3 Implement updateSingleSource(config, sourceName, flags) function
- [x] 3.4 Handle git sources: fetch latest, compare commits, display commit change
- [x] 3.5 Handle local sources: fetch from path, display refresh message
- [x] 3.6 Exit early if git source has no commit changes
- [x] 3.7 Parse latest manifest and validate structure

## 4. Guideline Update Processing

- [x] 4.1 Implement updateGuidelines() function using CompareGuidelines()
- [x] 4.2 Update metadata for updated guidelines (description, applicable_scenarios, prompts)
- [x] 4.3 Update metadata for unchanged guidelines (in case file content changed)
- [x] 4.4 Copy files using existing files.CopyGuidelineFiles() for updated/unchanged guidelines
- [x] 4.5 Display summary of updated guidelines with names
- [x] 4.6 Collect and update associated prompts for updated guidelines

## 5. New Guidelines Handling

- [x] 5.1 Implement handleNewGuidelines() function with policy parameter
- [x] 5.2 Display list of new guidelines with names and descriptions
- [x] 5.3 Implement interactive prompt "Add new guidelines? [y/N]" (default mode)
- [x] 5.4 Implement --add-new=all policy (auto-add all new guidelines)
- [x] 5.5 Implement --add-new=none policy (skip new guidelines)
- [x] 5.6 Add new guidelines to configuration when accepted
- [x] 5.7 Copy new guideline and prompt files to project directory
- [x] 5.8 Display confirmation for each added guideline

## 6. Removed Guidelines Handling

- [x] 6.1 Detect guidelines in config but not in latest manifest
- [x] 6.2 Display warning for each removed guideline
- [x] 6.3 Keep removed guidelines in configuration (do not auto-delete)
- [x] 6.4 Keep files in dnaspec/<source-name>/ directory (do not auto-delete)
- [x] 6.5 Continue update operation despite removed guidelines

## 7. Configuration and Metadata Updates

- [x] 7.1 Update commit hash for git sources after successful update
- [x] 7.2 Update all guideline metadata with latest values from manifest
- [x] 7.3 Use UpdateSourceInConfig() to replace source in configuration
- [x] 7.4 Use AtomicWriteProjectConfig() to save updated configuration
- [x] 7.5 Preserve all other sources unchanged during update

## 8. Update All Sources

- [x] 8.1 Implement updateAllSources(config, flags) function
- [x] 8.2 Handle empty sources list (display "No sources configured")
- [x] 8.3 Iterate through all sources and call updateSingleSource() for each
- [x] 8.4 Display progress header for each source ("=== Updating <name> ===")
- [x] 8.5 Continue on partial failures (collect errors but update remaining sources)
- [x] 8.6 Display error for each failed source
- [x] 8.7 Display final summary with success/failure counts
- [x] 8.8 Exit with error code if any sources failed

## 9. Dry Run Mode

- [x] 9.1 Add dry-run checks before all file copy operations
- [x] 9.2 Add dry-run checks before all configuration write operations
- [x] 9.3 Display preview summary showing updated/new/removed counts
- [x] 9.4 Display "No changes made (dry run)" message
- [x] 9.5 Ensure dry-run performs all read operations (fetch, parse, compare)

## 10. Progress Reporting and UX

- [x] 10.1 Display "Fetching latest from <url>..." or "Refreshing from local directory..." messages
- [x] 10.2 Display current and latest commit hashes for git sources
- [x] 10.3 Format change summary sections (Updated guidelines / New guidelines / Removed)
- [x] 10.4 Use lipgloss styles consistently (SuccessStyle, ErrorStyle, InfoStyle)
- [x] 10.5 Display "Run 'dnaspec update-agents' to regenerate agent files" after success
- [x] 10.6 Add clear error messages with context (source names, URLs, paths)
- [x] 10.7 Suggest corrective actions in error messages where applicable

## 11. Error Handling

- [x] 11.1 Handle source not found with list of available sources
- [x] 11.2 Handle project not initialized with suggestion to run init
- [x] 11.3 Handle git clone failures with descriptive messages
- [x] 11.4 Handle manifest validation failures
- [x] 11.5 Handle invalid --add-new values with error message
- [x] 11.6 Handle mutual exclusivity violations (source name + --all)
- [x] 11.7 Handle missing required arguments (no source name and no --all)
- [x] 11.8 Ensure temporary directory cleanup on all error paths

## 12. Unit Tests

- [x] 12.1 Test CompareGuidelines() with no changes scenario
- [x] 12.2 Test CompareGuidelines() with updated guidelines
- [x] 12.3 Test CompareGuidelines() with new guidelines
- [x] 12.4 Test CompareGuidelines() with removed guidelines
- [x] 12.5 Test CompareGuidelines() with mixed scenarios
- [x] 12.6 Test hasChanges() for description changes
- [x] 12.7 Test hasChanges() for applicable_scenarios changes
- [x] 12.8 Test hasChanges() for prompts changes
- [x] 12.9 Test FindSourceByName() for existing and non-existent sources
- [x] 12.10 Test UpdateSourceInConfig() preserves other sources
- [x] 12.11 Test UpdateSourceInConfig() updates specified source
- [x] 12.12 Test UpdateSourceInConfig() with non-existent source name

## 13. Integration Tests

- [x] 13.1 Create test fixtures with "old" and "new" manifests
- [x] 13.2 Test updating git source with changes (commit hash changed)
- [x] 13.3 Test updating git source with no changes (same commit)
- [x] 13.4 Test updating local directory source
- [x] 13.5 Test updating with new guidelines (interactive mode)
- [x] 13.6 Test updating with --add-new=all
- [x] 13.7 Test updating with --add-new=none
- [x] 13.8 Test update --all with multiple sources
- [x] 13.9 Test update --all with partial failures
- [x] 13.10 Test dry-run mode (no files or config modified)
- [x] 13.11 Test removed guidelines handling
- [x] 13.12 Test configuration preservation (other sources unchanged)

## 14. Error Scenario Tests

- [x] 14.1 Test source not found error
- [x] 14.2 Test project not initialized error
- [x] 14.3 Test invalid --add-new value error
- [x] 14.4 Test mutual exclusivity error (source name + --all)
- [x] 14.5 Test missing arguments error (no source name, no --all)
- [x] 14.6 Test manifest validation failure during update

## 15. Documentation and Polish

- [x] 15.1 Review and finalize command help text
- [x] 15.2 Add comprehensive examples to help text
- [x] 15.3 Update README.md with dnaspec update command examples
- [x] 15.4 Add entry to CHANGELOG.md
- [x] 15.5 Review and improve error messages for clarity
- [x] 15.6 Ensure output formatting is consistent with other commands

## Notes

- Reuse existing source fetching logic from add command (FetchGitSource, FetchLocalSource)
- Reuse existing file copying logic from add command (CopyGuidelineFiles)
- Use atomic config writes for all configuration updates
- Always cleanup temporary directories with defer
- Default to interactive mode for new guidelines (explicit is better)
- Dry-run should perform all reads but no writes
- Update preserves files for removed guidelines (user can manually delete)
- No drift detection - always overwrite files with latest from source
