# DNASpec

DNASpec is a tool that helps project developers integrate DNA (Development & Architecture) guidelines into their projects, and helps DNA repository maintainers create and validate manifest files.

## Installation

```bash
go install github.com/aviator5/dnaspec/cmd/dnaspec@latest
```

## Quick Start

### For Project Developers

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

### For DNA Repository Maintainers

1. Initialize a new manifest:
```bash
dnaspec manifest init
```

2. Edit `dnaspec-manifest.yaml` to add your guidelines and prompts

3. Create the referenced files in `guidelines/` and `prompts/` directories

4. Validate your manifest:
```bash
dnaspec manifest validate
```

## Commands

### Project Commands

Project commands help you integrate DNA guidelines into your projects.

#### `dnaspec init`

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

#### `dnaspec add`

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

#### `dnaspec remove`

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

#### `dnaspec list`

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

local-patterns (local-dir)
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

#### `dnaspec update`

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

#### `dnaspec update-agents`

Generate or update AI agent configuration files based on selected DNA guidelines.

```bash
# Interactive mode - select which agents to configure
dnaspec update-agents

# Non-interactive mode - use saved agent configuration
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

### Manifest Commands

Manifest commands help DNA repository maintainers create and validate manifests.

#### `dnaspec manifest init`

Initialize a new `dnaspec-manifest.yaml` file with example structure.

```bash
dnaspec manifest init
```

This command:
- Creates a new manifest file in the current directory
- Includes example guidelines and prompts with helpful comments
- Prevents overwriting an existing manifest file

**Example output:**
```
✓ Success: Created dnaspec-manifest.yaml

Next steps:
  1. Edit the manifest file to add your guidelines and prompts
  2. Create the referenced files in guidelines/ and prompts/ directories
  3. Run dnaspec manifest validate to check your manifest
```

### `dnaspec manifest validate`

Validate the `dnaspec-manifest.yaml` file in the current directory.

```bash
dnaspec manifest validate
```

This command checks:
- **Manifest structure**: Validates required fields (version, guidelines, prompts)
- **Guideline definitions**: Ensures all guidelines have name, file, description, and applicable_scenarios
- **Prompt definitions**: Ensures all prompts have name, file, and description
- **File references**: Verifies that all referenced files exist
- **Cross-references**: Checks that prompts referenced by guidelines are defined
- **Naming conventions**: Enforces spinal-case (lowercase with hyphens)
- **Path security**: Prevents absolute paths and path traversal attacks
- **Applicable scenarios**: Ensures guidelines have at least one applicable scenario (required for AGENTS.md generation)

**Example output (success):**
```
✓ Manifest is valid
```

**Example output (with errors):**
```
✗ Found 3 validation error(s):

  • guidelines[0].name: invalid naming format: 'MyGuideline' (expected spinal-case: lowercase letters and hyphens only)
  • guidelines[0].file: file not found: guidelines/missing.md
  • guidelines[1].prompts: guideline 'api-design' references non-existent prompt 'review'

Fix these errors and run dnaspec manifest validate again.
```

## Project Configuration Format

The `dnaspec.yaml` file in your project directory tracks which DNA sources you've added and which guidelines are active.

```yaml
version: 1

agents:
  - "claude-code"
  - "github-copilot"

sources:
  - name: "company-dna"
    type: "git-repo"
    url: "https://github.com/company/dna-guidelines"
    ref: "v1.2.0"
    commit: "abc123def456789..."
    guidelines:
      - name: "go-style"
        file: "guidelines/go-style.md"
        description: "Go code style conventions"
        applicable_scenarios:
          - "writing new Go code"
        prompts: ["code-review"]
    prompts:
      - name: "code-review"
        file: "prompts/code-review.md"
        description: "Review Go code"
        
  - name: "local-patterns"
    type: "local"
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

## Manifest File Format

The `dnaspec-manifest.yaml` file defines your project's guidelines and prompts:

