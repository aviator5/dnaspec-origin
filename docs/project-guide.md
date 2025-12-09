# DNASpec Project Developer Guide

Complete guide for integrating DNA guidelines into your projects.

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Project Commands Reference](#project-commands-reference)
  - [dnaspec init](#dnaspec-init)
  - [dnaspec add](#dnaspec-add)
  - [dnaspec remove](#dnaspec-remove)
  - [dnaspec list](#dnaspec-list)
  - [dnaspec update](#dnaspec-update)
  - [dnaspec update-agents](#dnaspec-update-agents)
  - [dnaspec validate](#dnaspec-validate)
  - [dnaspec sync](#dnaspec-sync)
- [Project Configuration](#project-configuration)
- [Using AI Agents with DNA](#using-ai-agents-with-dna)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)
- [Development](#development)

## Installation

### Download Pre-built Binary (Recommended)

Download the latest release for your platform from [GitHub Releases](https://github.com/aviator5/dnaspec/releases):

**macOS:**
```bash
# ARM (M1/M2/M3)
curl -L https://github.com/aviator5/dnaspec/releases/latest/download/dnaspec_darwin_arm64.tar.gz | tar xz
sudo mv dnaspec /usr/local/bin/

# Intel
curl -L https://github.com/aviator5/dnaspec/releases/latest/download/dnaspec_darwin_amd64.tar.gz | tar xz
sudo mv dnaspec /usr/local/bin/
```

**Linux:**
```bash
# x86_64
curl -L https://github.com/aviator5/dnaspec/releases/latest/download/dnaspec_linux_amd64.tar.gz | tar xz
sudo mv dnaspec /usr/local/bin/

# ARM64
curl -L https://github.com/aviator5/dnaspec/releases/latest/download/dnaspec_linux_arm64.tar.gz | tar xz
sudo mv dnaspec /usr/local/bin/
```

**Windows:**
```powershell
# Download from: https://github.com/aviator5/dnaspec/releases/latest
# Extract and add to PATH
```

### Install via Go

If you have Go 1.21+ installed:

```bash
go install github.com/aviator5/dnaspec/cmd/dnaspec@latest
```

### Verify Installation

```bash
dnaspec version
```

## Quick Start

### Single Source Setup

1. Initialize DNASpec in your project:
```bash
dnaspec init
```

2. Add DNA guidelines from a repository:
```bash
dnaspec add --git-repo https://github.com/company/dna-guidelines
```

3. Configure AI agents to use your guidelines:
```bash
dnaspec update-agents
```

4. Or add from a local directory:
```bash
dnaspec add /path/to/local/dna-guidelines
```

### Multiple Sources Setup

DNASpec supports adding guidelines from **multiple sources** to a single project:

```bash
# Add company-wide guidelines
dnaspec add --git-repo https://github.com/company/dna --name company

# Add team-specific patterns
dnaspec add --git-repo https://github.com/team/patterns --name team

# Add personal experimental guidelines
dnaspec add /Users/me/my-dna --name personal

# Generate agent configuration with all sources
dnaspec update-agents
```

**Benefits of multiple sources:**
- **Company-wide standards**: Shared across all organization projects
- **Team-specific patterns**: Tailored to your team's tech stack
- **Personal guidelines**: Your own best practices
- **Independent versioning**: Update each source on its own schedule
- **Namespace isolation**: No conflicts between sources (e.g., `company-review` vs `team-review`)

All guidelines from all sources are available to AI agents simultaneously.

## Project Commands Reference

### `dnaspec init`

Initialize a new `dnaspec.yaml` file in your project.

```bash
dnaspec init
```

This command:
- Creates a new `dnaspec.yaml` configuration file in the current directory
- Includes commented examples showing how to add DNA sources
- Prevents overwriting an existing configuration file

**Example output:**
```
✓ Success: Created dnaspec.yaml

Next steps:
  1. Run dnaspec add to add DNA guidelines from a repository or local directory
  2. Run dnaspec update-agents to generate agent configuration files
```

### `dnaspec add`

Add DNA guidelines from a git repository or local directory to your project.

**Add from git repository:**
```bash
# Add from default branch/tag
dnaspec add --git-repo https://github.com/company/dna-guidelines

# Add from specific branch or tag
dnaspec add --git-repo https://github.com/company/dna-guidelines --git-ref v1.2.0

# Add with custom source name
dnaspec add --git-repo https://github.com/company/dna-guidelines --name my-dna
```

**Add from local directory:**
```bash
# Add from local path
dnaspec add /path/to/local/dna-guidelines

# Add with custom source name
dnaspec add /path/to/local/dna-guidelines --name my-local-dna
```

**Selection options:**
```bash
# Interactive selection (default)
dnaspec add --git-repo https://github.com/company/dna-guidelines

# Add all guidelines without prompting
dnaspec add --git-repo https://github.com/company/dna-guidelines --all

# Add specific guidelines
dnaspec add --git-repo https://github.com/company/dna-guidelines --guideline go-style --guideline rest-api

# Preview changes without modifying files
dnaspec add --git-repo https://github.com/company/dna-guidelines --dry-run
```

This command:
- Clones the git repository (for git sources) or reads the local directory
- Parses the `dnaspec-manifest.yaml` file from the source
- Shows an interactive guideline selection (unless `--all` or `--guideline` flags are used)
- Copies selected guideline and prompt files to `dnaspec/<source-name>/` directory
- Updates `dnaspec.yaml` with source metadata and selected guidelines

**Flags:**
- `--git-repo <url>`: Git repository URL (https:// or git@)
- `--git-ref <ref>`: Git branch or tag to use (defaults to repository's default branch)
- `--name <name>`: Custom source name (defaults to derived name from URL/path)
- `--all`: Add all guidelines without interactive selection
- `--guideline <name>`: Add specific guideline by name (can be repeated)
- `--dry-run`: Preview changes without modifying files

**Example output:**
```
→ Cloning repository https://github.com/company/dna-guidelines...
✓ Repository cloned successfully

Select guidelines to add:
  [x] go-style - Go coding style conventions
  [x] rest-api - REST API design principles
  [ ] database-design - Database schema design guidelines

→ Copying 2 guidelines and 3 prompts to dnaspec/dna-guidelines/...
✓ Files copied successfully
✓ Updated dnaspec.yaml

Added source: dna-guidelines
  Guidelines: go-style, rest-api
  Prompts: code-review, documentation, api-review

Next steps:
  1. Review the added guidelines in dnaspec/dna-guidelines/
  2. Run dnaspec update-agents to generate agent configuration files
```

### `dnaspec remove`

Remove a DNA source from your project configuration. This command safely removes the source from `dnaspec.yaml`, deletes the source directory and all guideline files, and cleans up generated agent files (Claude commands and Copilot prompts).

**Basic usage:**
```bash
# Remove a source with confirmation prompt
dnaspec remove company-dna

# Remove a source without confirmation
dnaspec remove company-dna --force
```

This command:
- Shows what will be deleted (config entry, source directory, generated agent files)
- Prompts for confirmation before deletion (unless `--force` is used)
- Removes the source entry from `dnaspec.yaml`
- Deletes the `dnaspec/<source-name>/` directory
- Cleans up generated Claude command files (`.claude/commands/dnaspec/<source-name>-*.md`)
- Cleans up generated Copilot prompt files (`.github/prompts/dnaspec-<source-name>-*.prompt.md`)
- Handles missing files gracefully (idempotent operation)

**Flags:**
- `--force`, `-f`: Skip confirmation prompt

**Example output:**
```
The following will be deleted:
  - dnaspec.yaml entry for 'company-dna'
  - dnaspec/company-dna/ directory (5 guidelines, 3 prompts)
  - .claude/commands/dnaspec/company-dna-*.md (3 files)
  - .github/prompts/dnaspec-company-dna-*.prompt.md (2 files)

This cannot be undone. Continue? [y/N]: y

✓ Success: Removed source company-dna
  Cleaned up 5 file(s)

Next steps:
  Run dnaspec update-agents to regenerate AGENTS.md
```

**Error handling:**
- If the source doesn't exist, shows available sources
- If config file is missing, suggests running `dnaspec init`
- If file deletion fails, provides clear error messages

### `dnaspec list`

View the current DNA configuration for your project, showing all configured agents, sources, guidelines, and prompts.

```bash
dnaspec list
```

This command:
- Loads and displays the `dnaspec.yaml` configuration file
- Shows configured AI agents (Phase 1: Claude Code, GitHub Copilot)
- Lists all DNA sources with their type-specific metadata
- Displays guidelines and prompts for each source
- Provides a quick overview of your project's DNA setup

**Example output:**
```
Configured Agents (Phase 1):
  - claude-code
  - github-copilot

Sources:

company-dna (git-repo)
  URL: https://github.com/company/dna-guidelines
  Ref: v1.2.0
  Commit: abc123def456789...
  Guidelines:
    - go-style: Go coding style conventions
    - rest-api: REST API design principles
  Prompts:
    - code-review: Review Go code for style compliance
    - api-review: Review API designs

local-patterns (local-path)
  Path: /Users/me/my-dna-patterns
  Guidelines:
    - error-handling: Error handling patterns
```

**When no configuration exists:**
```
✗ Error: dnaspec.yaml not found

Run 'dnaspec init' to create a new configuration file.
```

**Use cases:**
- Verify which DNA sources are currently active
- Check which guidelines are available before running code reviews
- Confirm agent configuration before running `dnaspec update-agents`
- Debug configuration issues

### `dnaspec update`

Update DNA sources from their origins (git repositories or local directories) to fetch the latest guidelines and prompts.

**Update specific source:**
```bash
# Update a single source
dnaspec update my-company-dna

# Preview changes without writing files
dnaspec update my-company-dna --dry-run

# Add all new guidelines automatically
dnaspec update my-company-dna --add-new=all

# Skip new guidelines
dnaspec update my-company-dna --add-new=none
```

**Update all sources:**
```bash
# Update all sources at once
dnaspec update --all

# Update all with automatic new guideline handling
dnaspec update --all --add-new=all
```

This command:
- Fetches the latest manifest from the source (git clone or local directory read)
- Compares current configuration with latest manifest
- Updates metadata for existing guidelines (description, scenarios, prompts)
- Copies updated guideline and prompt files to `dnaspec/<source-name>/` directory
- Optionally adds new guidelines (interactive by default)
- Reports guidelines removed from source (but keeps local files)
- Updates `dnaspec.yaml` with new commit hashes (git sources) and metadata

**Flags:**
- `--all`: Update all configured sources
- `--dry-run`: Preview changes without modifying files
- `--add-new <policy>`: Policy for new guidelines (`all` or `none`). If not specified, prompts interactively.

**Example output:**
```
⏳ Fetching latest from https://github.com/company/dna...
✓ Current commit: abc123de
✓ Latest commit: def456ab (changed)

Updated guidelines:
  ✓ go-style (description changed)
  ✓ rest-api (content updated)

New guidelines available:
  - go-testing: Go testing patterns
  - go-errors: Error handling conventions

Removed from source:
  - old-guideline (no longer in manifest)

Add new guidelines? [y/N]: y

✓ Added go-testing
✓ Added go-errors

✓ Updated dnaspec.yaml

Run 'dnaspec update-agents' to regenerate agent files
```

**When sources are up to date:**
```
⏳ Fetching latest from https://github.com/company/dna...
✓ Current commit: abc123de
✓ Already at latest commit

All guidelines up to date.
```

### `dnaspec update-agents`

Generate or update AI agent configuration files based on selected DNA guidelines.

```bash
# Interactive mode - select which agents to configure
dnaspec update-agents

# Non-interactive mode - use saved selection
dnaspec update-agents --no-ask
```

This command:
- Shows an interactive checklist of available AI agents (Claude Code, GitHub Copilot)
- Saves your agent selection to `dnaspec.yaml`
- Generates agent-specific integration files based on your DNA guidelines
- Updates managed blocks while preserving custom content outside those blocks

**Generated Files:**

**AGENTS.md** (always):
- Contains references to all guidelines with their applicable scenarios
- Format: `@/dnaspec/<source-name>/<file>` with scenario bullet points
- Uses managed blocks (`<!-- DNASPEC:START/END -->`) that can be safely regenerated

**CLAUDE.md** (if Claude Code selected):
- Same content as AGENTS.md
- Enables Claude Code to discover project-specific instructions

**Claude Commands** (if Claude Code selected):
- `.claude/commands/dnaspec/<source-name>-<prompt-name>.md`
- Slash commands for guideline-based code reviews and tasks
- Includes frontmatter with name, description, category, and tags

**Copilot Prompts** (if GitHub Copilot selected):
- `.github/prompts/dnaspec-<source-name>-<prompt-name>.prompt.md`
- Prompts for guideline-based assistance
- Includes `$ARGUMENTS` placeholder for context

**Flags:**
- `--no-ask`: Use saved agent configuration without prompting (useful for CI/CD)

**Example output:**
```
Select AI agents to configure:
  [x] Claude Code
  [x] GitHub Copilot

→ Generating agent configuration files...
✓ Created AGENTS.md
✓ Created CLAUDE.md
✓ Created .claude/commands/dnaspec/dna-guidelines-code-review.md
✓ Created .claude/commands/dnaspec/dna-guidelines-documentation.md
✓ Created .github/prompts/dnaspec-dna-guidelines-code-review.prompt.md
✓ Created .github/prompts/dnaspec-dna-guidelines-documentation.prompt.md

Agent configuration updated successfully!

Your AI assistants can now access your DNA guidelines:
  • Claude Code: Type / and search for "dnaspec" commands
  • GitHub Copilot: Use GitHub Copilot Chat with dnaspec prompts
```

**When to run:**
- After adding new DNA sources with `dnaspec add`
- After changing which agents you want to use
- After updating guidelines in your DNA sources (with `--no-ask` flag)

**Managed Blocks:**

The command uses managed block markers to safely update generated content:

```markdown
<!-- DNASPEC:START -->
Generated content here
<!-- DNASPEC:END -->
```

Content outside these markers is preserved, so you can add custom instructions alongside generated guidelines.

### `dnaspec validate`

Validate the project configuration (`dnaspec.yaml`) without modifying any files.

```bash
dnaspec validate
```

This command checks:
- **YAML syntax and schema structure**: Ensures the configuration file is valid
- **Config version**: Verifies version is supported (currently version 1)
- **Source fields**: Checks all sources have required fields based on type
- **File references**: Verifies all guideline and prompt files exist in `dnaspec/` directory
- **Agent IDs**: Validates agent IDs are recognized (claude-code, github-copilot)
- **Duplicate names**: Checks for duplicate source names
- **Comprehensive error reporting**: Collects and displays all errors before exiting

**Example output (success):**
```
Validating dnaspec.yaml...
✓ YAML syntax valid
✓ Version 1 schema valid
✓ 2 sources configured
✓ All referenced files exist:
  - dnaspec/company-dna/guidelines/go-style.md
  - dnaspec/company-dna/guidelines/rest-api.md
  - dnaspec/company-dna/prompts/code-review.md
✓ All agent IDs recognized: claude-code, github-copilot
✓ Configuration is valid
```

**Example output (errors):**
```
Validating dnaspec.yaml...
✓ YAML syntax valid
✗ Found 3 validation error(s):

  • Source 'company-dna' (git-repo) missing required field: url
  • File not found: dnaspec/company-dna/guidelines/missing.md
  • Unknown agent ID: 'invalid-agent' (recognized: claude-code, github-copilot)

Validation failed with 3 error(s)
```

**Use cases:**
- Verify configuration before running other commands
- Debug configuration issues
- Check in CI/CD pipelines to catch errors early
- Ensure all referenced files exist after cloning a repository

### `dnaspec sync`

Update all DNA sources and regenerate agent files in a single non-interactive operation.

```bash
# Sync all sources and regenerate agent files
dnaspec sync

# Preview changes without writing files
dnaspec sync --dry-run
```

This command is a convenience wrapper that:
1. Updates all sources from their origins (equivalent to `dnaspec update --all --add-new=none`)
2. Regenerates all agent files (equivalent to `dnaspec update-agents --no-ask`)
3. Displays consolidated summary of all changes

**Non-interactive design:**
- Safe for CI/CD pipelines - never prompts for user input
- Uses saved agent configuration from `dnaspec.yaml`
- Does NOT add new guidelines automatically (uses `--add-new=none` policy)
- Exits early if any source update fails

**Flags:**
- `--dry-run`: Preview changes without modifying files

**Example output:**
```
Syncing all DNA sources...

Updating 2 sources...

=== Updating company-dna ===
⏳ Refreshing from https://github.com/company/dna...
✓ Updated 1 guideline
✓ No new guidelines available

=== Updating local-patterns ===
⏳ Refreshing from local directory...
✓ No changes (already up to date)

✓ All sources updated

Regenerating agent files...
Using saved agents: [claude-code, github-copilot]
✓ AGENTS.md
✓ CLAUDE.md
✓ Generated 2 Claude command(s)
✓ Generated 2 Copilot prompt(s)

✓ Sync complete
```

**When to use:**
- In CI/CD pipelines to keep DNA guidelines up to date
- Before starting work to ensure you have latest guidelines
- After pulling repository changes that might affect DNA sources
- When you want to update everything with a single command

**Comparison with individual commands:**

| Command | Interactive | Adds New Guidelines | Use Case |
|---------|------------|---------------------|----------|
| `dnaspec update --all` | Yes (prompts for new guidelines) | Optional | Manual updates with decisions |
| `dnaspec sync` | No | Never | CI/CD, automated workflows |

## Project Configuration

The `dnaspec.yaml` file in your project directory tracks which DNA sources you've added and which guidelines are active.

```yaml
version: 1

agents:
  - "claude-code"
  - "github-copilot"

sources:
  - name: "company-dna"          # Human-readable source name
    type: "git-repo"                # or "local-path"
    url: "https://github.com/company/dna-guidelines"
    ref: "v1.2.0"                   # Git tag or branch
    commit: "abc123def456789..."    # Resolved commit hash
    guidelines:
      - name: "go-style"            # spinal-case
        file: "guidelines/go-style.md"   # Relative to source root
        description: "Go code style conventions"
        applicable_scenarios:       # Used for generating AGENTS.md
          - "writing new Go code"
        prompts:                    # List of prompt names (not paths)
          - "code-review"
    prompts:
      - name: "code-review"   # spinal-case
        file: "prompts/code-review.md"
        description: "Review Go code"

  - name: "local-patterns"
    type: "local-path"
    path: "/Users/me/my-dna-patterns"
    guidelines:
      - name: "error-handling"
        file: "guidelines/error-handling.md"
        description: "Error handling patterns"
        applicable_scenarios:
          - "handling errors"
    prompts: []
```

### Configuration Fields

**Top-level:**
- `version`: Must be `1`
- `agents`: List of AI agents to generate configuration for (values: `"claude-code"`, `"github-copilot"`)
- `sources`: List of DNA sources added to this project

**Source (git-repo type):**
- `name`: Unique source identifier (derived from URL or custom via `--name`)
- `type`: Must be `"git-repo"`
- `url`: Git repository URL
- `ref`: Git branch or tag used (optional)
- `commit`: Git commit hash for tracking updates
- `guidelines`: List of selected guidelines from this source
- `prompts`: List of prompts referenced by selected guidelines

**Source (local type):**
- `name`: Unique source identifier (derived from path or custom via `--name`)
- `type`: Must be `"local"`
- `path`: Absolute path to local directory
- `guidelines`: List of selected guidelines from this source
- `prompts`: List of prompts referenced by selected guidelines

**Guideline:**
- `name`: Guideline identifier
- `file`: Relative path to guideline file (from source root)
- `description`: Brief description
- `applicable_scenarios`: List of scenarios where guideline applies
- `prompts`: List of prompt names referenced by this guideline

**Prompt:**
- `name`: Prompt identifier
- `file`: Relative path to prompt file (from source root)
- `description`: Brief description

### Source Name Derivation

When you don't specify `--name`, DNASpec automatically derives a source name:

**From git URL:**
- `https://github.com/company/dna-guidelines.git` → `dna-guidelines`
- `git@github.com:company/my-patterns.git` → `my-patterns`

**From local path:**
- `/Users/me/my-dna-patterns` → `my-dna-patterns`
- `/path/to/Company_DNA` → `company-dna`

Names are sanitized to lowercase with hyphens, removing special characters.

## Using AI Agents with DNA

After running `dnaspec update-agents`, your DNA guidelines become available to AI coding assistants through generated files.

### Claude Code Integration

**Generated files:**
- `CLAUDE.md` - Contains managed block with guideline references
- `.claude/commands/dnaspec/<source-name>-<prompt-name>.md` - Slash commands for each prompt

**Usage:**
1. Claude Code automatically reads `CLAUDE.md` when assisting in your project
2. Use slash commands to trigger guideline-based reviews: `/dnaspec-<source>-<prompt>`
3. Claude will reference the appropriate guidelines when executing prompts

**Example:**
```
Type: /dnaspec-company-dna-code-review

Claude will review your code against the guidelines in:
@/dnaspec/company-dna/guidelines/go-style.md
```

### GitHub Copilot Integration

**Generated files:**
- `AGENTS.md` - Contains managed block with guideline references
- `.github/prompts/dnaspec-<source-name>-<prompt-name>.prompt.md` - Copilot prompts

**Usage:**
1. GitHub Copilot reads prompts from `.github/prompts/` automatically
2. Invoke prompts in Copilot Chat
3. Copilot will apply the guidelines specified in the prompts

## Troubleshooting

### "dnaspec.yaml already exists"

**Problem:** You ran `dnaspec init` but a configuration file already exists.

**Solution:**
- If you want to keep the existing file, use `dnaspec add` to add more sources
- If you want to start fresh, rename or delete the existing file first:
  ```bash
  mv dnaspec.yaml dnaspec.yaml.bak
  dnaspec init
  ```

### "source with name 'X' already exists"

**Problem:** You're trying to add a source with a name that's already in your configuration.

**Solution:** Use the `--name` flag to specify a different name:
```bash
dnaspec add --git-repo https://github.com/company/dna --name company-dna-v2
```

### "git clone failed"

**Problem:** Unable to clone the git repository.

**Solutions:**
- **Network issues**: Check your internet connection
- **Authentication**: For private repositories, ensure you have SSH keys set up or use HTTPS with credentials
- **Invalid URL**: Verify the repository URL is correct
- **Timeout**: Large repositories may timeout; try using `--git-ref` to specify a tag/branch

### "dnaspec-manifest.yaml not found"

**Problem:** The source directory doesn't contain a valid DNASpec manifest.

**Solution:**
- Verify you're pointing to the correct repository/directory
- Check that the repository is actually a DNA guidelines repository
- If you're a maintainer, run `dnaspec manifest init` in that directory first

### "guideline 'X' not found in source"

**Problem:** You specified a guideline with `--guideline` that doesn't exist in the source.

**Solution:**
- Remove the `--guideline` flag to see all available guidelines interactively
- Check the source's manifest for correct guideline names
- Fix any typos in the guideline name

### "failed to copy file"

**Problem:** Unable to copy files to your project directory.

**Solutions:**
- **Permission denied**: Check that you have write permissions in the current directory
- **Disk full**: Ensure you have sufficient disk space
- **Path too long**: On Windows, file paths might exceed maximum length

### "git clone timed out"

**Problem:** The repository clone operation took longer than 5 minutes.

**Solutions:**
- Check your network connection
- Try cloning the repository manually first to diagnose issues
- For very large repositories, consider using a local directory instead:
  ```bash
  git clone https://github.com/company/dna /tmp/dna
  dnaspec add /tmp/dna
  ```

### "no sources configured"

**Problem:** You ran `dnaspec update-agents` but haven't added any DNA sources yet.

**Solution:** Add at least one DNA source first:
```bash
dnaspec add --git-repo https://github.com/company/dna-guidelines
```

### "failed to create directory"

**Problem:** Unable to create agent configuration directories (`.claude/commands/`, `.github/prompts/`).

**Solutions:**
- **Permission denied**: Check that you have write permissions in the current directory
- **Path conflicts**: Ensure no files exist with the same names as the directories
- **Disk full**: Ensure you have sufficient disk space

### "managed block corrupted"

**Problem:** The managed block markers in AGENTS.md or CLAUDE.md are incomplete or malformed.

**Solution:**
- Manually fix the markers to ensure they appear in pairs:
  ```markdown
  <!-- DNASPEC:START -->
  ...content...
  <!-- DNASPEC:END -->
  ```
- Or delete the file and run `dnaspec update-agents` again to regenerate it

## Examples

### Setting up a new project

```bash
# 1. Initialize project configuration
dnaspec init

# 2. Add DNA guidelines from company repo
dnaspec add --git-repo https://github.com/company/dna-guidelines --git-ref v1.0.0

# (Interactive: select guidelines)

# 3. Configure AI agents
dnaspec update-agents

# (Interactive: select Claude Code, GitHub Copilot)

# 4. Commit to version control
git add dnaspec.yaml dnaspec/ AGENTS.md CLAUDE.md .claude/ .github/
git commit -m "Add DNA guidelines"
```

### Adding guidelines from multiple sources

```bash
# Add company-wide guidelines
dnaspec add --git-repo https://github.com/company/dna --name company

# Add team-specific patterns
dnaspec add --git-repo https://github.com/team/patterns --name team

# Add local experimental guidelines
dnaspec add /Users/me/experimental-dna --name experimental

# Update agent files
dnaspec update-agents --no-ask
```

### Updating to latest guidelines

```bash
# Option 1: Update specific source
dnaspec update company
# (Shows changes, prompts for new guidelines)
dnaspec update-agents --no-ask

# Option 2: Update all sources
dnaspec update --all
dnaspec update-agents --no-ask

# Option 3: Sync everything (update all + regenerate agents)
dnaspec sync

# Commit changes
git add dnaspec.yaml dnaspec/ AGENTS.md CLAUDE.md .claude/ .github/
git commit -m "Update DNA guidelines to latest"
```

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test package
go test ./internal/core/validate/...
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/aviator5/dnaspec.git
cd dnaspec

# Build
go build -o dnaspec ./cmd/dnaspec

# Install
go install ./cmd/dnaspec
```
