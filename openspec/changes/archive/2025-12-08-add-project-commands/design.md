# Design: Project Commands (init and add)

## Context

Project developers need commands to initialize DNASpec in their projects and integrate DNA guidelines from various sources. This builds upon the manifest management foundation to enable the core DNASpec workflow.

**Constraints:**
- Go 1.25 language version
- Reuse existing manifest validation from `manifest-management` spec
- Support both git repositories and local directories as sources
- Security: prevent path traversal, validate git URLs, timeout git operations
- Must work cross-platform (Windows, macOS, Linux)

**Stakeholders:**
- Project developers integrating DNA guidelines into their projects
- DNA repository maintainers who need to test their manifests locally

## Goals / Non-Goals

**Goals:**
- Initialize empty project configuration with `dnaspec init`
- Add DNA sources from git repositories with `dnaspec add --git-repo <url>`
- Add DNA sources from local directories with `dnaspec add <path>`
- Interactive guideline selection for usability
- Non-interactive modes for CI/CD and scripting (`--all`, `--guideline` flags)
- Secure git operations (URL validation, timeouts, shallow clones)
- Namespace files under `dnaspec/<source-name>/` to prevent conflicts
- Atomic configuration updates to prevent corruption

**Non-Goals:**
- Symlink support for local directories (can be added later if needed)
- Git clone caching (optimization for future)
- Update command (separate change)
- Remove command (separate change)
- Agent file generation (separate change)
- Automatic guideline filtering based on project type

## Decisions

### Configuration File Structure

**Project Config (`dnaspec.yaml`):**
```yaml
version: 1

agents:
  - "claude-code"
  - "github-copilot"

sources:
  - name: "company-dna"
    type: "git-repo"
    url: "https://github.com/company/dna"
    ref: "v1.2.0"
    commit: "abc123def456789..."
    guidelines:
      - name: "go-style"
        file: "guidelines/go-style.md"
        description: "Go code style conventions"
        applicable_scenarios:
          - "writing new Go code"
        prompts: ["go-code-review"]
    prompts:
      - name: "go-code-review"
        file: "prompts/go-code-review.md"
        description: "Review Go code"
```

**Rationale:**
- Mirrors manifest structure for consistency
- Includes source metadata for updates and tracking
- Stores full guideline details (not just names) for offline use
- Commit hash enables detection of updates for git sources

### Source Name Derivation

**Algorithm:**
```go
func DeriveSourceName(gitURL, localPath string) string {
    var raw string
    if gitURL != "" {
        // Extract from URL: https://github.com/company/dna-guidelines.git → dna-guidelines
        raw = extractRepoName(gitURL)
    } else {
        // Extract from path: /Users/me/my-patterns → my-patterns
        raw = filepath.Base(localPath)
    }
    return SanitizeName(raw)
}

func SanitizeName(name string) string {
    // Convert to lowercase
    name = strings.ToLower(name)
    // Replace non-alphanumeric with hyphens
    name = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(name, "-")
    // Trim and collapse hyphens
    name = strings.Trim(name, "-")
    name = regexp.MustCompile(`-+`).ReplaceAllString(name, "-")
    return name
}
```

**Rationale:**
- Automatic naming reduces friction
- Consistent naming convention (lowercase with hyphens)
- Users can override with `--name` flag if needed
- Sanitization prevents filesystem issues

### Git Operations

**Library Choice: Command-line git vs go-git**

**Decision: Use command-line git via exec**

**Rationale:**
- Simpler implementation (no go-git dependency)
- More reliable (uses system git, which is battle-tested)
- Better SSH key handling (uses system SSH config)
- Easier authentication (HTTPS tokens, SSH keys work automatically)
- Smaller binary size

**Security Measures:**
```go
func ValidateGitURL(url string) error {
    // Only allow https:// and git@ (SSH)
    if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "git@") {
        return errors.New("only HTTPS and SSH URLs supported")
    }
    // Reject insecure git:// protocol
    if strings.HasPrefix(url, "git://") {
        return errors.New("git:// protocol not allowed (insecure)")
    }
    return nil
}

func CloneRepo(url, ref, tempDir string) (commit string, error) {
    // Validate URL first
    if err := ValidateGitURL(url); err != nil {
        return "", err
    }

    // Create timeout context (5 minutes)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    // Shallow clone for efficiency
    cmd := exec.CommandContext(ctx, "git", "clone", "--depth=1", "--single-branch", url, tempDir)
    if ref != "" {
        cmd.Args = append(cmd.Args, "--branch", ref)
    }

    // Run with timeout
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("git clone failed: %w", err)
    }

    // Get commit hash
    cmd = exec.CommandContext(ctx, "git", "-C", tempDir, "rev-parse", "HEAD")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

    return strings.TrimSpace(string(output)), nil
}
```

