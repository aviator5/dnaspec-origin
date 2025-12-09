# Design: Use Relative Paths for Local Sources

**Change ID:** `use-relative-local-paths`

## Overview

This document describes the technical design for transitioning from absolute paths to relative paths for local-path sources in DNASpec configurations.

## Architecture

### Current Implementation

```
User runs: dnaspec add /Users/me/project/local-dna
    ↓
Store in dnaspec.yaml:
    path: /Users/me/project/local-dna (absolute)
    ↓
When reading:
    Use path as-is
```

**Issues:**
- Path only works on original machine
- No validation against directory traversal
- Can't be shared via version control

### Proposed Implementation

```
User runs: dnaspec add /Users/me/project/local-dna
    ↓
Calculate relative path from project root:
    /Users/me/project/dnaspec.yaml (project root)
    /Users/me/project/local-dna
    → relative: local-dna
    ↓
Store in dnaspec.yaml:
    path: local-dna (relative)
    ↓
When reading:
    Resolve: project_root + relative_path
    Validate: resolved path is within project root
```

**Benefits:**
- Works on any machine
- Validates against directory traversal
- Shareable via version control

## Components

### 1. Path Utilities (New Package)

**Package:** `internal/core/paths`

**Purpose:** Centralize path handling logic

**Functions:**

```go
// MakeRelative converts an absolute path to relative from project root
// Returns error if path is not under project root
func MakeRelative(projectRoot, absPath string) (string, error)

// ResolveRelative converts a relative path to absolute based on project root
// Returns error if resolved path escapes project root
func ResolveRelative(projectRoot, relPath string) (string, error)

// ValidateLocalPath ensures a path is safe and within project root
// Works with both absolute and relative paths
func ValidateLocalPath(projectRoot, path string) error

// IsWithinProject checks if a path (after resolving symlinks) is within project root
func IsWithinProject(projectRoot, path string) (bool, error)
```

**Implementation Details:**

```go
func MakeRelative(projectRoot, absPath string) (string, error) {
    // 1. Clean both paths
    cleanRoot := filepath.Clean(projectRoot)
    cleanPath := filepath.Clean(absPath)

    // 2. Resolve symlinks
    realRoot, err := filepath.EvalSymlinks(cleanRoot)
    if err != nil {
        return "", err
    }
    realPath, err := filepath.EvalSymlinks(cleanPath)
    if err != nil {
        return "", err
    }

    // 3. Calculate relative path
    relPath, err := filepath.Rel(realRoot, realPath)
    if err != nil {
        return "", err
    }

    // 4. Validate doesn't escape (no leading ..)
    if strings.HasPrefix(relPath, "..") {
        return "", fmt.Errorf("path is outside project root")
    }

    // 5. Normalize ./ prefix
    relPath = strings.TrimPrefix(relPath, "./")

    return relPath, nil
}

func ResolveRelative(projectRoot, relPath string) (string, error) {
    // 1. Validate input is relative
    if filepath.IsAbs(relPath) {
        return "", fmt.Errorf("expected relative path, got absolute: %s", relPath)
    }

    // 2. Join with project root
    absPath := filepath.Join(projectRoot, relPath)

    // 3. Clean and resolve
    cleanPath := filepath.Clean(absPath)
    realPath, err := filepath.EvalSymlinks(cleanPath)
    if err != nil {
        // Path doesn't exist yet (ok for some operations)
        realPath = cleanPath
    }

    // 4. Validate within project root
    realRoot, _ := filepath.EvalSymlinks(filepath.Clean(projectRoot))
    if !strings.HasPrefix(realPath, realRoot+string(filepath.Separator)) &&
       realPath != realRoot {
        return "", fmt.Errorf("path escapes project root")
    }

    return realPath, nil
}
```

### 2. Configuration Updates

**File:** `internal/core/config/project.go`

**Changes:**

1. ~~Remove deprecation warnings from config loading~~ (loads silently)
2. Update save logic to use relative paths
3. Add migration helper

```go
// LoadProjectConfig loads and parses a project config file
func LoadProjectConfig(path string) (*ProjectConfig, error) {
    // ... existing load logic ...

    // Note: No warnings on load - only 'dnaspec validate' shows errors
    // All commands work silently with absolute paths to avoid noise

    return &config, nil
}

// MigrateToRelativePaths converts absolute paths to relative (in-place)
func (c *ProjectConfig) MigrateToRelativePaths(projectRoot string) error {
    for i := range c.Sources {
        source := &c.Sources[i]
        if source.Type == SourceTypeLocalPath && filepath.IsAbs(source.Path) {
            relPath, err := paths.MakeRelative(projectRoot, source.Path)
            if err != nil {
                return fmt.Errorf("source %s: %w", source.Name, err)
            }
            source.Path = relPath
        }
    }
    return nil
}
```

### 3. Command Updates

#### Add Command

**File:** `internal/cli/project/add.go`

**Changes:**

