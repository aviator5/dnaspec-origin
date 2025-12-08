# Change: Add Validate and Sync Commands

## Why
Users need two essential project-level commands:
1. **Validate** - Verify project configuration (`dnaspec.yaml`) to catch errors before running operations
2. **Sync** - Convenience command to update all sources and regenerate agent files in a single non-interactive operation, ideal for CI/CD pipelines

These commands fill critical gaps in the workflow:
- `validate` provides early error detection for misconfigured projects
- `sync` simplifies the common pattern of "update everything and regenerate agents"

## What Changes
- Add `dnaspec validate` command to validate project configuration
  - Check YAML syntax and schema version
  - Verify all source file references exist in `dnaspec/` directory
  - Validate agent IDs are recognized
  - Warn about symlinked sources with missing paths
  - Display comprehensive error messages or success confirmation
- Add `dnaspec sync` command as convenience wrapper
  - Execute `dnaspec update --all` to update all sources
  - Execute `dnaspec update-agents --no-ask` to regenerate agent files
  - Provide consolidated output showing all changes
  - Non-interactive design for CI/CD compatibility
- Register both commands in main CLI entry point

## Impact
- **Affected specs**: project-management
- **Affected code**:
  - New file: `internal/cli/project/validate.go`
  - New file: `internal/cli/project/sync.go`
  - New tests: `internal/cli/project/validate_test.go`
  - New tests: `internal/cli/project/sync_test.go`
  - Modified: `cmd/dnaspec/main.go` (register commands)
  - Reuse: `internal/core/config` for configuration loading
  - Reuse: `internal/core/validate` package validation patterns
- **User experience**: 
  - Provides early error detection via `validate` before operations
  - Simplifies maintenance workflow with `sync` command
  - Matches design specification in `docs/design.md` (lines 840-941)
- **Breaking changes**: None - purely additive features
