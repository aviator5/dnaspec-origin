# Implementation Tasks: Add Update Command

## 1. Core Comparison and Update Logic

- [ ] 1.1 Create internal/core/config/compare.go with GuidelineComparison struct (Updated, New, Removed, Unchanged fields)
- [ ] 1.2 Implement CompareGuidelines(current, manifest) function to categorize guidelines
- [ ] 1.3 Implement hasChanges() helper to detect metadata changes (description, scenarios, prompts)
- [ ] 1.4 Add FindSourceByName(config, name) function in internal/core/config/config.go
- [ ] 1.5 Add UpdateSourceInConfig(config, sourceName, updatedSource) function in internal/core/config/update.go
- [ ] 1.6 Write unit tests for comparison logic (all categorization scenarios)
- [ ] 1.7 Write unit tests for source lookup and update functions

## 2. Update Command Structure

- [ ] 2.1 Create internal/cli/project/update.go with NewUpdateCmd() cobra command
- [ ] 2.2 Define updateFlags struct with all, dryRun, addNew fields
- [ ] 2.3 Add command flags: --all, --dry-run, --add-new (with validation for "all" or "none")
- [ ] 2.4 Implement flag validation (require source name OR --all, not both)
- [ ] 2.5 Add command help text with usage examples matching docs/design.md
- [ ] 2.6 Register update command in cmd/dnaspec/main.go

## 3. Single Source Update Logic

- [ ] 3.1 Implement runUpdate() function with argument and flag validation
- [ ] 3.2 Load project configuration and check project is initialized
- [ ] 3.3 Implement updateSingleSource(config, sourceName, flags) function
- [ ] 3.4 Handle git sources: fetch latest, compare commits, display commit change
- [ ] 3.5 Handle local sources: fetch from path, display refresh message
- [ ] 3.6 Exit early if git source has no commit changes
- [ ] 3.7 Parse latest manifest and validate structure

## 4. Guideline Update Processing

- [ ] 4.1 Implement updateGuidelines() function using CompareGuidelines()
- [ ] 4.2 Update metadata for updated guidelines (description, applicable_scenarios, prompts)
- [ ] 4.3 Update metadata for unchanged guidelines (in case file content changed)
- [ ] 4.4 Copy files using existing files.CopyGuidelineFiles() for updated/unchanged guidelines
- [ ] 4.5 Display summary of updated guidelines with names
- [ ] 4.6 Collect and update associated prompts for updated guidelines

## 5. New Guidelines Handling

- [ ] 5.1 Implement handleNewGuidelines() function with policy parameter
- [ ] 5.2 Display list of new guidelines with names and descriptions
- [ ] 5.3 Implement interactive prompt "Add new guidelines? [y/N]" (default mode)
- [ ] 5.4 Implement --add-new=all policy (auto-add all new guidelines)
- [ ] 5.5 Implement --add-new=none policy (skip new guidelines)
- [ ] 5.6 Add new guidelines to configuration when accepted
- [ ] 5.7 Copy new guideline and prompt files to project directory
- [ ] 5.8 Display confirmation for each added guideline

## 6. Removed Guidelines Handling

- [ ] 6.1 Detect guidelines in config but not in latest manifest
- [ ] 6.2 Display warning for each removed guideline
- [ ] 6.3 Keep removed guidelines in configuration (do not auto-delete)
- [ ] 6.4 Keep files in dnaspec/<source-name>/ directory (do not auto-delete)
- [ ] 6.5 Continue update operation despite removed guidelines

## 7. Configuration and Metadata Updates

- [ ] 7.1 Update commit hash for git sources after successful update
- [ ] 7.2 Update all guideline metadata with latest values from manifest
- [ ] 7.3 Use UpdateSourceInConfig() to replace source in configuration
- [ ] 7.4 Use AtomicWriteProjectConfig() to save updated configuration
- [ ] 7.5 Preserve all other sources unchanged during update

## 8. Update All Sources

- [ ] 8.1 Implement updateAllSources(config, flags) function
- [ ] 8.2 Handle empty sources list (display "No sources configured")
- [ ] 8.3 Iterate through all sources and call updateSingleSource() for each
- [ ] 8.4 Display progress header for each source ("=== Updating <name> ===")
- [ ] 8.5 Continue on partial failures (collect errors but update remaining sources)
- [ ] 8.6 Display error for each failed source
- [ ] 8.7 Display final summary with success/failure counts
- [ ] 8.8 Exit with error code if any sources failed

## 9. Dry Run Mode

