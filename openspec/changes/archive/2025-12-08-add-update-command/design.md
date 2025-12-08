# Design: Add Update Command

## Architecture Overview

The `dnaspec update` command will be implemented in the CLI layer (`internal/cli/project/`) and will reuse existing core functionality from `internal/core/source/`, `internal/core/config/`, and `internal/core/files/`.

```
┌─────────────────────────────────────────┐
│         CLI Layer                       │
│  internal/cli/project/update.go         │
│  - Command definition (cobra)           │
│  - Flag parsing                          │
│  - User interaction                      │
│  - Progress reporting                    │
└────────────┬────────────────────────────┘
             │
             │ calls
             ▼
┌─────────────────────────────────────────┐
│      Core Domain Layer                  │
│  internal/core/source/                  │
│  - FetchGitSource()                     │
│  - FetchLocalSource()                   │
│                                         │
│  internal/core/config/                  │
│  - FindSourceByName()                   │
│  - CompareGuidelines()                  │
│  - UpdateSourceInConfig()               │
│                                         │
│  internal/core/files/                   │
│  - CopyGuidelineFiles()                 │
└─────────────────────────────────────────┘
```

## Command Interface

### Cobra Command Structure

```go
// NewUpdateCmd creates the update command
func NewUpdateCmd() *cobra.Command {
    var flags updateFlags

    cmd := &cobra.Command{
        Use:   "update [source-name]",
        Short: "Update source(s) from their origin",
        Args:  cobra.MaximumNArgs(1),
        RunE:  func(cmd *cobra.Command, args []string) error {
            return runUpdate(flags, args)
        },
    }

    cmd.Flags().BoolVar(&flags.all, "all", false, "Update all sources")
    cmd.Flags().BoolVar(&flags.dryRun, "dry-run", false, "Preview changes")
    cmd.Flags().StringVar(&flags.addNew, "add-new", "", "Policy for new guidelines (all|none)")

    return cmd
}
```

### Flag Validation

The command must validate:
1. Either `source-name` arg OR `--all` flag (not both, not neither)
2. If `--add-new` provided, value must be "all" or "none"
3. Project config must exist

## Core Update Algorithm

### Single Source Update Flow

```
1. Load project config
   ↓
2. Find source by name
   ↓
3. Fetch latest from origin
   ↓
4. Parse latest manifest
   ↓
5. Compare current vs latest
   ↓
6. Update existing guidelines
   ↓
7. Handle new guidelines
   ↓
8. Report removed guidelines
   ↓
9. Update config metadata
   ↓
10. Save config
```

### Comparison Logic

When comparing current config with latest manifest, we need to categorize guidelines:

```go
type GuidelineComparison struct {
    Updated []string  // Guidelines that exist in both but have changes
    New     []string  // Guidelines in manifest but not in config
    Removed []string  // Guidelines in config but not in manifest
    Unchanged []string // Guidelines with no changes
}

func CompareGuidelines(
    currentGuidelines []ProjectGuideline,
    manifestGuidelines []ManifestGuideline,
) GuidelineComparison {
    // Build maps for O(1) lookup
    currentMap := make(map[string]ProjectGuideline)
    manifestMap := make(map[string]ManifestGuideline)

    for _, g := range currentGuidelines {
        currentMap[g.Name] = g
    }
    for _, g := range manifestGuidelines {
        manifestMap[g.Name] = g
    }

    var result GuidelineComparison

    // Check each current guideline
    for name, current := range currentMap {
        if manifest, exists := manifestMap[name]; exists {
            if hasChanges(current, manifest) {
                result.Updated = append(result.Updated, name)
            } else {
                result.Unchanged = append(result.Unchanged, name)
            }
        } else {
            result.Removed = append(result.Removed, name)
        }
    }

    // Find new guidelines
    for name := range manifestMap {
        if _, exists := currentMap[name]; !exists {
            result.New = append(result.New, name)
        }
    }

    return result
}

func hasChanges(current ProjectGuideline, manifest ManifestGuideline) bool {
    if current.Description != manifest.Description {
        return true
    }
    if !slicesEqual(current.ApplicableScenarios, manifest.ApplicableScenarios) {
        return true
    }
    if !slicesEqual(current.Prompts, manifest.Prompts) {
        return true
    }
    // Note: File content changes will always be copied regardless
    return false
}
```

## Update Process Details

### Git Source Update

```go
func updateGitSource(source *ProjectSource) error {
    // 1. Clone at configured ref
    sourceInfo, cleanup, err := source.FetchGitSource(source.URL, source.Ref)
    if err != nil {
        return err
    }
    defer cleanup()

    // 2. Check if commit changed
    if sourceInfo.Commit == source.Commit {
        fmt.Println("Source is up to date")
        return nil
    }

    fmt.Printf("Updating from %s to %s\n",
        source.Commit[:8], sourceInfo.Commit[:8])

    // 3. Continue with update process...
    return updateSourceFiles(source, sourceInfo)
}
```

### Local Source Update

```go
func updateLocalSource(source *ProjectSource) error {
    // 1. Read from configured path
    sourceInfo, err := source.FetchLocalSource(source.Path)
    if err != nil {
        return err
    }

    // 2. No commit hash to check for local sources
    fmt.Println("Refreshing from local directory...")

    // 3. Continue with update process
    return updateSourceFiles(source, sourceInfo)
}
```

### File Update Strategy

For updated and unchanged guidelines:
```
1. Find guideline in manifest
2. Update metadata in config (description, scenarios, prompts)
3. Copy files from source to dnaspec/<source-name>/
4. Overwrite existing files (no merge, no drift detection)
```

