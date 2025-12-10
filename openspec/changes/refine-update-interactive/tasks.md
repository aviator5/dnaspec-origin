# Tasks: Refine Update Command Interactive Selection

## Phase 1: Remove --all Support

- [ ] Remove `--all` flag from `NewUpdateCmd`
- [ ] Remove `updateAllSources` function from `update.go`
- [ ] Update argument validation to require source name (remove --all checks)
- [ ] Update command help text and examples (remove --all references)
- [ ] Remove tests for `--all` flag behavior
- [ ] Update README.md to remove --all documentation

## Phase 2: Add Interactive Guideline Selection

- [ ] Add `SelectGuidelinesWithStatus` function to `internal/ui/selection.go`
  - Accept available (from manifest), existing (from config), and orphaned (in config but not in source) guidelines
  - Build options list with all three categories
  - Pre-select existing and orphaned guidelines
  - Add warning icon (⚠️) to orphaned guideline labels
  - Return selected guideline names

- [ ] Update `updateSingleSource` in `update.go`
  - Parse source manifest to get all available guidelines
  - Compare with current config to categorize (available, existing, orphaned)
  - Call `SelectGuidelinesWithStatus` with categorized guidelines
  - Apply selection and update config

- [ ] Remove `--add-new` flag and related logic
  - Remove flag definition
  - Remove `handleNewGuidelines` function
  - Remove addNewAll/addNewNone constants

## Phase 3: Update Tests

- [ ] Add unit tests for `SelectGuidelinesWithStatus`
  - Test option formatting with orphaned marker
  - Test pre-selection of existing and orphaned guidelines

- [ ] Update `update_test.go`
  - Remove tests for --all flag
  - Remove tests for --add-new flag
  - Add test for required source name validation
  - Add integration test for selection flow (with mocked UI)

- [ ] Test orphaned guideline handling
  - Test guideline categorization logic
  - Test orphaned guidelines appear in selection
  - Test deselecting orphaned guidelines removes them from config

## Phase 4: Update Documentation

### README.md

- [ ] Update `README.md`
  - Remove --all flag references
  - Remove --add-new flag references
  - Add interactive selection documentation
  - Add migration guide for scripts using --all

### docs/ Directory

- [ ] Update `docs/project-guide.md`
  - Remove --all flag examples (lines 351, 354, 624, 906)
  - Remove --add-new flag examples (lines 342, 345)
  - Update "dnaspec update" section (lines 329-405) with interactive selection flow
  - Remove batch update references (line 574)
  - Update comparison table (line 624) to remove --all row
  - Add interactive selection UI example with orphaned guideline display
  - Add migration guide for --all users

- [ ] Update `docs/design.md`
  - Remove --all flag examples and documentation (lines 539, 598)
  - Remove --add-new flag examples (lines 604, 607)
  - Update "dnaspec update" command section (lines 529-631) with new interactive flow
  - Update flow diagram to show interactive selection step
  - Remove references to batch updates

- [ ] Update command help text in code
  - Remove --all examples
  - Add interactive selection description
  - Update usage example

## Phase 5: Validation

- [ ] Run `make test` and ensure all tests pass
- [ ] Run `make lint` and fix any issues
- [ ] Manually test update command with:
  - Source with new guidelines
  - Source with removed guidelines (orphaned)
  - Source with no changes
  - --dry-run mode
