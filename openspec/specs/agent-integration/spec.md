# agent-integration Specification

## Purpose

Enable AI coding assistants (Claude Code, GitHub Copilot) to discover and use DNA guidelines by generating agent-specific integration files with context-aware instructions.

## Related Specifications

- **manifest-management**: Ensures applicable_scenarios field is populated and validated
- **project-management**: Provides dnaspec.yaml structure and ProjectConfig data model
## Requirements
### Requirement: Agent Selection and Configuration

The system SHALL provide a `dnaspec update-agents` command that allows users to select AI agents and persist the selection in configuration.

#### Scenario: Interactive agent selection

- **WHEN** user runs `dnaspec update-agents` without flags
- **THEN** display checklist of available agents (antigravity, claude-code, cursor, github-copilot, windsurf) in alphabetical order
- **AND** show currently selected agents as checked
- **AND** allow multi-select with keyboard navigation
- **AND** save selected agents to dnaspec.yaml agents array
- **AND** proceed to file generation

### Requirement: AGENTS.md Generation with Managed Blocks

The system SHALL generate AGENTS.md with guideline references and applicable scenarios, using managed blocks to preserve user content.

#### Scenario: Create new AGENTS.md file

- **WHEN** AGENTS.md does not exist and update-agents runs
- **THEN** create AGENTS.md with header and managed block containing guideline instructions
- **AND** include all guidelines from all sources with applicable scenarios

#### Scenario: Update existing AGENTS.md with managed block

- **WHEN** AGENTS.md exists with `<!-- DNASPEC:START -->` and `<!-- DNASPEC:END -->` markers
- **THEN** replace content between markers with updated guideline instructions
- **AND** preserve all content outside the markers exactly as-is

#### Scenario: Append to existing AGENTS.md without managed block

- **WHEN** AGENTS.md exists but does not contain managed block markers
- **THEN** append managed block at end of file
- **AND** preserve all existing content

#### Scenario: Guideline references with applicable scenarios

- **WHEN** AGENTS.md is generated
- **THEN** each guideline reference includes path format `@/dnaspec/<source-name>/<file>`
- **AND** each guideline lists applicable scenarios as bullet points
- **AND** scenarios are taken from guideline's applicable_scenarios field

#### Scenario: Multiple sources with guidelines

- **WHEN** project has multiple sources each with guidelines
- **THEN** AGENTS.md includes all guidelines from all sources
- **AND** guidelines are organized by source
- **AND** source names appear in paths to prevent ambiguity

### Requirement: CLAUDE.md Generation for Claude Code

The system SHALL generate CLAUDE.md with same content as AGENTS.md when Claude Code is selected.

#### Scenario: Generate CLAUDE.md when Claude selected

- **WHEN** user selects "claude-code" agent
- **THEN** generate CLAUDE.md with same content as AGENTS.md
- **AND** use same managed block preservation logic

#### Scenario: Skip CLAUDE.md when Claude not selected

- **WHEN** user does not select "claude-code" agent
- **THEN** do not generate or update CLAUDE.md

### Requirement: Claude Command File Generation

The system SHALL generate Claude slash command files for each prompt when Claude Code is selected.

#### Scenario: Generate Claude command for prompt

- **WHEN** "claude-code" is selected and source has prompt
- **THEN** create `.claude/commands/dnaspec/<source-name>-<prompt-name>.md`
- **AND** include frontmatter with name, description, category "DNASpec", and tags
- **AND** format command name as "DNASpec: <Source Name> <Prompt Name>"
- **AND** wrap prompt content in managed block markers
- **AND** create directory structure if missing

#### Scenario: Source namespacing prevents collisions

- **WHEN** multiple sources have prompts with same name
- **THEN** each generates separate file with source name in filename
- **AND** files do not overwrite each other

#### Scenario: Skip Claude commands when not selected

- **WHEN** "claude-code" is not selected
- **THEN** do not generate Claude command files

### Requirement: GitHub Copilot Prompt File Generation

The system SHALL generate Copilot prompt files for each prompt when GitHub Copilot is selected.

#### Scenario: Generate Copilot prompt for prompt

- **WHEN** "github-copilot" is selected and source has prompt
- **THEN** create `.github/prompts/dnaspec-<source-name>-<prompt-name>.prompt.md`
- **AND** include frontmatter with description
- **AND** include `$ARGUMENTS` placeholder before managed block
- **AND** wrap prompt content in managed block markers
- **AND** create directory structure if missing

#### Scenario: Source namespacing in Copilot prompts

- **WHEN** multiple sources have prompts with same name
- **THEN** each generates separate file with source name in filename
- **AND** files do not overwrite each other

#### Scenario: Skip Copilot prompts when not selected

- **WHEN** "github-copilot" is not selected
- **THEN** do not generate Copilot prompt files

### Requirement: Idempotent File Generation

The system SHALL generate files idempotently so running update-agents multiple times produces consistent results.

#### Scenario: Regenerate files without changes

- **WHEN** user runs `dnaspec update-agents --no-ask` multiple times without config changes
- **THEN** generated files have identical content each time
- **AND** user content outside managed blocks is preserved
- **AND** operation completes successfully

#### Scenario: Update after adding guidelines

