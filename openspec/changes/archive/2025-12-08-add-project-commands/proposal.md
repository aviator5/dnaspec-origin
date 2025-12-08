# Proposal: Add Project Commands

## Overview

Implement core project-level commands for DNASpec that enable users to initialize project configuration and add DNA guidelines from various sources (git repositories and local directories).

## Problem

Currently, DNASpec only supports manifest-level commands (`dnaspec manifest init`, `dnaspec manifest validate`) for DNA repository maintainers. Users need project-level commands to:

1. Initialize a project with `dnaspec.yaml` configuration
2. Add DNA guidelines from git repositories or local directories
3. Manage which guidelines are active in their projects

Without these commands, users cannot use DNASpec to integrate DNA guidelines into their projects.

## Solution

Implement two essential project commands:

### 1. `dnaspec init`
- Create empty `dnaspec.yaml` in current directory
- Include commented examples for user guidance
- Display next steps (run `dnaspec add` to add DNA sources)

### 2. `dnaspec add`
- Support two source types: git repositories and local directories
- Clone git repositories to temporary directory, parse manifest
- Read manifest from local directories directly
- Interactive guideline selection (or non-interactive with flags)
- Copy selected files to `dnaspec/<source-name>/` directory
- Update `dnaspec.yaml` with source metadata and selected guidelines

## Scope

### In Scope
- Project configuration data model (`dnaspec.yaml` structure)
- `dnaspec init` command implementation
- `dnaspec add` command with git repository support
- `dnaspec add` command with local directory support
- Source name derivation algorithm
- Interactive guideline selection
- Non-interactive modes (`--all`, `--guideline` flags)
- File copying to project's `dnaspec/` directory
- Configuration validation (basic structure validation)

### Out of Scope
- `dnaspec update` command (future work)
- `dnaspec remove` command (future work)
- `dnaspec update-agents` command (future work)
- `dnaspec list` command (future work)
- `dnaspec validate` command (future work)
- `dnaspec sync` command (future work)
- Symlink support for local directories (can be added later)
- Git clone caching (optimization for future)
- Agent file generation (separate change)

## Dependencies

- Existing manifest validation logic from `manifest-management` spec
- Go cobra CLI framework (already in use)
- Go git library for cloning repositories
- YAML parsing libraries (already in use)

## Risks

1. **Git clone failures**: Network issues or authentication problems could prevent adding sources
   - Mitigation: Clear error messages, timeout handling, URL validation

2. **Large repositories**: Cloning large DNA repositories could be slow
   - Mitigation: Use shallow clones (`--depth=1`), show progress indicators

3. **Path conflicts**: Multiple sources could have conflicting file paths
   - Mitigation: Namespace files under `dnaspec/<source-name>/` directories

4. **Security**: Malicious manifests could reference dangerous file paths
   - Mitigation: Reuse existing path validation from `manifest-management` spec

## Success Criteria

1. Users can run `dnaspec init` to create initial configuration
2. Users can add DNA sources from git repositories
3. Users can add DNA sources from local directories
4. Users can select specific guidelines interactively
5. Users can add all guidelines non-interactively with `--all` flag
6. Selected guideline files are copied to `dnaspec/<source-name>/` directory
7. `dnaspec.yaml` accurately reflects added sources and guidelines
8. Error messages are clear and actionable

## References

- Design document: `docs/design.md` (Commands Reference section)
- Existing spec: `manifest-management` (for validation logic)
- Project conventions: `openspec/project.md`
