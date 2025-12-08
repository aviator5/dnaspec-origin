# Proposal: Add Update-Agents Command

## Overview

Implement the `dnaspec update-agents` command that configures AI agent integrations and generates agent-specific files (AGENTS.md, CLAUDE.md, slash commands, prompts) based on selected DNA guidelines.

## Problem

After users add DNA guidelines with `dnaspec add`, they need a way to:

1. Select which AI agents (Claude Code, GitHub Copilot) to integrate with
2. Generate agent integration files that expose DNA guidelines to AI assistants
3. Regenerate these files when guidelines change or agents are added/removed

Without this command, DNA guidelines remain as files in `dnaspec/` directory but aren't surfaced to AI agents in a usable format. AI agents need:
- **Context-aware instructions** in AGENTS.md showing when to consult each guideline
- **Agent-specific command files** (Claude slash commands, Copilot prompts) for guideline-based code review
- **Managed blocks** that can be updated without destroying user customizations

## Solution

Implement `dnaspec update-agents` command with two modes:

### Interactive Mode (default)
1. Display checklist of available agents (Phase 1: Claude Code, GitHub Copilot)
2. Show currently selected agents as checked
3. Allow multi-select and save selection to `dnaspec.yaml`
4. Generate all agent files based on selection

### Non-Interactive Mode (`--no-ask` flag)
1. Use saved agent configuration from `dnaspec.yaml`
2. Generate all agent files without prompting
3. Useful for CI/CD and after adding/removing guidelines

### File Generation

**AGENTS.md** (always generated):
- Contains managed block with guideline references and applicable scenarios
- Format: `@/dnaspec/<source-name>/<file>` with bullet points for scenarios
- Preserves user content outside managed block markers

**CLAUDE.md** (if Claude Code selected):
- Same content as AGENTS.md
- Enables Claude Code to discover project instructions

**Claude Commands** (if Claude Code selected):
- `.claude/commands/dnaspec/<source-name>-<prompt-name>.md`
- Frontmatter with name, description, category, tags
- Prompt content wrapped in managed block
- Namespaced by source to prevent collisions

**Copilot Prompts** (if GitHub Copilot selected):
- `.github/prompts/dnaspec-<source-name>-<prompt-name>.prompt.md`
- Frontmatter with description
- Includes `$ARGUMENTS` placeholder
- Prompt content wrapped in managed block

## Scope

### In Scope
- Agent selection UI (interactive multi-select)
- Agent configuration persistence (agents array in dnaspec.yaml)
- Managed block generation algorithm using applicable_scenarios
- AGENTS.md file generation with create/append/replace logic
- CLAUDE.md file generation (same as AGENTS.md)
- Claude command file generation (frontmatter + content)
- Copilot prompt file generation (frontmatter + content)
- Managed block markers (`<!-- DNASPEC:START/END -->`)
- Directory creation for agent files
- File content preservation outside managed blocks
- `--no-ask` flag for non-interactive mode
- Source namespacing in filenames

### Out of Scope
- Future agent support (Windsurf, Cursor, Antigravity) - extensible design only
- Automatic regeneration on guideline changes (manual `update-agents` required)
- Selective prompt generation (all prompts for selected agents generated)
- Managed block migration or versioning
- Agent file cleanup when agents deselected (future: detect and offer to remove)

## Dependencies

- Existing project configuration (dnaspec.yaml) from add-project-commands change
- ProjectConfig, ProjectSource, ProjectGuideline, ProjectPrompt data structures
- Guideline files in dnaspec/<source-name>/ directories
- Manifest validation ensuring applicable_scenarios are not empty

## Risks

1. **Managed block corruption**: User edits inside managed blocks will be lost
   - Mitigation: Clear documentation about managed blocks, preserve outside content

2. **File conflicts**: Overwriting existing AGENTS.md or CLAUDE.md without managed blocks
   - Mitigation: Detect missing managed blocks and append instead of replace

3. **Large number of prompts**: Many sources × many prompts = many files
   - Mitigation: Namespace by source, organize in subdirectories

4. **Agent-specific format changes**: Different agents may require different formats
   - Mitigation: Abstract agent generators, test with real agents

## Success Criteria

1. Users can interactively select AI agents (Claude Code, GitHub Copilot)
2. Selected agents are persisted in dnaspec.yaml
3. AGENTS.md is generated with context-aware guideline references
4. CLAUDE.md is generated when Claude selected
5. Claude commands are generated with proper frontmatter and content
6. Copilot prompts are generated with proper frontmatter and $ARGUMENTS
7. User content outside managed blocks is preserved
8. Files can be regenerated without data loss (idempotent)
9. `--no-ask` flag works for CI/CD automation
10. Generated files are immediately usable by real AI agents

## References

- Design document: `docs/design.md` (Agent Integrations, update-agents command)
- Existing change: `add-project-commands` (provides dnaspec.yaml and data structures)
- Existing spec: `manifest-management` (ensures applicable_scenarios exist)
