# Implementation Tasks: Add Project Commands

## 1. Data Models and Configuration

- [x] 1.1 Define project configuration data structures (ProjectConfig, Source, Guideline, Prompt) in internal/core/config/project.go
- [x] 1.2 Implement YAML marshaling/unmarshaling with yaml.v3 tags for project config
- [x] 1.3 Create project configuration template with examples in internal/core/config/templates.go
- [x] 1.4 Implement LoadProjectConfig and SaveProjectConfig functions
- [x] 1.5 Implement source name derivation in internal/core/config/naming.go (DeriveSourceName, SanitizeName)
- [x] 1.6 Write unit tests for config loading/saving and name derivation

## 2. Init Command

- [x] 2.1 Create internal/cli/project/init.go with NewInitCmd() cobra command
- [x] 2.2 Implement check for existing dnaspec.yaml (prevent overwrite)
- [x] 2.3 Call CreateExampleProjectConfig to generate initial config file
- [x] 2.4 Add success message with next steps using lipgloss formatting
- [x] 2.5 Write unit tests for init command (success case and overwrite prevention)
- [x] 2.6 Register init command in cmd/dnaspec/main.go

## 3. Git Repository Operations

- [x] 3.1 Add go-git library dependency to go.mod (or use git command via exec)
- [x] 3.2 Implement git URL validation in internal/core/git/validation.go (reject git://, accept https:// and git@)
- [x] 3.3 Implement temporary directory management in internal/core/git/tempdir.go (unique PID+random ID, cleanup function)
- [x] 3.4 Implement git clone with shallow clone and timeout in internal/core/git/clone.go
- [x] 3.5 Write unit tests for URL validation and temp directory cleanup
- [x] 3.6 Write integration tests for git clone operation

## 4. Source Fetching

- [x] 4.1 Create internal/core/source/fetch.go with FetchGitSource function (clone, parse manifest, return cleanup)
- [x] 4.2 Implement FetchLocalSource function (validate path exists, parse manifest)
- [x] 4.3 Add manifest validation after fetching using existing manifest-management validation
- [x] 4.4 Add path security validation for all file references
- [x] 4.5 Write unit tests for both git and local source fetching
- [x] 4.6 Write integration tests with sample repositories and directories

## 5. Add Command Structure

- [x] 5.1 Create internal/cli/project/add.go with NewAddCmd() cobra command
- [x] 5.2 Define command flags: --git-repo, --git-ref, --name, --all, --guideline (repeatable), --dry-run
- [x] 5.3 Implement source type detection (git vs local path)
- [x] 5.4 Add flag validation (mutual exclusivity checks)
- [x] 5.5 Add command help text with usage examples

## 6. Guideline Selection

- [x] 6.1 Create internal/ui/selection.go with SelectGuidelines interactive function
- [x] 6.2 Display guideline checklist with name, description, and applicable scenarios
- [x] 6.3 Implement multi-select using huh or bubbletea
- [x] 6.4 Implement --all flag handler (select all guidelines)
- [x] 6.5 Implement --guideline flag handler (select specific by name, validate names exist)
- [x] 6.6 Handle user cancellation gracefully
- [x] 6.7 Write tests for selection logic

## 7. File Operations

- [x] 7.1 Create internal/core/files/copy.go with CopyGuidelineFiles function
- [x] 7.2 Implement directory creation for destination paths
- [x] 7.3 Preserve relative path structure from manifest (guidelines/, prompts/)
- [x] 7.4 Implement atomic config writes in internal/core/config/atomic.go (temp file + rename)
- [x] 7.5 Add error handling for disk space and permissions
- [x] 7.6 Implement rollback on partial copy failures
- [x] 7.7 Write unit tests for file copying and atomic writes

## 8. Configuration Management

- [x] 8.1 Create internal/core/config/update.go with AddSource function
- [x] 8.2 Implement duplicate source name detection
- [x] 8.3 Extract prompts referenced by selected guidelines
- [x] 8.4 Build source entry with metadata (type, url/path, ref, commit, guidelines, prompts)
- [x] 8.5 Append source to ProjectConfig sources array
- [x] 8.6 Save config atomically preserving existing sources
- [x] 8.7 Write unit tests for config update logic and round-trip (load, modify, save, load)

## 9. Add Command Integration

- [x] 9.1 Wire up add command: fetch source → select guidelines → copy files → update config
- [x] 9.2 Add defer cleanup for git temporary directories
- [x] 9.3 Implement error handling at each step with context
- [x] 9.4 Implement --dry-run mode (skip file writes, display preview)
- [x] 9.5 Add progress indicators for long operations (cloning, copying)
- [x] 9.6 Display success message with summary and next steps
- [x] 9.7 Register add command in cmd/dnaspec/main.go
- [x] 9.8 Write integration tests for full add workflow (git and local)

## 10. Error Handling and UX

- [x] 10.1 Review all error paths and add descriptive messages
- [x] 10.2 Add context to errors (URLs, paths, source names)
- [x] 10.3 Provide corrective action suggestions in error messages
- [x] 10.4 Use lipgloss for styled terminal output (success, errors, progress)
- [x] 10.5 Add timeout error handling for git operations
- [x] 10.6 Handle network failures gracefully
- [x] 10.7 Test error scenarios manually (duplicate names, missing manifest, invalid URLs, permissions)

## 11. Testing and Validation

- [x] 11.1 Achieve >80% code coverage for core logic (data models, git ops, file ops, config management)
- [x] 11.2 Write comprehensive unit tests for edge cases
- [x] 11.3 Create test fixtures (valid/invalid manifests, sample DNA repos)
- [x] 11.4 Write integration tests for full workflows
- [x] 11.5 Test security scenarios (path traversal, insecure URLs)
- [x] 11.6 Manual testing: init creates config, add with git repo, add with local path, interactive/non-interactive selection
- [x] 11.7 Manual testing: dry-run mode, error scenarios, cancellation

## 12. Documentation

- [x] 12.1 Update README with dnaspec init and dnaspec add command examples
- [x] 12.2 Document dnaspec.yaml configuration file structure
- [x] 12.3 Add troubleshooting section for common errors
- [x] 12.4 Create example project configurations
- [x] 12.5 Document flag usage and examples

## Notes

- Reuse existing manifest validation logic from manifest-management spec
- Use atomic file writes for all critical files (dnaspec.yaml)
- Namespace files under dnaspec/<source-name>/ to prevent conflicts
- Use shallow git clones (--depth=1) for efficiency
- Implement proper cleanup with defer for all temporary resources
- Clear, actionable error messages with context
- Display next steps after successful operations
