# Proposal: Use Relative Paths for Local Sources

**Change ID:** `use-relative-local-paths`
**Status:** Implemented
**Created:** 2025-12-09
**Implemented:** 2025-12-09

## Problem Statement

Currently, DNASpec stores absolute paths for local-path sources in `dnaspec.yaml`. This creates portability and collaboration issues:

1. **Not Machine-Portable**: Absolute paths like `/Users/vasily/dev/dnaspec/examples/dna` only work on the machine where they were created
2. **Version Control Issues**: Committing absolute paths to git means other team members can't use the configuration
3. **Security Concerns**: Absolute paths prevent validation that local sources remain within the project directory
4. **Inconsistent with Best Practices**: Most tools (git, npm, etc.) use relative paths for local references

As noted in the user request: "when local directory is used, it's better to store relative path since it should be the same on different machines. Additionally, that path should [not] leave the project dir."

## Current Behavior

**Current `dnaspec.yaml` example:**
```yaml
sources:
  - name: dna
    type: local-path
    path: /Users/vasily/dev/dnaspec/examples/dna  # Absolute path!
```

**Problems:**
- Won't work on other developers' machines
- Can reference paths outside the project directory
- Makes collaborative development difficult

## Proposed Solution

**Store relative paths instead of absolute paths** for local-path sources:

```yaml
sources:
  - name: dna
    type: local-path
    path: examples/dna  # Relative to project root!
```

**Key principles:**
1. **Relative to project root**: Paths are relative to the directory containing `dnaspec.yaml`
2. **Validation**: Ensure paths don't escape the project directory (no `..` traversal outside project root)
3. **Backward compatibility**: Support reading old configs with absolute paths (with warnings)
4. **Automatic conversion**: When adding local sources, automatically convert absolute paths to relative

## Benefits

1. **Portability**: Configuration works on any machine where the project is checked out
2. **Collaboration**: Team members can share configurations via version control
3. **Security**: Prevents accidental or malicious references to system directories
4. **Consistency**: Aligns with best practices from other development tools
5. **Monorepo-friendly**: Supports using local DNA sources within monorepos

## Scope

### In Scope

1. **Configuration format**: Use relative paths in `dnaspec.yaml` for `local-path` sources
2. **Path validation**: Ensure local paths don't escape project directory
3. **Path resolution**: Resolve relative paths based on `dnaspec.yaml` location
4. **Migration guidance**: Update documentation and examples
5. **Backward compatibility**: Commands continue working with absolute paths (with warnings), but `validate` is strict

### Out of Scope

1. Git repository paths (remain as URLs)
2. Symlink handling (existing `symlinked` field unchanged)
3. Manifest file paths (already relative within DNA sources)

## User Impact

### Breaking Changes

**None** - Full backward compatibility maintained:
- Old configs with absolute paths continue to work with all commands **without warnings**
- **`dnaspec validate` shows warnings** (not errors) for absolute paths - validation still passes
- New configs automatically use relative paths
- **`dnaspec add` shows warning and confirmation** for NEW sources outside project (before parsing manifest)
- All other commands (`list`, `update`, `sync`) operate silently with absolute paths

### Migration Path

For users with existing absolute path configurations:

1. **Detection**: `dnaspec validate` shows **warnings** for absolute paths and suggests running `update`
2. **Automatic**: Running `dnaspec update <source>` auto-converts to relative paths when possible (silently keeps absolute if outside project)
3. **Manual**: Users can edit `dnaspec.yaml` to convert absolute to relative paths
4. **Seamless**: All commands (`list`, `add`, `update`, `sync`) work silently with absolute paths - no warnings

## Examples

### Example 1: Adding Local Source

**Before:**
```bash
$ dnaspec add /Users/me/myproject/local-dna
# Saved: /Users/me/myproject/local-dna
```

**After:**
```bash
$ dnaspec add /Users/me/myproject/local-dna
# Saved as: local-dna (relative path)
```

### Example 2: Monorepo Structure

```
myproject/
├── dnaspec.yaml
├── packages/
│   ├── api/
│   └── web/
└── shared-dna/
    ├── dnaspec-manifest.yaml
    └── guidelines/
```

**Configuration:**
```yaml
sources:
  - name: shared
    type: local-path
    path: shared-dna  # Relative path works!
```

### Example 3: Adding Path Outside Project (Early Warning)

```bash
$ dnaspec add ../../system-files

⚠ Warning: Local source is outside project directory
  Project: /Users/me/myproject
  Source: /Users/me/system-files

This absolute path won't work on other machines.
Consider moving the source into your project directory.

Continue with absolute path? (y/N): n
# Cancelled by user
```

**Key behavior**: The warning and confirmation prompt appears **before** parsing `dnaspec-manifest.yaml`, preventing wasted time if user decides to cancel.

### Example 4: Validation Warnings

```bash
$ dnaspec validate

Validating dnaspec.yaml...
✓ YAML syntax valid
✓ Version 1 schema valid
✓ 1 sources configured
✓ All referenced files exist:
  - dnaspec/legacy-dna/guidelines/pattern.md

⚠ Found 1 warning(s):
  - Source 'legacy-dna' uses absolute path: /Users/me/external/dna
    Run 'dnaspec update legacy-dna' to auto-convert, or manually edit dnaspec.yaml

✓ Configuration is valid (with warnings)
```

## Alternatives Considered

### Alternative 1: Keep Absolute Paths Only

**Rejected because:**
- Doesn't solve portability issues
- Continues to create collaboration friction
- Doesn't align with user request or best practices

### Alternative 2: Support Both Absolute and Relative

**Rejected because:**
- Creates confusion about which to use
- Still allows absolute paths with their portability issues
- Better to have one clear approach

### Alternative 3: Use Environment Variables

**Rejected because:**
- Adds complexity for simple use case
- Doesn't solve the collaboration problem
- Overkill for local path references

## Success Criteria

1. ✅ New local sources added via CLI use relative paths
2. ✅ Relative paths resolve correctly from project root
3. ✅ Validation prevents paths escaping project directory
4. ✅ Documentation and examples updated
5. ✅ Backward compatibility maintained with warnings
6. ✅ Tests cover edge cases (symlinks, `.` and `..`, nested paths)

## Dependencies

None - this is an internal change to path handling.

## Timeline

**Estimated effort**: Small (2-4 hours)

**Phases:**
1. Update path validation logic
2. Update add/update commands to use relative paths
3. Update documentation and examples
4. Add tests for new behavior

## Open Questions

1. **Q**: Should we allow `.` prefix for relative paths (e.g., `./examples/dna`)?
   **A**: Yes, but normalize to remove `.` prefix for consistency

2. **Q**: How to handle symlinks that point outside the project?
   **A**: Validate the resolved path, not the symlink path

3. **Q**: Should `dnaspec validate` error or warn on absolute paths?
   **A**: Warn (not error). This maintains full backward compatibility - configs with absolute paths still pass validation. All other commands work silently with absolute paths to avoid noise.

4. **Q**: When should we show the "path outside project" warning in `dnaspec add`?
   **A**: Show warning and prompt for confirmation **before** loading the source or parsing `dnaspec-manifest.yaml`. This prevents wasted time if the user decides to cancel.

5. **Q**: Should existing commands show warnings when loading configs with absolute paths?
   **A**: No. Only `dnaspec add` warns when adding NEW sources. All other commands (`list`, `update`, `sync`) operate silently. Only `dnaspec validate` shows warnings.
