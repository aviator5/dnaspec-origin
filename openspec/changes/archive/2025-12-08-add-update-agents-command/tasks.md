# Implementation Tasks: Add Update-Agents Command

## 1. Agent Registry and Data Models

- [x] 1.1 Create internal/core/agents/registry.go with supported agent definitions (ID, display name, description)
- [x] 1.2 Define Phase1Agents constant list: ["claude-code", "github-copilot"]
- [x] 1.3 Implement GetAvailableAgents() function returning agent metadata
- [x] 1.4 Implement IsValidAgent(id string) validation function
- [x] 1.5 Write unit tests for agent registry and validation

## 2. Managed Block Utilities

- [x] 2.1 Create internal/core/files/managed.go with managed block constants (DNASPEC:START/END markers)
- [x] 2.2 Implement DetectManagedBlock(content string) function returning (hasBlock bool, startIdx, endIdx int)
- [x] 2.3 Implement ReplaceManagedBlock(content, newBlock string) function preserving outside content
- [x] 2.4 Implement CreateFileWithManagedBlock(newBlock string) function with header
- [x] 2.5 Implement AppendManagedBlock(content, newBlock string) function
- [x] 2.6 Write unit tests for all managed block operations

## 3. AGENTS.md Generation

- [x] 3.1 Create internal/core/agents/agents_md.go with GenerateAgentsMD(config ProjectConfig) function
- [x] 3.2 Implement managed block content generation using applicable_scenarios
- [x] 3.3 Format guideline paths as @/dnaspec/<source-name>/<file>
- [x] 3.4 Format applicable scenarios as bullet list under each guideline
- [x] 3.5 Add instructional text about managed blocks and DNA guidelines
- [x] 3.6 Implement file creation logic: create new / append to existing / replace block
- [x] 3.7 Write unit tests for content generation and file creation scenarios

## 4. CLAUDE.md Generation

- [x] 4.1 Create internal/core/agents/claude_md.go with GenerateClaudeMD(config ProjectConfig) function
- [x] 4.2 Reuse AGENTS.md content generation (same format)
- [x] 4.3 Implement file creation logic with managed blocks
- [x] 4.4 Write unit tests for Claude.md generation

## 5. Claude Command Generation

- [x] 5.1 Create internal/core/agents/claude_commands.go with GenerateClaudeCommand function
- [x] 5.2 Implement frontmatter generation (name, description, category: "DNASpec", tags)
- [x] 5.3 Format command name: "DNASpec: <Source Name> <Prompt Name>"
- [x] 5.4 Wrap prompt content in managed block markers
- [x] 5.5 Implement file path generation: .claude/commands/dnaspec/<source-name>-<prompt-name>.md
- [x] 5.6 Create directory structure if missing
- [x] 5.7 Write unit tests for command generation and frontmatter formatting

## 6. GitHub Copilot Prompt Generation

- [x] 6.1 Create internal/core/agents/copilot_prompts.go with GenerateCopilotPrompt function
- [x] 6.2 Implement frontmatter generation (description only)
- [x] 6.3 Add $ARGUMENTS placeholder before managed block
- [x] 6.4 Wrap prompt content in managed block markers
- [x] 6.5 Implement file path generation: .github/prompts/dnaspec-<source-name>-<prompt-name>.prompt.md
- [x] 6.6 Create directory structure if missing
- [x] 6.7 Write unit tests for prompt generation and frontmatter formatting

## 7. Agent Generator Orchestration

- [x] 7.1 Create internal/core/agents/generate.go with GenerateAgentFiles(config ProjectConfig, agents []string) function
- [x] 7.2 Always generate AGENTS.md regardless of selected agents
- [x] 7.3 Generate CLAUDE.md if "claude-code" in agents list
- [x] 7.4 For each source and prompt, generate Claude commands if "claude-code" selected
- [x] 7.5 For each source and prompt, generate Copilot prompts if "github-copilot" selected
- [x] 7.6 Return summary of generated files (counts by type)
- [x] 7.7 Implement error collection (continue on error, report all at end)
- [x] 7.8 Write unit tests for orchestration logic and file counts

## 8. Agent Selection UI

