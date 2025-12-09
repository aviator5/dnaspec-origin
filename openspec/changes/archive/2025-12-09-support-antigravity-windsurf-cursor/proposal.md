# Proposal: Support Antigravity, Windsurf, and Cursor Agents

## Problem Statement

DNASpec currently supports only Claude Code and GitHub Copilot (Phase 1 agents). Users want to integrate DNA guidelines with additional AI coding assistants: Antigravity, Windsurf, and Cursor. Each has its own file format and directory conventions for guidelines and prompts.

## Proposed Solution

Extend the agent registry and generation system to support three additional AI agents:

1. **Antigravity** - AI development assistant
   - Guidelines: AGENTS.md (shared)
   - Prompts: `.agent/workflows/dnaspec-<source-name>-<prompt-name>.md`

2. **Windsurf** - AI-powered code editor
   - Guidelines: AGENTS.md (shared)
   - Prompts: `.windsurf/workflows/dnaspec-<source-name>-<prompt-name>.md`
   - Additional: `auto_execution_mode: 3` in frontmatter

3. **Cursor** - AI-first code editor
   - Guidelines: AGENTS.md (shared)
   - Prompts: `.cursor/commands/dnaspec-<source-name>-<prompt-name>.md`
   - Additional: `name`, `id`, `category` fields in frontmatter

All agents will share the same AGENTS.md file for guideline discovery.

## User Experience Improvements

- **Alphabetical sorting**: Display agents in alphabetical order in selection UI for easier navigation
- **Comprehensive documentation**: Update all docs to reflect new agent support

## Goals

1. Extend agent registry with three new agents
2. Implement prompt file generators for each agent
3. Sort agent list alphabetically in UI
4. Update all documentation files

## Non-Goals

- Changing existing Claude Code or GitHub Copilot implementations
- Adding agents beyond these three (future work)
- Modifying AGENTS.md format (remains shared across all agents)

## Success Criteria

- Users can select antigravity, windsurf, or cursor in `dnaspec update-agents`
- Prompt files generate in correct directories with proper formats
- Agent list displays alphabetically
- Documentation accurately reflects new capabilities
- All tests pass