**Why this approach:**
- Defense in depth (URL validation + timeout)
- Shallow clone reduces bandwidth and disk usage
- Timeout prevents hanging on large repos or network issues
- Returns commit hash for tracking in config

### Temporary Directory Management

**Strategy:**
```go
func CreateTempCloneDir() (path string, cleanup func(), error) {
    // Generate unique ID
    randomBytes := make([]byte, 8)
    rand.Read(randomBytes)
    randomID := hex.EncodeToString(randomBytes)

    // Create unique temp directory
    tempDir := filepath.Join(
        os.TempDir(),
        "dnaspec",
        fmt.Sprintf("%d-%s", os.Getpid(), randomID),
    )

    if err := os.MkdirAll(tempDir, 0755); err != nil {
        return "", nil, err
    }

    // Cleanup function
    cleanup := func() {
        os.RemoveAll(tempDir)
    }

    return tempDir, cleanup, nil
}
```

**Usage pattern:**
```go
tempDir, cleanup, err := CreateTempCloneDir()
if err != nil {
    return err
}
defer cleanup()  // Always cleanup, even on error

// Use tempDir...
```

**Rationale:**
- Uses system temp directory (cross-platform)
- Unique PID + random ID prevents conflicts
- Cleanup function with defer ensures no leaks
- Works correctly even if process crashes (OS cleans temp dir eventually)

### Interactive Guideline Selection

**Library Choice: charmbracelet/huh**

Already in dependencies for terminal UI. Provides:
- Multi-select forms
- Good keyboard navigation
- Consistent with other terminal output
- Works in all terminal types

**Implementation:**
```go
func SelectGuidelines(available []ManifestGuideline) ([]ManifestGuideline, error) {
    if len(available) == 0 {
        return nil, errors.New("no guidelines available")
    }

    // Build options
    var options []huh.Option[string]
    for _, g := range available {
        label := fmt.Sprintf("%s - %s", g.Name, g.Description)
        options = append(options, huh.NewOption(label, g.Name))
    }

    // Create multi-select form
    var selected []string
    form := huh.NewForm(
        huh.NewGroup(
            huh.NewMultiSelect[string]().
                Title("Select guidelines to add:").
                Options(options...).
                Value(&selected),
        ),
    )

    if err := form.Run(); err != nil {
        return nil, err
    }

    // Filter to selected
    var result []ManifestGuideline
    for _, name := range selected {
        for _, g := range available {
            if g.Name == name {
                result = append(result, g)
                break
            }
        }
    }

    return result, nil
}
```

**Non-interactive modes:**
- `--all`: Select all guidelines without prompting
- `--guideline <name>`: Select specific guidelines by name (repeatable flag)
- Validate guideline names exist in manifest before proceeding

### File Operations

**File Copying Strategy:**
```go
func CopyGuidelineFiles(sourceDir, destDir string, guidelines []ManifestGuideline, prompts []ManifestPrompt) error {
    // Copy guidelines
    for _, g := range guidelines {
        src := filepath.Join(sourceDir, g.File)
        dst := filepath.Join(destDir, g.File)
        if err := copyFile(src, dst); err != nil {
            return fmt.Errorf("failed to copy %s: %w", g.File, err)
        }
    }

    // Copy prompts
    for _, p := range prompts {
        src := filepath.Join(sourceDir, p.File)
        dst := filepath.Join(destDir, p.File)
        if err := copyFile(src, dst); err != nil {
            return fmt.Errorf("failed to copy %s: %w", p.File, err)
        }
    }

    return nil
}

func copyFile(src, dst string) error {
    // Create destination directory
    if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
        return err
    }

    // Read source
    data, err := os.ReadFile(src)
    if err != nil {
        return err
    }

    // Write destination
    return os.WriteFile(dst, data, 0644)
}
```

**Rationale:**
- Preserves relative path structure from manifest
- Creates directories as needed
- Simple read/write approach (files are small markdown)
- Clear error messages with file paths

**Atomic Config Updates:**
```go
func AtomicWriteConfig(path string, config *ProjectConfig) error {
    // Marshal to YAML
    data, err := yaml.Marshal(config)
    if err != nil {
        return err
    }

    // Write to temp file
    tmpFile := path + ".tmp"
    if err := os.WriteFile(tmpFile, data, 0644); err != nil {
        return err
    }

    // Atomic rename
    return os.Rename(tmpFile, path)
}
```

