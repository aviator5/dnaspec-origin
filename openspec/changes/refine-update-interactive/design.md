# Design: Interactive Guideline Selection for Update Command

## Overview

This design removes the `--all` flag from `dnaspec update` and replaces the binary `--add-new` policy with an interactive multi-select UI that gives users complete control over guideline selection.

## Architecture

### Key Components

```
┌─────────────────────────────────────────────────────────────┐
│ dnaspec update <source-name>                                │
└───────────────────┬─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│ 1. Fetch source (git or local)                              │
│    - Load manifest from source                              │
│    - Compare with current project config                    │
└───────────────────┬─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Build guideline selection list                           │
│    - All guidelines from source manifest                    │
│    - Mark already-added guidelines as pre-checked           │
│    - Mark orphaned guidelines (in config, not in source)    │
└───────────────────┬─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Display interactive multi-select                         │
│    - Use charmbracelet/huh library                          │
│    - Format: "[name] - [description] [⚠️ if orphaned]"      │
│    - Pre-select existing guidelines                         │
└───────────────────┬─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Apply selection                                          │
│    - Update project config with selected guidelines         │
│    - Copy guideline and prompt files                        │
│    - Update commit hash (for git sources)                   │
└─────────────────────────────────────────────────────────────┘
```

### Guideline States

Each guideline can be in one of three states:

1. **Available**: Present in source manifest, not in current config
   - Display normally in selection list
   - Not pre-checked

2. **Existing**: Present in both source manifest and current config
   - Display normally in selection list
   - Pre-checked by default

3. **Orphaned**: Present in current config, missing from source manifest
   - Display with ⚠️ warning icon
   - Include in selection list at the end
   - Pre-checked by default (allow user to remove)

## Data Structures

### Selection Item

```go
type GuidelineSelectionItem struct {
    Name        string
    Description string
    IsOrphaned  bool
    IsExisting  bool
}
```

### Selection Result

The multi-select UI returns a list of selected guideline names. The system then:
1. Filters to guidelines that exist in the source manifest (non-orphaned)
2. Updates the project config with only the selected guidelines
3. Removes any unselected guidelines (including orphaned ones if unchecked)

## UI Flow

### Interactive Selection Display

```
Select guidelines to keep or add:
┌─────────────────────────────────────────────────────────────┐
│ Use space to select/deselect, enter to confirm              │
│                                                              │
│ [✓] go-style - Go coding style guidelines                   │
│ [✓] rest-api - REST API design patterns                     │
│ [ ] security-review - Security review checklist             │
│ [✓] old-guideline - Deprecated guideline ⚠️                 │
└─────────────────────────────────────────────────────────────┘
```

### Dry Run Mode

Dry run remains unchanged and shows a preview without the interactive selection:

```bash
$ dnaspec update my-source --dry-run

⏳ Fetching latest from https://github.com/...
✓ Current commit: abc123
✓ Latest commit: def456 (changed)

Available guidelines:
  - go-style: Go coding style guidelines
  - rest-api: REST API design patterns
  - security-review: Security review checklist

Already in config:
  ✓ go-style
  ✓ rest-api

Orphaned (in config but not in source):
  ⚠ old-guideline

=== Dry Run - Preview ===
No changes made (dry run)
```

## Implementation Notes

### Reuse Existing UI Components

The codebase already has multi-select functionality in `internal/ui/selection.go`:

```go
func SelectGuidelines(available []config.ManifestGuideline) ([]config.ManifestGuideline, error)
```

We'll extend this with a new function:

```go
func SelectGuidelinesWithStatus(
    available []config.ManifestGuideline,
    existing []config.ProjectGuideline,
    orphaned []config.ProjectGuideline,
) ([]config.ManifestGuideline, error)
```

### Validation

- **Before**: Validate that either source-name OR --all is provided
- **After**: Validate that source-name is provided (--all removed)

### File Operations

No changes to file copying logic. The same `files.CopyGuidelineFiles` function will be used, just with a different set of guidelines based on user selection.

## Error Handling

1. **No source name provided**: Exit with error "must specify a source name"
2. **Source not found**: Display available sources and exit
3. **Selection canceled**: Exit gracefully without changes
4. **Empty selection**: Allow (user can remove all guidelines from a source)

## Testing Strategy

### Unit Tests

1. Test guideline categorization (available, existing, orphaned)
2. Test selection item formatting with orphaned marker
3. Test validation rejects missing source name
4. Test validation rejects --all flag as invalid

### Integration Tests

1. Test full update flow with interactive selection (mock UI)
2. Test dry run displays correct preview
3. Test orphaned guidelines are preserved in selection
4. Test empty selection removes all guidelines from source

### Manual Testing

1. Update source with new guidelines and verify interactive selection
2. Update source where some guidelines were removed and verify orphaned display
3. Verify pre-selection of existing guidelines
4. Verify deselecting existing guidelines removes them from config
