# Proposal: Add Update Command

## Overview

Implement the `dnaspec update` command to enable users to update DNA sources from their origins (git repositories or local directories). This command will fetch the latest manifest, update selected guidelines, handle new/removed guidelines, and provide both interactive and non-interactive modes.

## Motivation

Currently, users can add DNA sources using `dnaspec add`, but there's no way to update existing sources when the upstream DNA repository changes. Users need to be able to:

1. Pull the latest changes from git-based DNA sources
2. Refresh local directory-based sources
3. Update guideline content and metadata
4. Handle new guidelines that become available
5. Clean up guidelines that are removed from source
6. Do this for a specific source or all sources at once

Without this command, users must manually remove and re-add sources to get updates, losing their configuration and requiring re-selection of guidelines.

## Scope

This change adds the `dnaspec update` command with the following capabilities:

**In Scope:**
- Update specific source by name: `dnaspec update <source-name>`
- Update all sources: `dnaspec update --all`
- Dry-run mode: `dnaspec update <source-name> --dry-run`
- Interactive handling of new guidelines (prompt user)
- Non-interactive modes: `--add-new=all` and `--add-new=none`
- Update git sources from configured ref
- Update local directory sources
- Detect and report updated, new, and removed guidelines
- Update configuration with new commit hashes
- Copy updated guideline and prompt files

**Out of Scope:**
- Changing the git ref for a source (user must edit config manually)
- Automatic drift detection of local modifications
- Merging local changes with upstream changes
- Updating symlinked sources (they update automatically)

## Approach

The implementation will follow these key principles:

1. **Reuse Existing Infrastructure**: Leverage existing source fetching, manifest parsing, and file copying code from the `add` command
2. **Clear Change Reporting**: Show users exactly what changed (updated, new, removed)
3. **Safe Defaults**: Interactive mode by default, require explicit flags for automation
4. **Atomic Operations**: Use atomic file writes to prevent corruption
5. **Proper Cleanup**: Always clean up temporary directories, even on error

The command will be implemented in three main phases:

### Phase 1: Update Single Source
Implement core update logic for updating a single named source, including:
- Source lookup and validation
- Fetching latest manifest
- Comparing current vs latest state
- Updating selected guidelines
- Handling new/removed guidelines
- Updating configuration

### Phase 2: Update All Sources
Extend to support `--all` flag that iterates over all configured sources and applies the single-source update logic to each.

### Phase 3: Polish and Testing
Add comprehensive tests, improve error messages, and ensure edge cases are handled.

## Dependencies

**Internal Dependencies:**
- Existing source fetching code (`internal/core/source/`)
- Manifest parsing and validation (`internal/core/config/`)
- File operations (`internal/core/files/`)
- Configuration management (`internal/core/config/`)
- UI helpers (`internal/ui/`)

**No External Dependencies**: All required functionality already exists in the codebase or standard library.

## Risk Assessment

**Low Risk Areas:**
- Fetching sources: Already implemented and tested in `add` command
- Manifest parsing: Existing validation ensures safety
- File copying: Atomic operations prevent corruption

**Medium Risk Areas:**
- Configuration updates: Must preserve all existing data while updating source metadata
- Temporary directory cleanup: Must handle both success and failure cases

**Mitigation Strategies:**
- Comprehensive testing of configuration updates
- Use defer for cleanup to ensure it runs even on error
- Validate manifest before making any changes
- Use dry-run mode for testing without side effects

## Alternatives Considered

### Alternative 1: Auto-update on every command
Automatically check for updates whenever any dnaspec command runs.

**Rejected because:**
- Adds latency to all commands
- Unexpected for users (should be explicit)
- Could cause conflicts if working offline

### Alternative 2: Make --add-new default to "all"
Automatically add all new guidelines without prompting.

**Rejected because:**
- Too aggressive, users might not want all guidelines
- Breaks principle of explicit control
- Interactive mode is safer default

### Alternative 3: Track local modifications
Detect if users modified files in dnaspec/ directory and prevent updates.

**Rejected because:**
- Adds complexity (drift detection)
- Files in dnaspec/ should be treated as managed (read-only)
- Design philosophy: explicit over automatic
- Users should use separate files for customizations

## Success Criteria

The implementation will be considered successful when:

1. ✅ Users can update a single source: `dnaspec update my-company-dna`
2. ✅ Users can update all sources: `dnaspec update --all`
3. ✅ Dry-run mode works: `dnaspec update <source> --dry-run`
4. ✅ Interactive mode prompts for new guidelines
5. ✅ Non-interactive modes work: `--add-new=all` and `--add-new=none`
6. ✅ Git sources update to latest commit at configured ref
7. ✅ Local directory sources refresh from filesystem
8. ✅ Configuration updates preserve existing data
9. ✅ Clear reporting of what changed
10. ✅ All tests pass
11. ✅ Documentation is complete

## Open Questions

None. The design is fully specified in docs/design.md and all implementation details are clear.