- [ ] 9.1 Add dry-run checks before all file copy operations
- [ ] 9.2 Add dry-run checks before all configuration write operations
- [ ] 9.3 Display preview summary showing updated/new/removed counts
- [ ] 9.4 Display "No changes made (dry run)" message
- [ ] 9.5 Ensure dry-run performs all read operations (fetch, parse, compare)

## 10. Progress Reporting and UX

- [ ] 10.1 Display "Fetching latest from <url>..." or "Refreshing from local directory..." messages
- [ ] 10.2 Display current and latest commit hashes for git sources
- [ ] 10.3 Format change summary sections (Updated guidelines / New guidelines / Removed)
- [ ] 10.4 Use lipgloss styles consistently (SuccessStyle, ErrorStyle, InfoStyle)
- [ ] 10.5 Display "Run 'dnaspec update-agents' to regenerate agent files" after success
- [ ] 10.6 Add clear error messages with context (source names, URLs, paths)
- [ ] 10.7 Suggest corrective actions in error messages where applicable

## 11. Error Handling

- [ ] 11.1 Handle source not found with list of available sources
- [ ] 11.2 Handle project not initialized with suggestion to run init
- [ ] 11.3 Handle git clone failures with descriptive messages
- [ ] 11.4 Handle manifest validation failures
- [ ] 11.5 Handle invalid --add-new values with error message
- [ ] 11.6 Handle mutual exclusivity violations (source name + --all)
- [ ] 11.7 Handle missing required arguments (no source name and no --all)
- [ ] 11.8 Ensure temporary directory cleanup on all error paths

## 12. Unit Tests

- [ ] 12.1 Test CompareGuidelines() with no changes scenario
- [ ] 12.2 Test CompareGuidelines() with updated guidelines
- [ ] 12.3 Test CompareGuidelines() with new guidelines
- [ ] 12.4 Test CompareGuidelines() with removed guidelines
- [ ] 12.5 Test CompareGuidelines() with mixed scenarios
- [ ] 12.6 Test hasChanges() for description changes
- [ ] 12.7 Test hasChanges() for applicable_scenarios changes
- [ ] 12.8 Test hasChanges() for prompts changes
- [ ] 12.9 Test FindSourceByName() for existing and non-existent sources
- [ ] 12.10 Test UpdateSourceInConfig() preserves other sources
- [ ] 12.11 Test UpdateSourceInConfig() updates specified source
- [ ] 12.12 Test UpdateSourceInConfig() with non-existent source name

## 13. Integration Tests

- [ ] 13.1 Create test fixtures with "old" and "new" manifests
- [ ] 13.2 Test updating git source with changes (commit hash changed)
- [ ] 13.3 Test updating git source with no changes (same commit)
- [ ] 13.4 Test updating local directory source
- [ ] 13.5 Test updating with new guidelines (interactive mode)
- [ ] 13.6 Test updating with --add-new=all
- [ ] 13.7 Test updating with --add-new=none
- [ ] 13.8 Test update --all with multiple sources
- [ ] 13.9 Test update --all with partial failures
- [ ] 13.10 Test dry-run mode (no files or config modified)
- [ ] 13.11 Test removed guidelines handling
- [ ] 13.12 Test configuration preservation (other sources unchanged)

## 14. Error Scenario Tests

- [ ] 14.1 Test source not found error
- [ ] 14.2 Test project not initialized error
- [ ] 14.3 Test invalid --add-new value error
- [ ] 14.4 Test mutual exclusivity error (source name + --all)
- [ ] 14.5 Test missing arguments error (no source name, no --all)
- [ ] 14.6 Test manifest validation failure during update

## 15. Documentation and Polish

- [ ] 15.1 Review and finalize command help text
- [ ] 15.2 Add comprehensive examples to help text
- [ ] 15.3 Update README.md with dnaspec update command examples
- [ ] 15.4 Add entry to CHANGELOG.md
- [ ] 15.5 Review and improve error messages for clarity
- [ ] 15.6 Ensure output formatting is consistent with other commands

## Notes

- Reuse existing source fetching logic from add command (FetchGitSource, FetchLocalSource)
- Reuse existing file copying logic from add command (CopyGuidelineFiles)
- Use atomic config writes for all configuration updates
- Always cleanup temporary directories with defer
- Default to interactive mode for new guidelines (explicit is better)
- Dry-run should perform all reads but no writes
- Update preserves files for removed guidelines (user can manually delete)
- No drift detection - always overwrite files with latest from source
