# Change: Add List Command

## Why
Users need a way to view their current DNA configuration to understand what sources, guidelines, prompts, and AI agents are active in their project. This is essential before running update or sync operations and for general configuration verification.

## What Changes
- Add `dnaspec list` command that displays configured agents, sources, guidelines, and prompts from `dnaspec.yaml`
- Implement formatted output showing:
  - Configured AI agents (Phase 1: Claude Code, GitHub Copilot)
  - All DNA sources with type-specific metadata (URL/path, ref, commit)
  - Guidelines and prompts for each source with names and descriptions
- Add comprehensive error handling for missing or malformed configuration files
- Register new command in main CLI entry point

## Impact
- **Affected specs**: project-management
- **Affected code**:
  - New file: `internal/cli/project/list.go`
  - Modified: `cmd/dnaspec/main.go` (register command)
  - New tests: `internal/cli/project/list_test.go`
- **User experience**: Provides visibility into current configuration, matching design specification in `docs/design.md` (lines 792-836)
- **Breaking changes**: None - purely additive feature
