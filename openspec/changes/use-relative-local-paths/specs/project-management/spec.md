# project-management Spec Delta

**Change ID:** `use-relative-local-paths`

## MODIFIED Requirements

### Requirement: Add DNA Source from Local Directory

The system SHALL provide ability to add DNA guidelines from a local directory using `dnaspec add <local-path>` command with **relative path storage**.

#### Scenario: Add source from local directory successfully

- **WHEN** user runs `dnaspec add <local-path>` with valid directory path
- **THEN** read `dnaspec-manifest.yaml` from directory
- **AND** derive source name from directory name
- **AND** **convert absolute input path to relative path from project root**
- **AND** **validate that relative path does not escape project directory**
- **AND** display interactive guideline selection
- **AND** copy selected guideline and prompt files to `dnaspec/<source-name>/` directory
- **AND** update `dnaspec.yaml` with source metadata (type, **relative path**, guidelines, prompts)
- **AND** display next steps

#### Scenario: Local directory outside project root

- **WHEN** user runs `dnaspec add <local-path>` where path is outside the project directory
- **THEN** display warning about path portability
- **AND** explain that absolute paths won't work on other machines
- **AND** prompt user for confirmation to continue
- **AND** if confirmed, store absolute path with deprecation note
- **AND** if cancelled, abort without modifying configuration

#### Scenario: Relative path input

- **WHEN** user runs `dnaspec add ./local-dna` or `dnaspec add local-dna`
- **THEN** resolve to absolute path for validation
- **AND** convert to relative path from project root for storage
- **AND** normalize by removing `./ ` prefix

### Requirement: Update Source from Origin

The system SHALL provide `dnaspec update <source-name>` command that fetches latest content from git repository or local directory.

#### Scenario: Update local source with absolute path

- **WHEN** updating source with type "local-path" that has absolute path
- **THEN** **attempt to convert absolute path to relative path**
- **AND** if conversion succeeds, update configuration with relative path
- **AND** display message "✓ Converted to relative path: <rel-path>"
- **AND** if conversion fails (outside project), keep absolute path with warning

### Requirement: Configuration Validation

The system SHALL provide `dnaspec validate` command that checks project configuration for errors and warnings.

#### Scenario: Validate configuration with absolute local paths

- **WHEN** configuration contains local-path source with absolute path
- **THEN** add error "Source '<name>' must use relative path, found absolute: <path>"
- **AND** suggest converting to relative path or running `dnaspec update <name>`
- **AND** include in errors count
- **AND** validation fails

#### Scenario: Validate local path within project

- **WHEN** configuration contains local-path source with relative path
- **THEN** resolve path relative to project root
- **AND** verify resolved path exists
- **AND** verify resolved path (after following symlinks) is within project directory
- **AND** if path escapes project, add error "Source '<name>' path escapes project directory"

#### Scenario: Validate symlink to outside project

- **WHEN** local-path source is symlink pointing outside project directory
- **THEN** resolve symlink to actual target
- **AND** validate target location is within project
- **AND** if target outside project, add error "Source '<name>' symlink points outside project directory"
- **AND** validation fails

## ADDED Requirements

### Requirement: Relative Path Storage for Local Sources

The system SHALL store local source paths as relative to the project root directory.

#### Scenario: Store relative path in configuration

- **WHEN** adding or updating local source within project directory
- **THEN** calculate relative path from project root to source directory
- **AND** store relative path (without `./ ` prefix) in `dnaspec.yaml`
- **AND** ensure path portability across different machines

#### Scenario: Resolve relative path when reading

- **WHEN** loading configuration with local-path source using relative path
- **THEN** resolve relative path from project root directory
- **AND** convert to absolute path for file operations
- **AND** validate resolved path is within project directory

### Requirement: Path Validation for Local Sources

The system SHALL validate that local source paths remain within the project directory.

#### Scenario: Prevent directory traversal in local paths

- **WHEN** resolving local source path (relative or absolute)
- **THEN** clean and normalize the path
- **AND** resolve symlinks to actual target locations
- **AND** verify final resolved path is within or equals project root
- **AND** reject paths that escape project directory via `..` or symlinks

#### Scenario: Allow paths within project subdirectories

- **WHEN** local source path resolves to subdirectory of project
- **THEN** accept path as valid
- **AND** allow paths like `shared/dna`, `packages/common/guidelines`, etc.

### Requirement: Backward Compatibility for Absolute Paths

The system SHALL maintain backward compatibility with existing configurations using absolute paths.

#### Scenario: Load configuration with absolute path

- **WHEN** loading `dnaspec.yaml` containing local-path source with absolute path
- **THEN** parse configuration successfully
- **AND** display deprecation warning to stderr
- **AND** include source name and absolute path in warning
- **AND** suggest running `dnaspec update <name>` to migrate
- **AND** continue normal operation with absolute path

#### Scenario: Use absolute path for file operations

- **WHEN** local-path source has absolute path (legacy configuration)
- **THEN** use absolute path directly for file operations
- **AND** skip relative path resolution
- **AND** warn on each operation about deprecation

### Requirement: Path Conversion During Migration

The system SHALL automatically convert absolute paths to relative paths when possible.

#### Scenario: Auto-convert during update command

- **WHEN** running `dnaspec update <source>` on local source with absolute path
- **THEN** check if absolute path is within project directory
- **AND** if yes, calculate relative path and update configuration
- **AND** display "✓ Converted to relative path: <rel-path>" message
- **AND** if no, keep absolute path and warn about portability

#### Scenario: Manual conversion via config edit

- **WHEN** user manually edits `dnaspec.yaml` to change absolute to relative path
- **THEN** on next command, resolve and validate new relative path
- **AND** if invalid, display clear error message
- **AND** if valid, continue normal operation

## MODIFIED Requirements (Documentation Impact)

### Requirement: List Project Sources and Guidelines

The system SHALL provide `dnaspec list` command that displays configured sources, guidelines, prompts, and agents.

#### Scenario: Display local source with relative path

- **WHEN** listing sources that include local-path sources with relative paths
- **THEN** display source name, type "local-path", and relative path
- **AND** format: "  Path: <relative-path> (relative to project root)"
- **AND** distinguish from absolute paths in display

#### Scenario: Display local source with absolute path

- **WHEN** listing sources that include local-path sources with absolute paths
- **THEN** display source name, type "local-path", and absolute path
- **AND** add "[deprecated]" marker next to absolute paths
- **AND** format: "  Path: <absolute-path> [deprecated]"

## Notes

### Related Specifications

This change affects the `project-management` capability only. No changes required to:
- `manifest-management` - Manifest paths are already relative within DNA sources
- `agent-integration` - Agent file generation is unaffected

### Migration Path for Users

Users with existing absolute paths can migrate by:
1. Running `dnaspec update <source-name>` to auto-convert
2. Manually editing `dnaspec.yaml` to change paths to relative
3. Running `dnaspec validate` to verify the change

### Security Considerations

Path validation prevents:
- Directory traversal attacks via `../../etc/passwd` paths
- Accidental reference to system directories
- Symlink attacks pointing outside project

All validation occurs after resolving symlinks to prevent circumvention.

### Cross-Platform Compatibility

Implementation uses `filepath` package which handles:
- Windows vs Unix path separators
- Case sensitivity differences
- Path normalization across platforms

No platform-specific code required.
