# Implementation Tasks: Use Relative Paths for Local Sources

**Change ID:** `use-relative-local-paths`

## Task List

### 1. Create Path Utilities Package

**Description:** Implement core path handling functions in new `internal/core/paths` package

**Files:**
- Create: `internal/core/paths/paths.go`
- Create: `internal/core/paths/paths_test.go`

**Implementation:**
- [x] Create package directory and files
- [x] Implement `MakeRelative(projectRoot, absPath string) (string, error)`
  - Clean and normalize both paths
  - Resolve symlinks to get real paths
  - Calculate relative path from project root
  - Validate doesn't start with `..` (outside project)
  - Strip `./ ` prefix from result
- [x] Implement `ResolveRelative(projectRoot, relPath string) (string, error)`
  - Validate input is relative (not absolute)
  - Join with project root
  - Clean and resolve symlinks
  - Validate result is within project root
- [x] Implement `ValidateLocalPath(projectRoot, path string) error`
  - Handle both absolute and relative paths
  - Resolve to absolute path
  - Check within project boundaries
- [x] Implement `IsWithinProject(projectRoot, path string) (bool, error)`
  - Resolve all symlinks
  - Compare paths with proper prefix checking
- [x] Write comprehensive unit tests
  - Test simple subdirectories
  - Test nested subdirectories
  - Test paths outside project (should fail)
  - Test symlinks within project
  - Test symlinks pointing outside project
  - Test edge cases: `.`, `..`, `./foo`, etc.
  - Test cross-platform paths (use filepath package)

**Validation:**
- `go test ./internal/core/paths/... -v`
- All tests pass
- Coverage > 80%

**Dependencies:** None

---

### 2. Update Configuration Loading with Deprecation Warnings

**Description:** Add deprecation warnings when loading configs with absolute paths

**Files:**
- Modify: `internal/core/config/project.go`

**Implementation:**
- [x] Update `LoadProjectConfig` function
  - After parsing YAML, iterate sources
  - Check if `Type == SourceTypeLocalPath && filepath.IsAbs(Path)`
  - Print deprecation warning to stderr for each absolute path
  - Include source name, path, and migration suggestion
  - Continue loading (maintain backward compatibility)
- [x] Add `MigrateToRelativePaths(projectRoot string) error` method
  - Iterate all sources
  - For local-path sources with absolute paths
  - Call `paths.MakeRelative()` to convert
  - Update source.Path if successful
  - Return error if any conversion fails

**Validation:**
- Load test config with absolute path
- Verify warning printed to stderr
- Verify config still loads successfully

**Dependencies:** Task 1 (paths package)

---

### 3. Update Add Command to Use Relative Paths

**Description:** Modify `dnaspec add` to convert and store relative paths

**Files:**
- Modify: `internal/cli/project/add.go`

**Implementation:**
- [x] In local path handling section
  - Convert input path to absolute (for validation)
  - Get project root from config path
  - **EARLY CHECK**: Call `paths.MakeRelative(projectRoot, absPath)` BEFORE loading source
  - If successful: store relative path in source.Path
  - If fails (outside project):
    - Display warning about portability (with empty lines for visual clarity)
    - **Prompt user for confirmation BEFORE calling `source.FetchLocalSource()`**
    - If confirmed: store absolute path and proceed to load source
    - If cancelled: abort operation (saves time by not parsing manifest)
  - **THEN** proceed to load source with `source.FetchLocalSource()`
- [x] Update success message to indicate relative path
- [x] Handle relative path input (e.g., `./local-dna`)
  - Resolve to absolute first
  - Then convert back to relative from project root
  - Normalize by removing `./ ` prefix

**Validation:**
- Run `dnaspec add ./examples/dna` in test project
- Verify config stores `path: examples/dna`
- Run `dnaspec add /outside/path`
- Verify warning shown and confirmation required

**Dependencies:** Task 1 (paths package), Task 2 (config updates)

---

### 4. Update Source Fetching to Resolve Relative Paths

**Description:** Modify source fetching to handle relative paths correctly

**Files:**
- Modify: `internal/core/source/fetch.go`

**Implementation:**
- [x] In `FetchSource` function for local-path type
  - Check if source.Path is absolute or relative
  - If relative: call `paths.ResolveRelative(projectRoot, source.Path)`
  - If absolute: use as-is (backward compatibility)
  - Use resolved absolute path for subsequent operations
  - Add validation that resolved path is within project (for relative paths)
