# Implementation Tasks: Add Remove Command

## 1. Command Structure

- [x] 1.1 Create internal/cli/project/remove.go with NewRemoveCmd() cobra command
- [x] 1.2 Define command structure with Use, Short, Long, and Example fields
- [x] 1.3 Add source name as required positional argument (Args: cobra.ExactArgs(1))
- [x] 1.4 Add --force/-f flag to skip confirmation prompt
- [x] 1.5 Add RunE function to execute remove command logic
- [x] 1.6 Verify file compiles without errors

## 2. Configuration Loading and Source Lookup

- [x] 2.1 Load project configuration using config.LoadProjectConfig("dnaspec.yaml")
- [x] 2.2 Handle missing configuration file with error message
- [x] 2.3 Suggest running "dnaspec init" when config not found
- [x] 2.4 Find source by name in config.Sources array
- [x] 2.5 Handle source not found error - display available source names
- [x] 2.6 Exit with appropriate error codes for error scenarios

## 3. Impact Display

- [x] 3.1 Implement displayImpact() function to show what will be deleted
- [x] 3.2 Count files in dnaspec/<source-name>/ directory (guidelines and prompts)
- [x] 3.3 Use filepath.Glob() to find generated Claude command files (.claude/commands/dnaspec/<source-name>-*.md)
- [x] 3.4 Use filepath.Glob() to find generated Copilot prompt files (.github/prompts/dnaspec-<source-name>-*.prompt.md)
- [x] 3.5 Display formatted list showing:
  - dnaspec.yaml entry
  - Source directory with file counts
  - Generated agent files with counts
- [x] 3.6 Handle missing directories gracefully (may not exist)
- [x] 3.7 Verify output matches design doc format (docs/design.md lines 660-668)

## 4. Confirmation Prompt

- [x] 4.1 Implement confirmation prompt "This cannot be undone. Continue? [y/N]"
- [x] 4.2 Read user input (y/Y for yes, anything else for no)
- [x] 4.3 Skip prompt if --force flag is set
- [x] 4.4 If user cancels, display cancellation message and exit with code 0
- [x] 4.5 Ensure confirmation happens BEFORE any file deletion

## 5. File Deletion - Generated Agent Files

- [x] 5.1 Implement deleteGeneratedFiles() function
- [x] 5.2 Delete Claude command files matching .claude/commands/dnaspec/<source-name>-*.md
- [x] 5.3 Delete Copilot prompt files matching .github/prompts/dnaspec-<source-name>-*.prompt.md
- [x] 5.4 Handle missing directories gracefully (skip if doesn't exist)
- [x] 5.5 Count deleted files for success message
- [x] 5.6 Handle file deletion errors with clear error messages

## 6. File Deletion - Source Directory

- [x] 6.1 Delete dnaspec/<source-name>/ directory using os.RemoveAll()
- [x] 6.2 Handle missing directory gracefully (idempotent operation)
- [x] 6.3 Handle permission errors with clear error message
- [x] 6.4 Exit with error if directory deletion fails

## 7. Configuration Update

- [x] 7.1 Filter out removed source from config.Sources array
- [x] 7.2 Use config.AtomicWriteProjectConfig() to safely update config
- [x] 7.3 Handle config write errors (critical - suggests manual cleanup)
- [x] 7.4 Verify config update happens LAST (after all file deletions)

## 8. Success Message and Next Steps

- [x] 8.1 Display success message using ui.SuccessStyle
- [x] 8.2 Include count of cleaned up files
- [x] 8.3 Suggest running 'dnaspec update-agents' to regenerate AGENTS.md
- [x] 8.4 Exit with code 0

## 9. Command Registration

- [x] 9.1 Register NewRemoveCmd() in cmd/dnaspec/main.go
- [x] 9.2 Add after other project commands (init, add, update, update-agents, list)
- [x] 9.3 Verify command appears in "dnaspec --help" output
- [x] 9.4 Verify command can be invoked with "dnaspec remove <source-name>"

## 10. Unit Tests

- [x] 10.1 Create internal/cli/project/remove_test.go
- [x] 10.2 Test successful removal with confirmation (simulate user input 'y')
- [x] 10.3 Test successful removal with --force flag (no confirmation)
- [x] 10.4 Test user cancellation (simulate user input 'n')
- [x] 10.5 Test source not found error
- [x] 10.6 Test missing configuration file error
- [x] 10.7 Test cleanup of source directory
- [x] 10.8 Test cleanup of generated Claude command files
- [x] 10.9 Test cleanup of generated Copilot prompt files
- [x] 10.10 Test config update removes source entry correctly
- [x] 10.11 Test idempotent behavior (source directory already deleted)
- [x] 10.12 Run tests with "go test ./internal/cli/project/" and verify all pass

## 11. Integration Tests

- [x] 11.1 Create test fixture with sample dnaspec.yaml and source files
- [x] 11.2 Test end-to-end: add source → remove source → verify cleanup
- [x] 11.3 Test removal with generated agent files present
- [x] 11.4 Test removal with no generated agent files
- [x] 11.5 Verify all files are properly deleted
- [x] 11.6 Verify configuration is updated correctly

## 12. Manual Testing

- [x] 12.1 Build binary: go build -o dnaspec ./cmd/dnaspec
- [x] 12.2 Create test project with dnaspec init
- [x] 12.3 Add a DNA source with dnaspec add
- [x] 12.4 Run dnaspec update-agents to generate agent files
- [x] 12.5 Test remove with confirmation (answer 'y')
- [x] 12.6 Verify source directory deleted
- [x] 12.7 Verify generated agent files deleted
- [x] 12.8 Verify dnaspec.yaml updated correctly
- [x] 12.9 Test remove with --force flag (no prompt)
- [x] 12.10 Test remove with cancellation (answer 'n')
- [x] 12.11 Test error cases (missing source, missing config)
- [x] 12.12 Verify all error messages are clear and helpful

## 13. Documentation

- [x] 13.1 Update README.md to document the remove command
- [x] 13.2 Add usage examples for `dnaspec remove`
- [x] 13.3 Document the --force flag
- [x] 13.4 Update command list/table if present

## Notes

- Follow existing command patterns from internal/cli/project/list.go and add.go
- Use config.LoadProjectConfig() and config.AtomicWriteProjectConfig()
- Use internal/ui styles for consistent terminal output
- Delete files in order: generated agent files → source directory → config entry
- Use filepath.Glob() for pattern matching generated files
- Handle missing directories gracefully (idempotent operations)
- Exit codes: 0 for success or cancellation, non-zero for errors