For removed guidelines:
```
1. Warn user that guideline was removed from source
2. Keep files and config (don't auto-delete)
3. User can manually remove if desired
```

## New Guidelines Handling

### Interactive Mode (Default)

```go
if len(comparison.New) > 0 && addNewPolicy == "" {
    fmt.Println("\nNew guidelines available:")
    for _, name := range comparison.New {
        guideline := findGuideline(manifest, name)
        fmt.Printf("  - %s: %s\n", name, guideline.Description)
    }

    answer := promptYesNo("Add new guidelines?")
    if answer {
        addNewPolicy = "all"
    } else {
        addNewPolicy = "none"
    }
}
```

### Non-Interactive Modes

```go
switch addNewPolicy {
case "all":
    for _, name := range comparison.New {
        guideline := findGuideline(manifest, name)
        source.Guidelines = append(source.Guidelines, convertToProjectGuideline(guideline))
        fmt.Printf("✓ Added %s\n", name)
    }
case "none":
    fmt.Printf("ℹ Skipped %d new guidelines\n", len(comparison.New))
default:
    // No new guidelines or user declined
}
```

## Configuration Update

### Atomic Config Update

```go
func UpdateSourceInConfig(cfg *ProjectConfig, sourceName string, updatedSource ProjectSource) error {
    // Find and replace source
    for i, src := range cfg.Sources {
        if src.Name == sourceName {
            cfg.Sources[i] = updatedSource
            return nil
        }
    }
    return fmt.Errorf("source '%s' not found", sourceName)
}

// Then use atomic write
err := config.AtomicWriteProjectConfig(projectConfigFileName, cfg)
```

### Metadata to Update

For git sources:
- `commit`: Update to new commit hash
- `guidelines`: Update metadata (description, scenarios, prompts)
- All other fields remain unchanged (url, ref, name, type)

For local sources:
- `guidelines`: Update metadata
- All other fields remain unchanged

## Update All Sources

```go
func updateAllSources(cfg *ProjectConfig, flags updateFlags) error {
    if len(cfg.Sources) == 0 {
        fmt.Println("No sources configured")
        return nil
    }

    fmt.Printf("Updating %d sources...\n\n", len(cfg.Sources))

    var errors []error
    for _, source := range cfg.Sources {
        fmt.Printf("=== Updating %s ===\n", source.Name)

        if err := updateSingleSource(cfg, source.Name, flags); err != nil {
            errors = append(errors, fmt.Errorf("%s: %w", source.Name, err))
            fmt.Println(ui.ErrorStyle.Render("✗ Failed"), err)
        }

        fmt.Println()
    }

    if len(errors) > 0 {
        return fmt.Errorf("failed to update %d sources", len(errors))
    }

    fmt.Println(ui.SuccessStyle.Render("✓ All sources updated"))
    return nil
}
```

## Dry Run Mode

Dry run should:
1. Perform all read operations (fetch source, parse manifest, compare)
2. Display what would change
3. Skip all write operations (file copying, config updates)

```go
if flags.dryRun {
    fmt.Println(ui.InfoStyle.Render("\n=== Dry Run - Preview ==="))
    fmt.Println("Would update:", len(comparison.Updated), "guidelines")
    fmt.Println("Would add:", len(comparison.New), "guidelines")
    fmt.Println("Removed from source:", len(comparison.Removed), "guidelines")
    fmt.Println("\nNo changes made (dry run)")
    return nil
}
```

## Error Handling

### Source Not Found
```go
source := findSourceByName(cfg, sourceName)
if source == nil {
    fmt.Println(ui.ErrorStyle.Render("✗ Error:"), "Source not found:", sourceName)
    fmt.Println("\nAvailable sources:")
    for _, src := range cfg.Sources {
        fmt.Printf("  - %s\n", src.Name)
    }
    return fmt.Errorf("source not found")
}
```

### Git Clone Failure
```go
sourceInfo, cleanup, err := source.FetchGitSource(source.URL, source.Ref)
if err != nil {
    return fmt.Errorf("failed to fetch git source: %w", err)
}
defer cleanup()
```

### Manifest Validation Failure
```go
if err := manifest.Validate(); err != nil {
    return fmt.Errorf("invalid manifest: %w", err)
}
```

## Output Format

### Successful Update
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

### No Changes
```
⏳ Fetching latest from https://github.com/company/dna...
✓ Current commit: abc123de
✓ Already at latest commit

All guidelines up to date.
```

## Testing Strategy

### Unit Tests
- `CompareGuidelines()` with various scenarios
- `UpdateSourceInConfig()` preserves other sources
- `hasChanges()` detects all types of changes

### Integration Tests
- Update git source with changes
- Update git source with no changes
- Update local source
- Update with new guidelines (interactive and non-interactive)
- Update all sources
- Dry run mode

### Test Data
Reuse test fixtures from `add` command tests:
- `testdata/valid-manifest/` for git sources
- Temp directories for local sources

## Performance Considerations

- Git clones use `--depth=1` (shallow clone)
- Reuse existing source fetching code (already optimized)
- Config updates are atomic but not cached (acceptable for update frequency)
- No performance concerns for typical use cases (1-10 sources)

## Security Considerations

All security measures from `add` command apply:
- Git URL validation (reject git:// protocol)
- Path traversal protection in manifest
- Atomic config writes prevent corruption
- Temporary directory cleanup prevents leaks

No new security concerns introduced by update command.