```go
// When adding local source
if localPath != "" {
    // Convert to absolute for validation
    absPath, err := filepath.Abs(localPath)
    if err != nil {
        return err
    }

    // EARLY CHECK: Validate path is within project BEFORE loading manifest
    projectRoot := filepath.Dir(configPath)
    relPath, err := paths.MakeRelative(projectRoot, absPath)
    if err != nil {
        // Path is outside project - warn and confirm BEFORE parsing manifest
        fmt.Fprintf(os.Stderr, "\n")
        fmt.Fprintf(os.Stderr, "⚠ Warning: Local source is outside project directory\n")
        fmt.Fprintf(os.Stderr, "  Project: %s\n", projectRoot)
        fmt.Fprintf(os.Stderr, "  Source:  %s\n", absPath)
        fmt.Fprintf(os.Stderr, "\n")
        fmt.Fprintf(os.Stderr, "This absolute path won't work on other machines.\n")
        fmt.Fprintf(os.Stderr, "Consider moving the source into your project directory.\n")
        fmt.Fprintf(os.Stderr, "\n")

        // Get confirmation BEFORE proceeding to load source
        if !confirm("Continue with absolute path?") {
            return fmt.Errorf("cancelled by user")
        }

        // Store as absolute path
        relPath = absPath
    }

    // NOW proceed to load source and parse manifest...
    sourceInfo, err := source.FetchLocalSource(absPath)
    // ... rest of add logic
}
```

**Key change**: The warning and confirmation prompt happens **before** `source.FetchLocalSource()` is called, so users don't waste time waiting for manifest parsing if they decide to cancel.

#### Update Command

**File:** `internal/cli/project/update.go`

**Changes:**

```go
// Auto-migrate absolute to relative during update
if source.Type == config.SourceTypeLocalPath && filepath.IsAbs(source.Path) {
    projectRoot := filepath.Dir(configPath)
    relPath, err := paths.MakeRelative(projectRoot, source.Path)
    if err == nil {
        fmt.Fprintf(os.Stderr, "✓ Converted to relative path: %s\n", relPath)
        source.Path = relPath
        needsSave = true
    }
}
```

#### Validate Command

**File:** `internal/cli/project/validate.go`

**Changes:**

```go
for _, source := range config.Sources {
    if source.Type == config.SourceTypeLocalPath {
        // Warn on absolute paths (not error - maintains backward compatibility)
        if filepath.IsAbs(source.Path) {
            warnings = append(warnings, fmt.Sprintf(
                "Source '%s' uses absolute path: %s\n" +
                "    Consider manually editing dnaspec.yaml to use a relative path",
                source.Name, source.Path,
            ))
        } else {
            // Validate relative path resolves within project
            projectRoot := filepath.Dir(configPath)
            if err := paths.ValidateLocalPath(projectRoot, source.Path); err != nil {
                errors = append(errors, fmt.Sprintf(
                    "Source '%s' path validation failed: %v",
                    source.Name, err,
                ))
            }
        }
    }
}
```

**Key change**: Uses **warnings** (not errors) for absolute paths. This maintains full backward compatibility - validation passes even with warnings. All other commands work silently with absolute paths.

### 4. Source Fetching Updates

**File:** `internal/core/source/fetch.go`

**Changes:**

```go
func FetchSource(source *config.ProjectSource, projectRoot string) (*SourceInfo, error) {
    switch source.Type {
    case config.SourceTypeLocalPath:
        // Resolve relative path
        sourcePath := source.Path
        if !filepath.IsAbs(sourcePath) {
            resolvedPath, err := paths.ResolveRelative(projectRoot, sourcePath)
            if err != nil {
                return nil, fmt.Errorf("failed to resolve path: %w", err)
            }
            sourcePath = resolvedPath
        }

        // Validate path is within project (if relative)
        if !filepath.IsAbs(source.Path) {
            if err := paths.ValidateLocalPath(projectRoot, sourcePath); err != nil {
                return nil, err
            }
        }

        // ... rest of local path handling ...
    }
}
```

## Migration Strategy

### Phase 1: Add Support for Relative Paths (Non-Breaking)

1. Implement `internal/core/paths` package
2. Update add/update commands to write relative paths
3. Update source fetching to handle relative paths
4. Add tests

**Result:** New configs use relative paths, old configs still work

### Phase 2: Warn About Absolute Paths

1. Add deprecation warnings when loading absolute paths
2. Update `dnaspec validate` to **warn** (not error) about absolute paths
3. Early warning in `dnaspec add` before parsing manifest
4. Update documentation to recommend relative paths

**Result:** Users aware of deprecation, can migrate at their own pace

### Phase 3: Auto-Migration (Optional)

1. Add migration helper to `dnaspec update`
2. Auto-convert on update operations
3. Document migration path

**Result:** Easy migration for existing users

## Edge Cases

### Edge Case 1: Path Outside Project

**Scenario:** User tries to add `/etc/dna` from `/Users/me/project`

