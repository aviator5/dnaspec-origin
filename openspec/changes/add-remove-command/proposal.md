# Change: Add Remove Command

## Why
Users need a way to remove DNA sources from their project configuration when they no longer need them. This is essential for:
- Cleaning up experimental or temporary DNA sources
- Removing deprecated guideline sources
- Maintaining a clean project configuration
- Reclaiming disk space from unused guideline files

Currently, users can add and update DNA sources, but there's no way to remove them without manually editing configuration files and deleting directories.

## What Changes
- Add `dnaspec remove <source-name>` command that safely removes a DNA source from the project
- Implement confirmation prompt (skippable with `--force` flag) to prevent accidental deletions
- Display detailed impact showing what will be deleted before confirmation
- Remove source entry from `dnaspec.yaml`
- Delete source directory (`dnaspec/<source-name>/`)
- Clean up generated agent files (Claude commands and Copilot prompts)
- Provide clear success message and suggest running `dnaspec update-agents` to regenerate AGENTS.md
- Add comprehensive error handling for missing sources and file operations

## Impact
- **Affected specs**: project-management
- **Affected code**:
  - New file: `internal/cli/project/remove.go`
  - Modified: `cmd/dnaspec/main.go` (register command)
  - New tests: `internal/cli/project/remove_test.go`
- **User experience**: Completes the source management lifecycle (init → add → update → remove), matching design specification in `docs/design.md` (lines 602-676)
- **Breaking changes**: None - purely additive feature
- **Dependencies**: Uses existing config management and file operations infrastructure