**Why atomic writes:**
- Prevents corruption if process crashes during write
- Rename is atomic on all major filesystems
- Critical for configuration files

### Prompt Extraction

**Algorithm:**
```go
func ExtractReferencedPrompts(selectedGuidelines []ManifestGuideline, allPrompts []ManifestPrompt) []ManifestPrompt {
    // Build set of referenced prompt names
    referenced := make(map[string]bool)
    for _, g := range selectedGuidelines {
        for _, pName := range g.Prompts {
            referenced[pName] = true
        }
    }

    // Filter prompts to referenced ones
    var result []ManifestPrompt
    for _, p := range allPrompts {
        if referenced[p.Name] {
            result = append(result, p)
        }
    }

    return result
}
```

**Rationale:**
- Only copy prompts that are actually used by selected guidelines
- Reduces unnecessary files in project
- Maintains referential integrity

### Directory Structure

**Project directory after `dnaspec add`:**
```
my-project/
├── dnaspec.yaml
├── dnaspec/
│   ├── company-dna/           # Source name as namespace
│   │   ├── guidelines/
│   │   │   └── go-style.md
│   │   └── prompts/
│   │       └── go-code-review.md
│   └── team-patterns/         # Multiple sources supported
│       ├── guidelines/
│       │   └── microservices.md
│       └── prompts/
│           └── arch-review.md
```

**Rationale:**
- Source name as namespace prevents file conflicts
- Preserves relative path structure from manifest
- Clear organization for users
- Easy to see which guidelines come from which source

### Error Handling and User Experience

**Progress Indicators:**
```go
// Use lipgloss for styled output
func displayProgress(message string) {
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
    fmt.Println(style.Render("⏳ " + message))
}

func displaySuccess(message string) {
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
    fmt.Println(style.Render("✓ " + message))
}

func displayError(message string) {
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
    fmt.Println(style.Render("✗ " + message))
}
```

**Error Context:**
```go
// Always include context in errors
return fmt.Errorf("failed to clone repository %s: %w", url, err)
return fmt.Errorf("source with name '%s' already exists, use --name to specify different name", name)
return fmt.Errorf("guideline '%s' not found in manifest (available: %s)", name, strings.Join(available, ", "))
```

**Next Steps Messages:**
```
✓ Added source 'company-dna' with 2 guidelines and 1 prompt
✓ Copied files to dnaspec/company-dna/

Next steps:
  Run 'dnaspec update-agents' to configure AI agents
```

## Risks / Trade-offs

### Risk: Git clone failures due to network/auth issues

**Mitigation:**
- Clear error messages distinguishing network vs auth failures
- Document authentication setup (SSH keys, HTTPS tokens)
- Timeout prevents hanging indefinitely
- URL validation catches common mistakes

### Risk: Large repositories slow to clone

**Mitigation:**
- Shallow clone (--depth=1) significantly reduces size
- Progress indicators show operation is running
- Document that large repos may take time
- Future: could add clone caching

### Risk: Multiple sources with same derived name

**Mitigation:**
- Check for duplicate names before adding
- Clear error message with suggestion to use --name flag
- Users can override automatic naming

### Risk: File copy failures (permissions, disk space)

**Mitigation:**
- Check errors at each step
- Atomic config write (don't update config if copy fails)
- Clear error messages with file paths
- Future: could add rollback of partial copies

### Trade-off: Command-line git vs go-git library

**Chose: Command-line git**

**Pros:**
- Simpler implementation
- Better auth handling (uses system SSH/git config)
- Smaller binary
- More reliable (battle-tested)

**Cons:**
- Requires git to be installed
- Slightly slower (process spawn overhead)
- Less control over output

**Decision: Pros outweigh cons for this use case**

### Trade-off: Interactive selection vs always prompt for each guideline

**Chose: Multi-select form**

**Pros:**
- Users see all options at once
- Faster for selecting multiple guidelines
- Clear visual interface
- Can add/remove selections easily

**Cons:**
- Requires terminal UI library
- Not as simple as Y/N prompts

**Decision: Better UX justifies the complexity**

## Migration Plan

Not applicable - this is initial implementation of project commands.

**Future Compatibility:**
- Config version field (currently 1) allows future schema changes
- Source type field allows new source types in future
- Atomic writes ensure we never corrupt existing configs

## Open Questions

None - implementation path is clear from design decisions above.
