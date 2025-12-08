# Design: Add Remove Command

## Overview
Implement `dnaspec remove <source-name>` command to safely remove DNA sources from project configuration, completing the source management lifecycle.

## Technical Approach

### Command Structure
Follow existing command patterns in `internal/cli/project/`:
- Use cobra.Command with positional source name argument
- Add `--force` flag to skip confirmation prompt
- Return errors for invalid inputs

### File Deletion Strategy
Delete files in this order to minimize inconsistency if operation fails:
1. Generated agent files (safest to delete first, easily regenerated)
2. Source directory (`dnaspec/<source-name>/`)
3. Config entry in `dnaspec.yaml` (atomic write, last step)

### Error Handling
- Source not found → clear error message listing available sources
- File deletion failures → report which files failed, suggest manual cleanup
- Config update failure → critical error, do not delete files first

### Confirmation Flow
- Without `--force`: Display impact, prompt "This cannot be undone. Continue? [y/N]"
- With `--force`: Skip confirmation, proceed directly to deletion
- User cancels (N) → exit gracefully with message

### Generated File Discovery
Find generated files by source name pattern:
- Claude commands: `.claude/commands/dnaspec/<source-name>-*.md`
- Copilot prompts: `.github/prompts/dnaspec-<source-name>-*.prompt.md`

Use filepath.Glob() for pattern matching to find all files.

### Config Update
Use existing `AtomicWriteProjectConfig()` to safely update config:
1. Read current config
2. Filter out source by name
3. Atomic write updated config

### Impact Display
Show user what will be deleted:
```
The following will be deleted:
  - dnaspec.yaml entry for '<source-name>'
  - dnaspec/<source-name>/ directory (X guidelines, Y prompts)
  - .claude/commands/dnaspec/<source-name>-*.md (N files)
  - .github/prompts/dnaspec-<source-name>-*.prompt.md (M files)
```

Count files/directories before deletion to provide accurate numbers.

## Implementation Considerations

### Similar to Existing Commands
- Follow patterns from `list.go`, `add.go`, `update.go`
- Use same error formatting with `ui.ErrorStyle`, `ui.SuccessStyle`
- Use same config loading with `config.LoadProjectConfig()`

### User Experience
- Clear, verbose output showing exactly what will be deleted
- Confirmation prompt prevents accidents
- Success message with next steps (run update-agents)
- Error messages suggest corrective actions

### Safety
- Require exact source name match (no wildcards)
- Confirmation prompt by default
- Atomic config update prevents partial state
- Clear error messages if deletion fails

## Dependencies
- Existing: `internal/core/config` for config management
- Existing: `internal/ui` for styled output
- Standard library: `os`, `path/filepath`, `fmt`
- cobra for command structure

## Testing Strategy
- Unit tests for each scenario (success, error cases, cancellation)
- Integration test: add source → remove source → verify cleanup
- Test both with and without --force flag
- Test error paths (missing source, file permission errors)
