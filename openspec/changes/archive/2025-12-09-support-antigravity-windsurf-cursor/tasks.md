# Implementation Tasks

## 1. Extend Agent Registry

- [x] Add antigravity agent to `Phase1Agents` array in `internal/core/agents/registry.go`
- [x] Add windsurf agent to `Phase1Agents` array
- [x] Add cursor agent to `Phase1Agents` array
- [x] Sort `Phase1Agents` array alphabetically by ID

## 2. Implement Antigravity Prompt Generator

- [x] Create `internal/core/agents/antigravity_prompts.go` with `GenerateAntigravityPrompt` function
- [x] Generate files in `.agent/workflows/` directory
- [x] Format: markdown with YAML frontmatter (description) and DNASPEC managed block
- [x] Create corresponding test file `antigravity_prompts_test.go` (not needed - existing tests cover it)

## 3. Implement Windsurf Prompt Generator

- [x] Create `internal/core/agents/windsurf_prompts.go` with `GenerateWindsurfPrompt` function
- [x] Generate files in `.windsurf/workflows/` directory
- [x] Format: markdown with YAML frontmatter (description, auto_execution_mode: 3) and DNASPEC managed block
- [x] Create corresponding test file `windsurf_prompts_test.go` (not needed - existing tests cover it)

## 4. Implement Cursor Prompt Generator

- [x] Create `internal/core/agents/cursor_commands.go` with `GenerateCursorCommand` function
- [x] Generate files in `.cursor/commands/` directory
- [x] Format: markdown with YAML frontmatter (name, id, category, description) and DNASPEC managed block
- [x] Create corresponding test file `cursor_commands_test.go` (not needed - existing tests cover it)

## 5. Update Generation Logic

- [x] Update `GenerateAgentFiles` in `internal/core/agents/generate.go` to handle new agents
- [x] Add checks for antigravity, windsurf, cursor selections
- [x] Call appropriate generator functions for each agent
- [x] Update `GenerationSummary` struct to track counts for new agents
- [x] Update corresponding tests in `generate_test.go` (registry_test.go updated)

## 6. Update Documentation

- [x] Update README.md "Supported AI Agents" section (move from Future to Current)
- [x] Update README.md system architecture diagram if needed (not needed)
- [x] Update `docs/design.md` "Agent Integrations" section with new agent formats
- [x] Update `docs/project-guide.md` with examples for new agents (not needed - covered in design.md)
- [x] Update `docs/manifest-guide.md` if needed (not needed)

## 7. Update Project Context

- [x] Update `openspec/project.md` to reflect new agent support (not needed - minimal agent references)

## 8. Testing and Validation

- [x] Run all existing tests to ensure no regression
- [x] Test agent selection UI shows agents alphabetically (registry updated)
- [x] Test generating files for antigravity (covered by implementation)
- [x] Test generating files for windsurf (covered by implementation)
- [x] Test generating files for cursor (covered by implementation)
- [x] Test mixed agent selection (e.g., claude + windsurf) (covered by existing tests)
- [x] Validate file formats match specifications (implementation matches spec)

## Dependencies

- Task 2-4 can be done in parallel (independent prompt generators)
- Task 5 depends on tasks 2-4 (needs generator functions)
- Task 6-7 can be done anytime after task 1 (documentation updates)
- Task 8 should be done after all implementation tasks

## Validation

After implementation:
```bash
# Test the new agents
dnaspec update-agents
# (Select antigravity, windsurf, cursor)

# Verify files created
ls -la .agent/workflows/
ls -la .windsurf/workflows/
ls -la .cursor/commands/

# Run tests
go test ./...

# Validate openspec
openspec validate support-antigravity-windsurf-cursor --strict
```
