# Project Context

## Purpose

DNASpec is a CLI tool for DNA repository maintainers to create and validate manifest files that define project guidelines and prompts.

**Key Goals:**
- Enable DNA repository maintainers to define reusable development guidelines
- Provide validation for manifest structure and file references
- Support spinal-case naming conventions for consistency
- Ensure security through path validation and traversal protection

## Tech Stack

**Language & Runtime:**
- Go 1.25

**Core Dependencies:**
- `github.com/spf13/cobra` - CLI framework and command structure
- `gopkg.in/yaml.v3` - YAML parsing for manifest files
- `github.com/charmbracelet/lipgloss` - Terminal UI styling and formatting
- `github.com/stretchr/testify` - Testing assertions and mocks

**Build & Distribution:**
- Standard Go toolchain (`go build`, `go install`)
- Distributed as single binary via `go install`

## Project Conventions

### Code Style

**Naming Conventions:**
- **Package names**: Lowercase, single word (e.g., `config`, `validate`, `manifest`)
- **Struct names**: PascalCase (e.g., `Manifest`, `ManifestGuideline`)
- **Function names**: PascalCase for exported, camelCase for private
- **Manifest identifiers**: **spinal-case** (lowercase with hyphens, e.g., `go-style`, `rest-api`)
  - This is a critical convention enforced by validation
  - Valid: `go-style`, `code-review-123`
  - Invalid: `GoStyle`, `go_style`, `Go-Style`

**File Organization:**
- `cmd/` - CLI entry points
- `internal/` - Internal packages not exposed to external users
  - `internal/core/` - Domain logic (config, validation)
  - `internal/cli/` - Command implementations
  - `internal/ui/` - Terminal UI and styling

**Error Handling:**
- Return errors explicitly, don't panic in library code
- Use descriptive error messages with context
- Wrap errors with additional context where helpful

### Architecture Patterns

DNASpec follows **clean architecture** with clear layer separation:

**Core Domain Layer** (`internal/core/`):
- Configuration management (manifest structs, YAML parsing)
- Domain entities (Manifest, ManifestGuideline, ManifestPrompt)
- Schema definitions and templates
- Validation logic (manifest validation, path security, naming conventions)
- **No dependencies on CLI, UI, or external services**

**CLI / UI Layer** (`internal/cli/`, `internal/ui/`):
- Command parsing and flag handling (cobra)
- Terminal output formatting (lipgloss)
- User interaction patterns
- Depends on core domain layer

**File Operations:**
- All critical file writes should be atomic (write to temp file, then rename)
- Prevents corrupted files from interrupted operations

**Key Patterns:**
- Struct-based configuration with YAML tags
- Validation functions return error lists for comprehensive feedback
- Template-based file generation (see `internal/core/config/templates.go`)

### Testing Strategy

**Framework:**
- Standard Go testing (`testing` package)
- `github.com/stretchr/testify` for assertions and test utilities

**Test Organization:**
- Tests live alongside implementation files (`*_test.go`)
- Test packages mirror implementation packages
- Example test files: `validator_test.go`, `manifest_test.go`, `init_test.go`

**Test Execution:**
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/core/validate/...
```

**Testing Focus Areas:**
- Validation logic (path security, naming conventions, manifest structure)
- Template generation
- Error handling and edge cases

### Git Workflow

**Branch Strategy:**
- Feature branches for new work (e.g., `add-manifest-commands`)
- Descriptive branch names based on feature/change

**Commit Messages:**
- Clear, concise descriptions
- Action-oriented (e.g., "Add openspec add-manifest-commands proposal")
- Use conventional commit style where appropriate

**Development Flow:**
1. Create feature branch
2. Implement changes with tests
3. Commit with descriptive messages
4. Merge when complete

## Domain Context

**DNA (Development Norms & Architecture):**
- Reusable development guidelines, patterns, and best practices
- Can be shared across multiple projects
- Examples: code style conventions, architectural patterns, API design guidelines

**Manifests:**
- YAML configuration files (`dnaspec-manifest.yaml`)
- Define available guidelines and prompts in a DNA repository
- Required fields: `version`, `guidelines`, `prompts`
- Must validate before use

**Guidelines:**
- Individual markdown files containing development guidance
- Each has: name (spinal-case), file path, description, applicable scenarios
- File paths must start with `guidelines/`
- Referenced in manifest for discovery and validation

**Prompts:**
- Reusable instructions for AI agents
- Each has: name (spinal-case), file path, description
- File paths must start with `prompts/`
- Can be referenced by guidelines

**Validation Rules:**
- All guideline/prompt names must be unique
- Names must use spinal-case format
- File paths must be relative and within expected directories
- No path traversal (`..`) allowed
- All referenced files must exist
- Guidelines must have at least one applicable scenario

## Important Constraints

### Security Constraints

**Path Traversal Protection:**
- Manifest file paths are validated to prevent directory traversal attacks
- Paths must be relative (no absolute paths like `/etc/passwd`)
- No `..` components allowed in paths
- Paths must be within `guidelines/` or `prompts/` directories
- Implemented in validation layer

**File Path Requirements:**
- ✓ Valid: `guidelines/go-style.md`, `prompts/review.md`
- ✗ Invalid: `/etc/passwd`, `../other/file.md`, `guidelines/../../etc/passwd`

### Naming Constraints

**Spinal-Case Requirement:**
- All guideline and prompt names MUST use spinal-case
- Only lowercase letters, numbers, and hyphens allowed
- Enforced during validation
- This is a **hard requirement** for consistency across DNA repositories

### File Organization Constraints

**Directory Structure:**
```
DNA Repository Root/
├── dnaspec-manifest.yaml  (required)
├── guidelines/            (required directory for guideline files)
│   └── *.md
└── prompts/               (required directory for prompt files)
    └── *.md
```

## External Dependencies

**None at runtime** - DNASpec is a self-contained CLI tool with no external service dependencies.

**Build Dependencies:**
- Go 1.25 toolchain
- Standard library packages
- Third-party Go modules (managed via `go.mod`):
  - cobra (CLI framework)
  - yaml.v3 (YAML parsing)
  - lipgloss (terminal styling)
  - testify (testing utilities)