- [x] 8.1 Create internal/ui/agent_selection.go with SelectAgents function
- [x] 8.2 Display available agents with descriptions using huh or bubbletea
- [x] 8.3 Show currently selected agents as checked (from config)
- [x] 8.4 Implement multi-select with keyboard navigation
- [x] 8.5 Handle user cancellation (return error)
- [x] 8.6 Return selected agent IDs
- [x] 8.7 Write tests for selection logic (if testable)

## 9. Configuration Update

- [x] 9.1 Add SaveAgents(config *ProjectConfig, agents []string) function to internal/core/config/update.go
- [x] 9.2 Update config.Agents field with selected agents
- [x] 9.3 Use atomic write when saving config
- [x] 9.4 Write unit tests for agent configuration updates

## 10. Update-Agents Command

- [x] 10.1 Create internal/cli/project/update_agents.go with NewUpdateAgentsCmd() cobra command
- [x] 10.2 Add --no-ask flag for non-interactive mode
- [x] 10.3 Check dnaspec.yaml exists (error if not: "Run dnaspec init first")
- [x] 10.4 Load project configuration
- [x] 10.5 If no sources configured, display helpful message and exit
- [x] 10.6 Interactive mode: call SelectAgents, save to config
- [x] 10.7 Non-interactive mode: use config.Agents, error if empty
- [x] 10.8 Call GenerateAgentFiles with config and agents
- [x] 10.9 Display generation summary with file counts
- [x] 10.10 Add success message with lipgloss formatting
- [x] 10.11 Register command in cmd/dnaspec/main.go

## 11. Error Handling and UX

- [x] 11.1 Handle missing dnaspec.yaml gracefully
- [x] 11.2 Handle empty sources list gracefully
- [x] 11.3 Handle missing applicable_scenarios (should be prevented by manifest validation)
- [x] 11.4 Add descriptive errors for file write failures (permissions, disk space)
- [x] 11.5 Display warnings for prompts without content files
- [x] 11.6 Use lipgloss for styled output (success, errors, summaries)
- [x] 11.7 Add progress indicators for bulk file generation

## 12. Testing and Validation

- [x] 12.1 Achieve >80% code coverage for core agent generation logic
- [x] 12.2 Create test fixtures (sample config with multiple sources, guidelines, prompts)
- [x] 12.3 Test managed block detection and replacement edge cases
- [x] 12.4 Test file creation scenarios (new file, existing without block, existing with block)
- [x] 12.5 Test prompt generation for all supported agents
- [x] 12.6 Test source name namespacing (multiple sources with same prompt names)
- [x] 12.7 Write integration tests for full update-agents workflow
- [x] 12.8 Manual testing: run update-agents after dnaspec add
- [x] 12.9 Manual testing: verify generated files work with real Claude Code
- [x] 12.10 Manual testing: verify generated files work with real GitHub Copilot
- [x] 12.11 Manual testing: test --no-ask flag behavior
- [x] 12.12 Manual testing: test idempotency (run multiple times, check preservation)

## 13. Integration with Existing Commands

- [x] 13.1 Update dnaspec add command to suggest "Run 'dnaspec update-agents' to configure AI agents"
- [x] 13.2 Update dnaspec init command to mention update-agents in next steps
- [x] 13.3 Ensure update-agents can be run multiple times safely (idempotent)

## 14. Documentation

- [x] 14.1 Update README with dnaspec update-agents command example
- [x] 14.2 Document --no-ask flag usage
- [x] 14.3 Document managed block concept and preservation rules
- [x] 14.4 Add examples of generated AGENTS.md
- [x] 14.5 Add examples of generated Claude commands and Copilot prompts
- [x] 14.6 Document how to add custom content to AGENTS.md outside managed blocks
- [x] 14.7 Document agent file locations and formats

## Notes

- Always generate AGENTS.md even if no agents selected (base instructions for all AI)
- Use applicable_scenarios from guidelines to create context-aware instructions
- Preserve user content outside managed blocks (critical UX requirement)
- Namespace prompt files by source name to prevent collisions
- Atomic writes for all file operations
- Clear error messages with actionable suggestions
- Design agent generators to be extensible (future: Windsurf, Cursor)
- Test with real AI agents (Claude Code, GitHub Copilot) to validate formats
