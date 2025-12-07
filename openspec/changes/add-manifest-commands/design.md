# Design: Manifest Management Commands

## Context

DNA repository maintainers need tools to create and validate manifest files. This is the first implementation of DNASpec CLI, establishing patterns for all future commands.

**Constraints:**
- Go 1.25 language version
- Must prevent path traversal security vulnerabilities
- Clear, actionable error messages for users
- Terminal-friendly output formatting

**Stakeholders:**
- DNA repository maintainers (primary users)
- Project developers who consume DNA guidelines

## Goals / Non-Goals

**Goals:**
- Scaffold valid manifest files quickly with `manifest init`
- Catch all manifest errors before publishing with `manifest validate`
- Prevent security vulnerabilities (path traversal attacks)
- Provide clear, actionable error messages
- Establish code organization patterns for future commands (grouped subcommands)

**Non-Goals:**
- Interactive manifest editing (future enhancement)
- Auto-fixing invalid manifests (user must fix manually)
- Manifest migration between versions (not needed for v1)
- CI/CD integration (explicitly postponed)

## Decisions

### Library Selection

**CLI Framework: Cobra**
- Industry standard for Go CLIs (kubectl, hugo, github-cli use it)
- Excellent help text generation
- Easy subcommand structure
- Well-maintained and stable

**YAML Processing: gopkg.in/yaml.v3**
- Standard YAML library for Go
- Preserves comments and formatting
- Supports struct tags for marshaling
- Better error messages than v2

**Terminal UI: Charmbracelet (lipgloss, huh)**
- Modern, consistent terminal output formatting
- lipgloss for styled text and colors
- huh for future interactive prompts (in other commands)
- Widely adopted in Go CLI ecosystem

**Git Operations: go-git/v5**
- Pure Go implementation (no git binary dependency)
- Needed for future `add` and `update` commands
- Include now to establish dependency pattern

**Testing: testify + moq**
- **github.com/stretchr/testify** for test assertions
  - Use `assert` package for regular assertions
  - Use `require` package for critical assertions that should halt test
  - Clear, readable test failures
- **github.com/matryer/moq** for mock generation
  - Lightweight, interface-based mocking
  - Generated mocks from interfaces
  - Simple, idiomatic Go code generation

### Directory Structure

Following [golang-standards/project-layout](https://github.com/golang-standards/project-layout):

```
dnaspec/
├── cmd/
│   └── dnaspec/               # Main application
│       └── main.go            # Entry point
├── internal/                  # Private application code
│   ├── cli/                   # CLI layer (command parsing, flags)
│   │   ├── root.go            # Root command setup
│   │   └── manifest/          # Manifest command group
│   │       ├── manifest.go    # Group command
│   │       ├── init.go        # manifest init subcommand
│   │       └── validate.go    # manifest validate subcommand
│   ├── core/                  # Core domain logic
│   │   ├── config/            # Configuration and schema
│   │   │   ├── manifest.go    # Manifest data structures
│   │   │   └── templates.go   # Template generation
│   │   └── validate/          # Validation engine
│   │       ├── validator.go   # Manifest validation logic
│   │       └── errors.go      # Validation errors
│   └── ui/                    # UI/formatting layer
│       └── styles.go          # Terminal output styles (lipgloss)
├── go.mod
├── go.sum
└── README.md
```

**Rationale:**
- `/cmd/dnaspec` - Main application entry point (standard Go layout)
- `/internal` - Private code not importable by other projects
- Separate CLI layer from core domain (testability, reusability)
- Grouped subcommands under manifest/ (Cobra best practice)
- Clean separation of concerns (CLI, core, UI)
- Follows golang-standards/project-layout conventions

### Path Security Validation

**Threat Model:** Malicious manifest could reference files outside intended directories
```yaml
guidelines:
  - name: "evil"
    file: "../../../etc/passwd"  # Attack!
```

**Defense Strategy:**
1. **Reject absolute paths:** Only allow relative paths
2. **Reject path traversal:** No ".." components allowed
3. **Enforce directory prefix:** Must start with "guidelines/" or "prompts/"
4. **Use filepath.Clean():** Normalize before validation

**Implementation:**
```go
func validateManifestPath(path string) error {
    if filepath.IsAbs(path) {
        return fmt.Errorf("absolute paths not allowed: %s", path)
    }
    if strings.Contains(path, "..") {
        return fmt.Errorf("path traversal not allowed: %s", path)
    }
    clean := filepath.Clean(path)
    if !strings.HasPrefix(clean, "guidelines/") &&
       !strings.HasPrefix(clean, "prompts/") {
        return fmt.Errorf("path must be within guidelines/ or prompts/: %s", path)
    }
    return nil
}
```

**Why this approach:**
- Simple, auditable security checks
- Defense in depth (multiple validation layers)
- Clear error messages for legitimate mistakes
- Prevents all common path traversal techniques

### Error Handling and Reporting

**Validation Error Structure:**
```go
type ValidationError struct {
    Field   string // e.g., "guidelines[0].name"
    Message string // Human-readable error
}

type ValidationResult struct {
    Errors []ValidationError
    Valid  bool
}
```

**Output Format:**
```
Validating dnaspec-manifest.yaml...
✗ Error: Guideline 'go-style' references non-existent prompt 'missing-prompt'
✗ Error: File not found: guidelines/missing.md
✗ Error: Duplicate guideline name 'go-service'

Validation failed with 3 errors
```

**Rationale:**
- Collect ALL errors in one pass (don't fail fast)
- Users can fix multiple issues at once
- Clear context for each error
- Visual indicators (✓/✗) for quick scanning

### Atomic File Writes

**Problem:** Process could crash during file write, leaving corrupt file

**Solution for manifest init:**
1. Write content to temporary file (e.g., `.dnaspec-manifest.yaml.tmp`)
2. On success, atomically rename to target filename
3. Rename is atomic on all major filesystems

```go
tmpFile := ".dnaspec-manifest.yaml.tmp"
if err := os.WriteFile(tmpFile, content, 0644); err != nil {
    return err
}
return os.Rename(tmpFile, "dnaspec-manifest.yaml")
```

### Manifest Template Content

**Include in template:**
- Version 1 schema
- Two example guidelines with all required fields
- One example prompt
- Helpful comments explaining each field
- Example of applicable_scenarios (critical for users to understand)

**Example structure:**
```yaml
version: 1

guidelines:
  - name: "example-guideline"
    file: "guidelines/example.md"
    description: "Example guideline description"
    applicable_scenarios:
      - "example scenario 1"
      - "example scenario 2"
    prompts:
      - "example-prompt"

prompts:
  - name: "example-prompt"
    file: "prompts/example.md"
    description: "Example prompt description"
```

## Risks / Trade-offs

### Risk: Overly strict validation blocks legitimate use cases

**Mitigation:**
- Start with strict rules, can relax later if needed
- Validation rules match security best practices
- Clear error messages help users understand constraints

### Risk: Path validation has bypass

**Mitigation:**
- Multiple validation layers (defense in depth)
- Use standard library functions (filepath.Clean, filepath.IsAbs)
- Add security-focused test cases
- Document security model in code comments

### Risk: Error messages not helpful enough

**Mitigation:**
- Include context (field name, guideline/prompt name)
- Suggest fixes where possible
- Test with real users during development
- Iterate on message clarity

## Migration Plan

Not applicable - this is initial implementation, no migration needed.

## Open Questions

None - implementation path is clear from design.md specification.