- **WHEN** user adds new guidelines with `dnaspec add` then runs `dnaspec update-agents --no-ask`
- **THEN** AGENTS.md and CLAUDE.md updated to include new guidelines
- **AND** new prompt files generated for new prompts
- **AND** existing files remain unchanged
- **AND** user content outside managed blocks preserved

### Requirement: Generation Summary Display

The system SHALL display a summary of generated files after successful completion.

#### Scenario: Display generation summary

- **WHEN** update-agents completes successfully
- **THEN** display summary with:
  - "✓ Updated AGENTS.md"
  - "✓ Updated CLAUDE.md" (if Claude selected)
  - "✓ Generated X Claude commands" (if Claude selected)
  - "✓ Generated X Copilot prompts" (if Copilot selected)
  - "✓ Generated X Antigravity prompts" (if Antigravity selected)
  - "✓ Generated X Windsurf workflows" (if Windsurf selected)
  - "✓ Generated X Cursor commands" (if Cursor selected)
- **AND** display success message

### Requirement: Managed Block Marker Format

The system SHALL use consistent managed block markers for all generated content.

#### Scenario: Standard managed block markers

- **WHEN** any file is generated with managed block
- **THEN** use exactly `<!-- DNASPEC:START -->` as start marker
- **AND** use exactly `<!-- DNASPEC:END -->` as end marker
- **AND** place markers on their own lines

### Requirement: Error Handling for Missing Prerequisites

The system SHALL provide clear error messages when prerequisites are missing.

#### Scenario: Config file not found

- **WHEN** user runs `dnaspec update-agents` without dnaspec.yaml in current or parent directories
- **THEN** exit with error "dnaspec.yaml not found. Run 'dnaspec init' first."

#### Scenario: Missing prompt file

- **WHEN** prompt references file that does not exist
- **THEN** display warning "Prompt file not found: {path}"
- **AND** continue processing other prompts
- **AND** include warning in summary

#### Scenario: File permission errors

- **WHEN** file generation fails due to permission error
- **THEN** display error with file path and reason
- **AND** continue processing other files
- **AND** exit with non-zero code

### Requirement: Antigravity Agent Support

The system SHALL generate prompt files for Antigravity when selected.

#### Scenario: Generate Antigravity prompt for prompt

- **WHEN** "antigravity" is selected and source has prompt
- **THEN** create `.agent/workflows/dnaspec-<source-name>-<prompt-name>.md`
- **AND** include frontmatter with description field
- **AND** wrap prompt content in managed block markers (`<!-- DNASPEC:START -->` and `<!-- DNASPEC:END -->`)
- **AND** create directory structure if missing

#### Scenario: Source namespacing in Antigravity prompts

- **WHEN** multiple sources have prompts with same name
- **THEN** each generates separate file with source name in filename
- **AND** files do not overwrite each other

#### Scenario: Skip Antigravity prompts when not selected

- **WHEN** "antigravity" is not selected
- **THEN** do not generate Antigravity prompt files

### Requirement: Windsurf Agent Support

The system SHALL generate prompt files for Windsurf when selected with auto-execution configuration.

#### Scenario: Generate Windsurf prompt for prompt

- **WHEN** "windsurf" is selected and source has prompt
- **THEN** create `.windsurf/workflows/dnaspec-<source-name>-<prompt-name>.md`
- **AND** include frontmatter with description field
- **AND** include frontmatter with `auto_execution_mode: 3` field
- **AND** wrap prompt content in managed block markers (`<!-- DNASPEC:START -->` and `<!-- DNASPEC:END -->`)
- **AND** create directory structure if missing

#### Scenario: Source namespacing in Windsurf prompts

- **WHEN** multiple sources have prompts with same name
- **THEN** each generates separate file with source name in filename
- **AND** files do not overwrite each other

#### Scenario: Skip Windsurf prompts when not selected

- **WHEN** "windsurf" is not selected
- **THEN** do not generate Windsurf prompt files

### Requirement: Cursor Agent Support

The system SHALL generate command files for Cursor when selected with complete metadata.

#### Scenario: Generate Cursor command for prompt

- **WHEN** "cursor" is selected and source has prompt
- **THEN** create `.cursor/commands/dnaspec-<source-name>-<prompt-name>.md`
- **AND** include frontmatter with name field formatted as `/dnaspec-<source-name>-<prompt-name>`
- **AND** include frontmatter with id field matching prompt filename without extension
- **AND** include frontmatter with category field set to "DNASpec"
- **AND** include frontmatter with description field
- **AND** wrap prompt content in managed block markers (`<!-- DNASPEC:START -->` and `<!-- DNASPEC:END -->`)
- **AND** create directory structure if missing

#### Scenario: Source namespacing in Cursor commands

- **WHEN** multiple sources have prompts with same name
- **THEN** each generates separate file with source name in filename
- **AND** files do not overwrite each other

#### Scenario: Skip Cursor commands when not selected

- **WHEN** "cursor" is not selected
- **THEN** do not generate Cursor command files

### Requirement: Alphabetical Agent Display

The system SHALL display available agents in alphabetical order by display name.

#### Scenario: Agent list sorted alphabetically

- **WHEN** user runs `dnaspec update-agents` without flags
- **THEN** display agent selection list sorted alphabetically by display name
- **AND** maintain consistent ordering across all invocations