```yaml
version: 1

guidelines:
  - name: go-style
    file: guidelines/go-style.md
    description: Go coding style guidelines
    applicable_scenarios:
      - Writing Go code
      - Code reviews
    prompts:
      - code-review
      - documentation

  - name: rest-api
    file: guidelines/rest-api.md
    description: REST API design principles
    applicable_scenarios:
      - Designing APIs
      - API documentation

prompts:
  - name: code-review
    file: prompts/code-review.md
    description: Code review checklist

  - name: documentation
    file: prompts/documentation.md
    description: Documentation standards
```

### Required Fields

**Manifest:**
- `version`: Must be `1`

**Guideline:**
- `name`: Unique identifier in spinal-case (e.g., `go-style`)
- `file`: Relative path starting with `guidelines/`
- `description`: Brief description of the guideline
- `applicable_scenarios`: List of scenarios where this guideline applies (at least one required)

**Prompt (optional):**
- `prompts`: List of prompt names that complement this guideline

**Prompt:**
- `name`: Unique identifier in spinal-case (e.g., `code-review`)
- `file`: Relative path starting with `prompts/`
- `description`: Brief description of the prompt

### Naming Conventions

Names must follow **spinal-case**:
- ✓ Valid: `go-style`, `rest-api`, `code-review-123`
- ✗ Invalid: `GoStyle` (camelCase), `go_style` (snake_case), `Go-Style` (uppercase)

### File Paths

File paths must:
- Be relative (not absolute)
- Start with `guidelines/` or `prompts/`
- Not contain path traversal (`..`)

Examples:
- ✓ Valid: `guidelines/go-style.md`, `prompts/review.md`
- ✗ Invalid: `/etc/passwd`, `../other/file.md`, `guidelines/../../etc/passwd`

## Validation Rules

### Structure Validation
- Version must be specified and equal to 1
- Guidelines and prompts arrays must be present (can be empty)

### Guideline Validation
- All required fields must be present
- Names must be unique across all guidelines
- Names must use spinal-case format
- File paths must follow security rules
- Referenced files must exist
- Must have at least one applicable scenario

### Prompt Validation
- All required fields must be present
- Names must be unique across all prompts
- Names must use spinal-case format
- File paths must follow security rules
- Referenced files must exist

### Cross-Reference Validation
- Any prompt referenced in a guideline's `prompts` field must be defined in the `prompts` section

## Troubleshooting

### Project Commands

#### "dnaspec.yaml already exists"

**Problem:** You ran `dnaspec init` but a configuration file already exists.

**Solution:**
- If you want to keep the existing file, use `dnaspec add` to add more sources
- If you want to start fresh, rename or delete the existing file first:
  ```bash
  mv dnaspec.yaml dnaspec.yaml.bak
  dnaspec init
  ```

#### "source with name 'X' already exists"

**Problem:** You're trying to add a source with a name that's already in your configuration.

**Solution:** Use the `--name` flag to specify a different name:
```bash
dnaspec add --git-repo https://github.com/company/dna --name company-dna-v2
```

#### "git clone failed"

**Problem:** Unable to clone the git repository.

**Solutions:**
- **Network issues**: Check your internet connection
- **Authentication**: For private repositories, ensure you have SSH keys set up or use HTTPS with credentials
- **Invalid URL**: Verify the repository URL is correct
- **Timeout**: Large repositories may timeout; try using `--git-ref` to specify a tag/branch

#### "dnaspec-manifest.yaml not found"

**Problem:** The source directory doesn't contain a valid DNASpec manifest.

**Solution:**
- Verify you're pointing to the correct repository/directory
- Check that the repository is actually a DNA guidelines repository
- If you're a maintainer, run `dnaspec manifest init` in that directory first

#### "guideline 'X' not found in source"

**Problem:** You specified a guideline with `--guideline` that doesn't exist in the source.

**Solution:**
- Remove the `--guideline` flag to see all available guidelines interactively
- Check the source's manifest for correct guideline names
- Fix any typos in the guideline name

