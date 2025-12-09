# Spec Delta: Agent Integration

## ADDED Requirements

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

## MODIFIED Requirements

### Requirement: Agent Selection and Configuration

The system SHALL provide a `dnaspec update-agents` command that allows users to select AI agents and persist the selection in configuration.

#### Scenario: Interactive agent selection

- **WHEN** user runs `dnaspec update-agents` without flags
- **THEN** display checklist of available agents (antigravity, claude-code, cursor, github-copilot, windsurf) in alphabetical order
- **AND** show currently selected agents as checked
- **AND** allow multi-select with keyboard navigation
- **AND** save selected agents to dnaspec.yaml agents array
- **AND** proceed to file generation

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

## Related Specifications

- **project-management**: Provides agent configuration storage in dnaspec.yaml
- **manifest-management**: Provides prompt definitions and metadata