- [x] Update error messages to include path details

**Validation:**
- Create config with relative path
- Run command that fetches source
- Verify path resolves correctly
- Verify files accessed successfully

**Dependencies:** Task 1 (paths package)

---

### 5. Update Update Command with Auto-Migration

**Description:** Add automatic conversion of absolute to relative paths during update

**Files:**
- Modify: `internal/cli/project/update.go`

**Implementation:**
- [x] In update command logic
  - After fetching source, check if local-path with absolute path
  - Get project root from config path
  - Call `paths.MakeRelative(projectRoot, source.Path)`
  - If successful:
    - Update source.Path to relative path
    - Set flag indicating config needs saving
    - Print success message: "✓ Converted to relative path: <relPath>"
  - If fails (outside project):
    - Print warning about keeping absolute path
    - Suggest moving source into project
- [x] Ensure config saved if migration occurred

**Validation:**
- Create config with absolute path pointing to project subdirectory
- Run `dnaspec update <source>`
- Verify path converted to relative
- Verify config file updated
- Verify message printed

**Dependencies:** Task 1 (paths package), Task 4 (source fetching)

---

### 6. Update Validate Command with Path Checks

**Description:** Add validation warnings for absolute paths (warnings, not errors - backward compatible)

**Files:**
- Modify: `internal/cli/project/validate.go`

**Implementation:**
- [x] Add local path validation section
  - Create separate `warnings` slice (in addition to `errors`)
  - For each local-path source:
    - **If absolute path: add WARNING (not error)** - maintains backward compatibility!
    - Format: "Source '<name>' uses absolute path\n    Path: <path>\n    Run 'dnaspec update <name>' to convert to relative path"
    - Get project root from config path
    - Call `paths.ValidateLocalPath(projectRoot, source.Path)`
    - If validation fails: add ERROR with details
    - If path escapes project: add ERROR (this is a real problem)
- [x] Update output formatting to show warnings separately
  - Display warnings after errors
  - "⚠ Warning: Source '<name>' uses absolute path"
  - "  Path: <path>"
  - "  Run 'dnaspec update <name>' to convert to relative path"
  - Show count: "Found N warning(s). Run suggested commands to fix."
- [x] Warnings do NOT cause validation to fail (exit code 0 if only warnings)

**Validation:**
- Create config with absolute path
- Run `dnaspec validate`
- **Verify WARNING shown (not error)**
- Verify exit code is zero (validation passes with warnings)
- Create config with path escaping project
- Run `dnaspec validate`
- Verify ERROR shown and validation fails with non-zero exit code

**Dependencies:** Task 1 (paths package)

---

### 7. Update List Command Output

**Description:** Update `dnaspec list` to indicate path type

**Files:**
- Modify: `internal/cli/project/list.go`

**Implementation:**
- [x] In source display section for local-path sources
  - Check if path is absolute or relative
  - If relative: show "Path: <path> (relative to project root)"
  - If absolute: show "Path: <path> [deprecated]"
- [x] Update formatting to clearly distinguish types

**Validation:**
- Create config with both absolute and relative paths
- Run `dnaspec list`
- Verify correct labels shown

**Dependencies:** None (read-only display)

---

### 8. Update Documentation

**Description:** Update all documentation to reflect relative path usage

