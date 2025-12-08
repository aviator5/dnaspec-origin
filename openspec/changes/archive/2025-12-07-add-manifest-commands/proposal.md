# Change: Add Manifest Management Commands

## Why

DNA repository maintainers need tools to create and validate their `dnaspec-manifest.yaml` files. These commands enable DNA repository authors to initialize a manifest with proper structure and validate it before publishing, ensuring that consumers can successfully integrate the DNA guidelines into their projects.

## What Changes

- Add `dnaspec manifest` command group with subcommands
- Add `dnaspec manifest init` subcommand to scaffold a new manifest file with example structure
- Add `dnaspec manifest validate` subcommand to validate manifest structure, file references, and constraints
- Implement manifest schema validation (version, required fields, naming conventions)
- Implement file reference validation (guidelines and prompts files must exist)
- Implement cross-reference validation (prompt references from guidelines must exist)
- Support spinal-case naming convention for guidelines and prompts
- Validate that `applicable_scenarios` is not empty (critical for AGENTS.md generation)

## Impact

- **Affected specs**: Creates new capability `manifest-management`
- **Affected code**:
  - New CLI command group under `internal/cli/manifest/` (manifest.go, init.go, validate.go)
  - New core logic under `internal/core/config/` for manifest schema and validation
  - New validation engine under `internal/core/validate/` for manifest validation rules
  - Main entry point in `cmd/dnaspec/main.go`
- **Users affected**: DNA repository maintainers
- **No breaking changes**: This is new functionality
