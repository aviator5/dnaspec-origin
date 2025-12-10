# Refine Update Command: Interactive Guideline Selection

## Problem Statement

The current `dnaspec update` command supports updating all sources at once via the `--all` flag, but this creates a coarse-grained update experience that doesn't give users fine control over which guidelines to add, update, or keep when refreshing a source. Specifically:

1. **No per-guideline selection**: Users cannot selectively choose which guidelines from a source manifest to add or update
2. **Batch updates with --all**: The `--all` flag updates every source without interaction, making it hard to review changes per source
3. **Missing guidelines not highlighted**: When guidelines exist in the project config but are missing from the source manifest, users have no clear visual indication
4. **Binary add-new policy**: The `--add-new=all|none` flag forces an all-or-nothing decision for new guidelines

## Proposed Solution

Remove support for `--all` flag and mandate one-source-at-a-time updates with interactive guideline selection. For each source update, the command will:

1. **Parse source manifest** and display all available guidelines
2. **Pre-check existing guidelines** that are already in the project configuration
3. **Highlight orphaned guidelines** (present in config but missing in source) with a warning icon
4. **Allow interactive selection** via a multi-select UI using the existing `charmbracelet/huh` library

This gives users complete control over which guidelines to keep, add, or remove while providing clear visual feedback about the state of each guideline.

## Scope

### In Scope

- Remove `--all` flag and `updateAllSources` function from update command
- Remove `--add-new` flag (replaced by interactive selection)
- Add interactive multi-select UI for guideline selection
- Parse source manifest and build available guidelines list
- Pre-check guidelines already in project configuration
- Highlight orphaned guidelines (in config but not in source) with warning indicator
- Update argument validation to require source name

### Out of Scope

- Batch updates of multiple sources (removed capability)
- Changes to `--dry-run` flag behavior (remains unchanged for single-source preview)
- Changes to prompt handling or file copying logic
- Changes to git source fetching or local source loading

## User Impact

**Breaking Change**: Users relying on `dnaspec update --all` will need to update each source individually:

```bash
# Before
dnaspec update --all

# After
dnaspec update source1
dnaspec update source2
```

**Migration Path**: Scripts using `--all` can iterate over sources:

```bash
for source in $(yq '.sources[].name' dnaspec.yaml); do
  dnaspec update "$source"
done
```

## Dependencies

- Existing `internal/ui/selection.go` multi-select functionality (already uses `huh`)
- Existing manifest parsing and guideline comparison logic
- Existing source fetching and file copying infrastructure

## Alternatives Considered

1. **Keep --all but add interactive mode**: This would allow batch updates but prompt for each source. Rejected because it creates a confusing UX (too many prompts) and doesn't address the core issue of coarse-grained control.

2. **Add --select flag alongside --all**: This would make selection opt-in. Rejected because it increases flag complexity and leaves the default behavior unchanged.

3. **Use separate command for interactive selection**: Create `dnaspec update-select` or similar. Rejected because it fragments the update workflow unnecessarily.