**Files:**
- Modify: `docs/design.md`
- Modify: `examples/project-config-local-only.yaml`
- Modify: `examples/project-config-multi-source.yaml`
- Modify: `dnaspec.yaml` (project's own config)

**Implementation:**
- [x] Update design.md
  - Change all local-path examples to use relative paths
  - Add section explaining relative vs absolute paths
  - Document deprecation of absolute paths
  - Add migration instructions
  - Update "Local Path Handling" section
- [x] Update example configs
  - Change `path: "/Users/me/..."` to `path: "local-dna"`
  - Add comments explaining relative paths
  - Show example directory structure
- [x] Update project's own dnaspec.yaml
  - Change current absolute path to relative
  - Verify it still works
- [x] Add migration guide
  - How to convert existing configs
  - What to do if source is outside project
  - Troubleshooting common issues

**Validation:**
- Review all documentation changes
- Verify examples are accurate
- Test example configs work as shown

**Dependencies:** All implementation tasks completed

---

### 9. Add Integration Tests

**Description:** Add end-to-end tests for relative path functionality

**Files:**
- Modify: `internal/cli/project/integration_test.go`
- Create: `internal/cli/project/add_relative_paths_test.go`

**Implementation:**
- [x] Test add command with relative path input
  - Create temp project directory
  - Create local DNA source in subdirectory
  - Run add command
  - Verify relative path stored in config
  - Verify path resolves correctly
- [x] Test add command with absolute path (inside project)
  - Run add with absolute path to project subdirectory
  - Verify converted to relative path
- [x] Test add command with absolute path (outside project)
  - Run add with path outside project
  - Verify warning shown
  - Verify confirmation required
- [x] Test update command migration
  - Create config with absolute path (inside project)
  - Run update command
  - Verify auto-converted to relative
- [x] Test validate command (warning-based validation)
  - Test with relative path (should pass with no warnings)
  - Test with absolute path (should **warn** but still pass - exit code 0)
  - Test with path outside project (should error and exit non-zero)
- [x] Test source fetching with relative paths
  - Create config with relative path
  - Run command that fetches source
  - Verify works correctly
- [x] Test backward compatibility
  - Create config with absolute path
  - Verify still works (with warnings)

**Validation:**
- `go test ./internal/cli/project/... -v -run TestRelative`
- All integration tests pass
- Cover main user journeys

**Dependencies:** All implementation tasks

---

### 10. Manual Testing and Edge Cases

**Description:** Manually test edge cases and real-world scenarios

**Test Cases:**
- [x] Test with symlink to directory inside project
  - Create symlink in project to subdirectory
  - Add via symlink path
  - Verify relative path stored
  - Verify resolves correctly
- [x] Test with symlink to directory outside project
  - Create symlink to outside directory
  - Add via symlink
  - Verify warning/error shown
- [x] Test with `./` and `../` in paths
  - Add `./local-dna`
  - Verify normalized to `local-dna`
  - Try `../outside` (should warn/error)
- [x] Test in monorepo structure
  - Create nested project structure
  - Add local sources at various levels
  - Verify relative paths work correctly
- [x] Test on different platforms
  - Verify works on macOS
  - Verify works on Linux
  - (Ideally) verify works on Windows
- [x] Test migration from absolute to relative
  - Create config with absolute paths
  - Run update command
  - Verify conversion successful
  - Verify all commands still work

**Validation:**
- Document results of each test
- Fix any issues found
- Add regression tests for bugs discovered

**Dependencies:** All previous tasks

---

## Task Ordering and Parallelization

**Phase 1: Foundation (Sequential)**
1. Create path utilities package (Task 1)

**Phase 2: Core Implementation (Can parallelize)**
2. Update config loading (Task 2)
3. Update add command (Task 3)
4. Update source fetching (Task 4)

**Phase 3: Additional Commands (Can parallelize)**
5. Update update command (Task 5)
6. Update validate command (Task 6)
7. Update list command (Task 7)

**Phase 4: Polish (Sequential)**
8. Update documentation (Task 8)
9. Add integration tests (Task 9)
10. Manual testing (Task 10)

## Estimated Effort

- **Task 1:** 2-3 hours (core functionality + comprehensive tests)
- **Task 2:** 30 minutes (simple addition)
- **Task 3:** 1 hour (user interaction + validation)
- **Task 4:** 30 minutes (straightforward resolution)
- **Task 5:** 30 minutes (similar to Task 3)
- **Task 6:** 45 minutes (validation logic)
- **Task 7:** 15 minutes (display only)
- **Task 8:** 1 hour (documentation updates)
- **Task 9:** 1-2 hours (comprehensive integration tests)
- **Task 10:** 1 hour (manual testing)

**Total:** ~8-10 hours

## Success Criteria

- [x] All unit tests pass
- [x] All integration tests pass
- [x] Manual testing scenarios complete
- [x] Documentation updated and accurate
- [x] Backward compatibility maintained
- [x] No breaking changes to existing configurations
- [x] Clear migration path documented
- [x] Example configs updated
- [x] `openspec validate use-relative-local-paths --strict` passes

## Risk Mitigation

**Risk:** Breaking existing configurations
- **Mitigation:** Maintain full backward compatibility, only warn about deprecation

**Risk:** Edge cases with symlinks
- **Mitigation:** Comprehensive testing, resolve symlinks before validation

**Risk:** Cross-platform path issues
- **Mitigation:** Use filepath package, test on multiple platforms

**Risk:** User confusion during migration
- **Mitigation:** Clear warnings, helpful error messages, good documentation
