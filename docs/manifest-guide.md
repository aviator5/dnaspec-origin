# DNASpec Manifest Guide

Complete guide for creating and maintaining DNA repositories.

## Table of Contents

- [What is a DNA Repository](#what-is-a-dna-repository)
- [Getting Started](#getting-started)
- [Manifest Commands Reference](#manifest-commands-reference)
  - [dnaspec manifest init](#dnaspec-manifest-init)
  - [dnaspec manifest validate](#dnaspec-manifest-validate)
- [Manifest Configuration](#manifest-configuration)
- [Creating Guidelines](#creating-guidelines)
- [Creating Prompts](#creating-prompts)
- [Guidelines and Prompts Best Practices](#guidelines-and-prompts-best-practices)
- [Validation Rules](#validation-rules)
- [Naming Conventions](#naming-conventions)
- [Publishing DNA Repositories](#publishing-dna-repositories)
- [Troubleshooting](#troubleshooting)
- [Examples](#examples)

## What is a DNA Repository

A DNA repository is a git repository or local directory that contains:
- **Guidelines**: Markdown files defining development standards and architectural patterns
- **Prompts**: Instructions for AI agents on how to apply the guidelines
- **Manifest**: A `dnaspec-manifest.yaml` file that declares available guidelines and prompts

DNA repositories serve as centralized sources of reusable development knowledge that can be integrated into multiple projects. They enable teams to:
- Share best practices across projects
- Maintain consistent coding standards
- Provide AI agents with context-aware guidance
- Version control development guidelines with git

## Getting Started

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

## Manifest Commands Reference

### `dnaspec manifest init`

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

## Manifest Configuration

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

### File Paths

File paths must:
- Be relative (not absolute)
- Start with `guidelines/` or `prompts/`
- Not contain path traversal (`..`)

Examples:
- ✓ Valid: `guidelines/go-style.md`, `prompts/review.md`
- ✗ Invalid: `/etc/passwd`, `../other/file.md`, `guidelines/../../etc/passwd`

## Creating Guidelines

Guidelines are markdown files that define development standards, architectural patterns, and best practices.

### Guideline Structure

```markdown
# [Guideline Name]

## Overview
Brief description of what this guideline covers.

## When to Use
- Scenario 1
- Scenario 2

## Standards/Patterns
Detailed content about the guideline.

### Example
Code or configuration examples.

## References
Links to additional resources.
```

### Best Practices for Guidelines

1. **Be Specific**: Provide concrete, actionable guidance
2. **Use Examples**: Include code examples showing correct and incorrect usage
3. **Explain Why**: Don't just say "do this" - explain the reasoning
4. **Keep Focused**: One guideline per topic (e.g., separate guidelines for style vs. architecture)
5. **Link Related Guidelines**: Reference other guidelines when relevant
6. **Update Regularly**: Keep guidelines current with evolving best practices

### Example Guideline

```markdown
# Go Code Style

## Overview
Go coding style conventions for writing clean, idiomatic Go code.

## When to Use
- Writing new Go code
- Refactoring existing Go code
- Code reviews

## Naming Conventions

### Variables
- Use `camelCase` for local variables
- Use `PascalCase` for exported names
- Use short names for small scopes: `i` for loop counter, `r` for reader

**Good:**
```go
func processUser(userID int) {
    u := getUserByID(userID)
    // ...
}
```

**Bad:**
```go
func processUser(user_id int) {  // Don't use snake_case
    UserObject := getUserByID(user_id)  // Don't capitalize local vars
    // ...
}
```

## Error Handling
Always check errors and handle them appropriately.

**Good:**
```go
data, err := ioutil.ReadFile(filename)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```

## References
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
```

## Creating Prompts

Prompts are instructions that tell AI agents how to apply guidelines. They become slash commands in Claude Code and chat prompts in GitHub Copilot.

### Prompt Structure

```markdown
[Clear instruction for what the AI agent should do]

Check for:
- Specific thing 1
- Specific thing 2
- Specific thing 3

Reference the [guideline-name] guideline at @/dnaspec/<source>/<guideline-file>.

Provide specific feedback with line numbers.
```

### Best Practices for Prompts

1. **Be Action-Oriented**: Start with a clear verb (Review, Validate, Suggest, Check)
2. **List Specific Checks**: Break down what the AI should look for
3. **Reference Guidelines**: Tell the AI which guidelines to consult
4. **Request Specificity**: Ask for line numbers, code examples, etc.
5. **Keep Concise**: Prompts should be short and focused

### Example Prompt

```markdown
Review the Go code against the go-style guideline.

Check for:
- Proper naming conventions (camelCase for locals, PascalCase for exports)
- Error handling patterns (checking all errors, using %w for wrapping)
- Code organization (package structure, function size)
- Documentation (public functions and types have comments)
- Idiomatic Go patterns

Reference the go-style guideline at @/dnaspec/company-dna/guidelines/go-style.md.

Provide specific feedback with line numbers and suggest corrections.
```

## Guidelines and Prompts Best Practices

### How Guidelines and Prompts Work Together

- **Guidelines** define the "what": standards, patterns, and best practices
- **Prompts** define the "how": actionable tasks that apply the guidelines
- A single guideline can have multiple prompts (e.g., review, validate, generate)
- A single prompt can reference multiple guidelines

### Examples of Guideline/Prompt Pairs

**Example 1: Code Style**
- Guideline: `go-style.md` - Defines Go coding conventions
- Prompts:
  - `go-code-review` - Review code for style compliance
  - `go-format-suggest` - Suggest formatting improvements

**Example 2: API Design**
- Guideline: `rest-api.md` - REST API design principles
- Prompts:
  - `api-design-review` - Review API design
  - `api-documentation` - Generate API documentation
  - `api-validation` - Validate API spec against guidelines

**Example 3: Security**
- Guideline: `security-practices.md` - Security best practices
- Prompts:
  - `security-audit` - Audit code for security issues
  - `security-checklist` - Generate security checklist

### When to Create a New Guideline vs. Prompt

**Create a new guideline when:**
- You have a distinct topic with its own standards
- The content is substantial enough to stand alone
- Multiple prompts might reference it

**Create a new prompt when:**
- You have a specific task AI agents should perform
- The task can be clearly defined and scoped
- You want users to be able to invoke it directly

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

## Naming Conventions

### Spinal-case (for guidelines and prompts)
- Use lowercase letters
- Separate words with hyphens
- Numbers are allowed
- Examples: `go-style`, `rest-api`, `database-migrations`, `error-handling`

**Valid names:**
- ✓ `go-style`
- ✓ `rest-api`
- ✓ `code-review-123`

**Invalid names:**
- ✗ `GoStyle` (camelCase)
- ✗ `go_style` (snake_case)
- ✗ `Go-Style` (uppercase)
- ✗ `go style` (spaces)

## Publishing DNA Repositories

### Initial Publication

```bash
# In your DNA repository
dnaspec manifest init

# Edit manifest and create guideline/prompt files
# ...

# Validate
dnaspec manifest validate

# Commit
git add dnaspec-manifest.yaml guidelines/ prompts/
git commit -m "Initial DNA guidelines"

# Tag
git tag v1.0.0

# Push
git remote add origin https://github.com/company/dna.git
git push -u origin main --tags
```

### Updating Published DNA

```bash
# Make changes to guidelines/prompts
# ...

# Update manifest if needed
# ...

# Validate
dnaspec manifest validate

# Commit
git add .
git commit -m "Update go-style guideline"

# Tag new version
git tag v1.1.0
git push --tags
```

### Versioning Strategy

**Use semantic versioning:**
- **Major (v2.0.0)**: Breaking changes (removing guidelines, changing structure)
- **Minor (v1.1.0)**: Adding new guidelines or prompts
- **Patch (v1.0.1)**: Fixing typos, clarifying content

**Best practices:**
- Tag stable releases
- Use branches for major versions (v1.x, v2.x)
- Document changes in CHANGELOG.md
- Announce breaking changes to users

## Troubleshooting

### "Manifest file already exists"

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

### Minimal Manifest

```yaml
version: 1

guidelines:
  - name: go-style
    file: guidelines/go-style.md
    description: Go coding style conventions
    applicable_scenarios:
      - "writing Go code"
    prompts:
      - code-review

prompts:
  - name: code-review
    file: prompts/code-review.md
    description: Review code for style
```

### Complete Manifest

```yaml
version: 1

guidelines:
  - name: go-style
    file: guidelines/go-style.md
    description: Go coding style conventions
    applicable_scenarios:
      - "writing new Go code"
      - "refactoring existing Go code"
      - "code reviews"
    prompts:
      - go-code-review
      - go-format-suggest

  - name: go-service
    file: guidelines/go-service.md
    description: Go service layering and structure
    applicable_scenarios:
      - "designing Go services"
      - "organizing Go packages"
    prompts:
      - architecture-review

  - name: rest-api
    file: guidelines/rest-api.md
    description: REST API design principles
    applicable_scenarios:
      - "designing API endpoints"
      - "implementing HTTP handlers"
    prompts:
      - api-design-review
      - api-documentation

prompts:
  - name: go-code-review
    file: prompts/go-code-review.md
    description: Review Go code against go-style guideline

  - name: go-format-suggest
    file: prompts/go-format-suggest.md
    description: Suggest Go code formatting improvements

  - name: architecture-review
    file: prompts/architecture-review.md
    description: Review service architecture

  - name: api-design-review
    file: prompts/api-design-review.md
    description: Review API design against REST principles

  - name: api-documentation
    file: prompts/api-documentation.md
    description: Generate API documentation
```

### Creating a DNA Repository from Scratch

```bash
# Create DNA repo structure
mkdir company-dna
cd company-dna
git init

# Initialize manifest
dnaspec manifest init

# Create directory structure
mkdir -p guidelines prompts

# Create a guideline
cat > guidelines/go-style.md << 'EOF'
# Go Code Style

## Overview
Go coding style conventions for writing clean, idiomatic Go code.

## When to Use
- Writing new Go code
- Refactoring existing Go code
- Code reviews

## Naming Conventions

### Variables
- Use `camelCase` for local variables
- Use `PascalCase` for exported names

### Functions
- Use descriptive names
- Start with verbs for functions that perform actions

## Error Handling
Always check errors and handle them appropriately.

**Good:**
```go
data, err := ioutil.ReadFile(filename)
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}
```

## References
- [Effective Go](https://golang.org/doc/effective_go)
EOF

# Create a prompt
cat > prompts/go-code-review.md << 'EOF'
Review the Go code against the go-style guideline.

Check for:
- Proper naming conventions
- Error handling patterns
- Code organization
- Documentation

Reference the go-style guideline.

Provide specific feedback with line numbers.
EOF

# Update manifest
cat > dnaspec-manifest.yaml << 'EOF'
version: 1

guidelines:
  - name: go-style
    file: guidelines/go-style.md
    description: Go coding style conventions
    applicable_scenarios:
      - "writing new Go code"
      - "refactoring existing Go code"
      - "code reviews"
    prompts:
      - go-code-review

prompts:
  - name: go-code-review
    file: prompts/go-code-review.md
    description: Review Go code against go-style guideline
EOF

# Validate
dnaspec manifest validate

# Commit and tag
git add .
git commit -m "Initial DNA guidelines"
git tag v1.0.0

# Push to GitHub
git remote add origin https://github.com/company/dna.git
git push -u origin main --tags
```