**Handling:**
```go
relPath, err := paths.MakeRelative(projectRoot, absPath)
// err: "path is outside project root"

// Show warning, require confirmation
// Store as absolute with warning in config
```

### Edge Case 2: Symlink to Outside Project

**Scenario:** `project/dna` → symlink to `/usr/share/dna`

**Handling:**
```go
// ResolveRelative follows symlinks
realPath, _ := filepath.EvalSymlinks(absPath)
// realPath: /usr/share/dna

// Validate against project root
if !isWithinProject(realPath, projectRoot) {
    return error("symlink points outside project")
}
```

### Edge Case 3: Relative Path with `..`

**Scenario:** Path is `../sibling-project/dna`

**Handling:**
```go
absPath := filepath.Join(projectRoot, "../sibling-project/dna")
relPath, err := paths.MakeRelative(projectRoot, absPath)
// err: "path is outside project root"

// Reject or warn
```

### Edge Case 4: Current Directory Reference

**Scenario:** User specifies `./dna` or `dna`

**Handling:**
```go
// Normalize to remove ./ prefix
relPath = strings.TrimPrefix(relPath, "./")
// Store: dna
```

### Edge Case 5: Windows Paths

**Scenario:** `C:\Users\me\project\dna`

**Handling:**
```go
// filepath package handles cross-platform paths
// No special handling needed
```

## Testing Strategy

### Unit Tests

**Package:** `internal/core/paths`

```go
func TestMakeRelative(t *testing.T) {
    tests := []struct {
        name        string
        projectRoot string
        absPath     string
        want        string
        wantErr     bool
    }{
        {
            name:        "simple subdirectory",
            projectRoot: "/Users/me/project",
            absPath:     "/Users/me/project/dna",
            want:        "dna",
        },
        {
            name:        "nested subdirectory",
            projectRoot: "/Users/me/project",
            absPath:     "/Users/me/project/shared/dna",
            want:        "shared/dna",
        },
        {
            name:        "outside project",
            projectRoot: "/Users/me/project",
            absPath:     "/Users/me/other/dna",
            wantErr:     true,
        },
        {
            name:        "parent directory",
            projectRoot: "/Users/me/project",
            absPath:     "/Users/me/dna",
            wantErr:     true,
        },
    }
    // ...
}
```

### Integration Tests

**File:** `internal/cli/project/add_test.go`

```go
func TestAddLocalSource_RelativePaths(t *testing.T) {
    // Test that add command creates relative paths
}

func TestAddLocalSource_OutsideProject(t *testing.T) {
    // Test handling of paths outside project
}
```

## Documentation Updates

### Files to Update

1. **docs/design.md**
   - Update configuration examples to use relative paths
   - Update local-path source description
   - Add note about absolute path deprecation

2. **examples/project-config-local-only.yaml**
   - Change: `path: "/Users/me/my-dna-patterns"`
   - To: `path: "local-dna"` (with accompanying directory structure)

3. **examples/project-config-multi-source.yaml**
   - Update local source example to relative path

4. **README.md** (if exists)
   - Update examples

### New Documentation

Add section to design.md:

```markdown
## Local Path Handling

**Relative Paths (Recommended):**
DNASpec stores local source paths as relative to the project root:

```yaml
sources:
  - name: shared-dna
    type: local-path
    path: shared/dna  # Relative to project root
```

**Benefits:**
- Works on any machine
- Safe for version control
- Prevents accidental system directory references

**Absolute Paths (Deprecated):**
Absolute paths still work but will show deprecation warnings:
```yaml
  - name: legacy
    type: local-path
    path: /Users/me/dna  # Deprecated!
```

Run `dnaspec update legacy` to convert to relative path.
```

## Security Considerations

### Path Traversal Prevention

**Issue:** Malicious relative paths could try to escape project

**Mitigation:**
```go
// Validate after joining with project root
absPath := filepath.Join(projectRoot, relPath)
cleanPath := filepath.Clean(absPath)

// Ensure result is within project
if !strings.HasPrefix(cleanPath, projectRoot+string(filepath.Separator)) {
    return error("path escapes project root")
}
```

### Symlink Security

**Issue:** Symlinks could point outside project

**Mitigation:**
```go
// Always resolve symlinks before validation
realPath, err := filepath.EvalSymlinks(path)
if err != nil {
    return err
}

// Validate resolved path
if !isWithinProject(realPath, projectRoot) {
    return error("symlink target outside project")
}
```

## Performance Considerations

**Impact:** Minimal
- Path resolution: O(1) operation
- Symlink resolution: Filesystem operation (cached by OS)
- No network operations involved

**Optimization:**
- Cache project root discovery
- Reuse path resolution results where possible

## Rollback Plan

If issues arise:

1. **Immediate**: Add flag `--allow-absolute` to bypass validation
2. **Short-term**: Add config option `allow_absolute_paths: true`
3. **Long-term**: Fix underlying issue and remove workaround

Config would never be in a bad state because:
- Old absolute paths continue to work (with warnings)
- New relative paths are validated before save
- Atomic writes prevent partial updates