#### "failed to copy file"

**Problem:** Unable to copy files to your project directory.

**Solutions:**
- **Permission denied**: Check that you have write permissions in the current directory
- **Disk full**: Ensure you have sufficient disk space
- **Path too long**: On Windows, file paths might exceed maximum length

#### "git clone timed out"

**Problem:** The repository clone operation took longer than 5 minutes.

**Solutions:**
- Check your network connection
- Try cloning the repository manually first to diagnose issues
- For very large repositories, consider using a local directory instead:
  ```bash
  git clone https://github.com/company/dna /tmp/dna
  dnaspec add /tmp/dna
  ```

#### "no sources configured"

**Problem:** You ran `dnaspec update-agents` but haven't added any DNA sources yet.

**Solution:** Add at least one DNA source first:
```bash
dnaspec add --git-repo https://github.com/company/dna-guidelines
```

#### "failed to create directory"

**Problem:** Unable to create agent configuration directories (`.claude/commands/`, `.github/prompts/`).

**Solutions:**
- **Permission denied**: Check that you have write permissions in the current directory
- **Path conflicts**: Ensure no files exist with the same names as the directories
- **Disk full**: Ensure you have sufficient disk space

#### "managed block corrupted"

**Problem:** The managed block markers in AGENTS.md or CLAUDE.md are incomplete or malformed.

**Solution:**
- Manually fix the markers to ensure they appear in pairs:
  ```markdown
  <!-- DNASPEC:START -->
  ...content...
  <!-- DNASPEC:END -->
  ```
- Or delete the file and run `dnaspec update-agents` again to regenerate it

### Manifest Commands

#### "Manifest file already exists"

**Problem:** You ran `dnaspec manifest init` but a manifest file already exists.

**Solution:**
- If you want to keep the existing file, edit it directly instead of running `init`
- If you want to start fresh, rename or delete the existing file first:
  ```bash
  mv dnaspec-manifest.yaml dnaspec-manifest.yaml.bak
  dnaspec manifest init
  ```

### "file not found: guidelines/..."

**Problem:** Your manifest references a file that doesn't exist.

**Solution:** Create the missing file:
```bash
mkdir -p guidelines
touch guidelines/your-guideline.md
```

Or update the manifest to reference an existing file.

### "invalid naming format"

**Problem:** A guideline or prompt name doesn't follow spinal-case conventions.

**Solution:** Rename to use only lowercase letters, numbers, and hyphens:
- Change `MyGuideline` to `my-guideline`
- Change `API_Design` to `api-design`
- Change `codeReview` to `code-review`

### "path traversal not allowed"

**Problem:** A file path contains `..` which could reference files outside the intended directories.

**Solution:** Use relative paths within `guidelines/` or `prompts/`:
- Change `../../etc/passwd` to a proper path like `guidelines/security.md`
- Remove any `..` path components

### "references non-existent prompt"

**Problem:** A guideline references a prompt that isn't defined in the `prompts` section.

**Solution:** Either:
- Add the missing prompt to the `prompts` section
- Remove the prompt reference from the guideline
- Fix the prompt name if it's a typo

### "empty applicable_scenarios"

**Problem:** A guideline has an empty `applicable_scenarios` list.

**Solution:** Add at least one scenario where this guideline applies:
```yaml
guidelines:
  - name: my-guideline
    file: guidelines/my-guideline.md
    description: My guideline
    applicable_scenarios:
      - "When writing new features"
      - "During code reviews"
```

The `applicable_scenarios` field is required because it's used to generate AGENTS.md files that help AI assistants understand when to apply your guidelines.

## Examples

See the `examples/` directory for complete example manifests:

- `examples/minimal-manifest.yaml`: Minimal valid manifest with one guideline
- `examples/complete-manifest.yaml`: Full-featured manifest with multiple guidelines and prompts
- `examples/go-project-manifest.yaml`: Example for a Go project

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

## License

[License information to be added]
