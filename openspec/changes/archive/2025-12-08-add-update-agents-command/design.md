# Design: Add Update-Agents Command

## Architectural Overview

The `update-agents` command bridges DNA guidelines (markdown files) with AI agent integrations (agent-specific file formats). It implements a **generator pattern** where each agent type has a dedicated generator that produces files in that agent's expected format.

```
┌─────────────────┐
│  dnaspec.yaml   │
│  (sources +     │
│   guidelines)   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────┐
│  update-agents Command      │
│  ┌───────────────────────┐  │
│  │ Agent Selection UI    │  │
│  └───────────┬───────────┘  │
│              ▼              │
│  ┌───────────────────────┐  │
│  │ Agent File Generator  │  │
│  │  Orchestrator         │  │
│  └───────┬───────────────┘  │
│          │                  │
│  ┌───────┴───────────┐      │
│  ▼                   ▼      │
│  Claude Generator  Copilot  │
│                   Generator │
└─────────┬──────────┬────────┘
          │          │
          ▼          ▼
    .claude/    .github/prompts/
    CLAUDE.md   (copilot files)

    AGENTS.md (always generated)
```

## Key Design Decisions

### 1. Managed Block System

**Problem**: Users need to add custom project-specific instructions to AGENTS.md and CLAUDE.md, but DNASpec also needs to update these files when guidelines change.

**Solution**: Managed blocks with clear markers:

```markdown
# My Project Instructions

(user content here)

<!-- DNASPEC:START -->
... auto-generated content ...
<!-- DNASPEC:END -->

(more user content here)
```

**Implementation Strategy**:

1. **Detection**: Scan file for `<!-- DNASPEC:START -->` and `<!-- DNASPEC:END -->` markers
2. **Replacement**: If markers found, replace content between them, preserve everything else
3. **Append**: If file exists but no markers, append managed block at end
4. **Create**: If file doesn't exist, create with minimal header + managed block

**Trade-offs**:
- ✅ Users can add custom content safely
- ✅ Predictable update behavior
- ✅ Clear ownership boundaries
- ❌ Users who edit inside markers lose changes (documented risk)
- ❌ Requires marker parsing logic

**Alternative Considered**: Separate files (e.g., `AGENTS.md` for DNASpec, `AGENTS_CUSTOM.md` for users)
- Rejected because: AI agents expect single file, users would need to maintain references

### 2. Agent Abstraction

**Problem**: Different AI agents have different file formats and locations. We need to support multiple agents without coupling the core logic.

**Solution**: Agent generator pattern with separate modules per agent:

```
internal/core/agents/
├── registry.go          # Agent definitions and discovery
├── generate.go          # Orchestration (calls generators)
├── agents_md.go         # AGENTS.md (all agents)
├── claude_md.go         # CLAUDE.md (Claude only)
├── claude_commands.go   # Claude slash commands
└── copilot_prompts.go   # Copilot prompt files
```

**Generator Interface** (conceptual, not enforced as Go interface):

```go
// Each generator implements:
// GenerateFiles(config ProjectConfig, sourceName string, prompt ProjectPrompt) error
```

**Orchestration** (`generate.go`):

```go
func GenerateAgentFiles(config ProjectConfig, agents []string) error {
    // Always generate AGENTS.md
    GenerateAgentsMD(config)

    // Generate agent-specific files
    for _, agent := range agents {
        switch agent {
        case "claude-code":
            GenerateClaudeMD(config)
            for each source {
                for each prompt {
                    GenerateClaudeCommand(source, prompt)
                }
            }
        case "github-copilot":
            for each source {
                for each prompt {
                    GenerateCopilotPrompt(source, prompt)
                }
            }
        }
    }
}
```

**Trade-offs**:
- ✅ Easy to add new agents (add new generator module)
- ✅ Clear separation of concerns
- ✅ Each agent format is isolated
- ❌ More files/modules than a single monolithic generator
- ❌ Some duplication (e.g., frontmatter parsing)

**Alternative Considered**: Plugin system with agent generators as plugins
- Rejected because: Overkill for initial implementation, Go doesn't have great plugin story

### 3. Source Namespacing

**Problem**: Multiple sources might have prompts with the same name, causing file conflicts.

**Solution**: Include source name in all generated filenames:

- Claude: `.claude/commands/dnaspec/<source-name>-<prompt-name>.md`
- Copilot: `.github/prompts/dnaspec-<source-name>-<prompt-name>.prompt.md`

**Example**:
```
Sources: company-dna, team-patterns
Both have: code-review prompt

Generated:
.claude/commands/dnaspec/company-dna-code-review.md
.claude/commands/dnaspec/team-patterns-code-review.md
```

**Trade-offs**:
- ✅ Eliminates conflicts completely
- ✅ Clear source attribution in filenames
- ✅ User can tell which source a prompt came from
- ❌ Longer filenames
- ❌ Renaming source breaks existing files (future: cleanup command)

**Alternative Considered**: Directory structure by source
```
.claude/commands/dnaspec/company-dna/code-review.md
.claude/commands/dnaspec/team-patterns/code-review.md
```
- Rejected because: Claude Code expects flat command structure

### 4. Applicable Scenarios Usage

**Problem**: AI agents need to know **when** to consult each guideline, not just that guidelines exist.

**Solution**: Use `applicable_scenarios` from manifest to generate context-aware instructions in AGENTS.md:

```markdown
When working on the codebase, open and refer to the following DNA guidelines as needed:
- `@/dnaspec/company-dna/guidelines/go-style.md` for
   * writing new Go code
   * refactoring existing Go code
- `@/dnaspec/company-dna/guidelines/rest-api.md` for
   * designing API endpoints
   * implementing HTTP handlers
```

**Benefits**:
- AI agents see explicit triggers for when to open guidelines
- Reduces cognitive load (agent knows which guideline for which task)
- Makes guidelines discoverable based on current work context

**Requirements**:
- Manifest validation must enforce non-empty applicable_scenarios (already done in manifest-management spec)
- Guideline authors must write good, specific scenarios

### 5. File Generation Strategy

**Problem**: Generated files need to be updated when guidelines change, but must preserve user content and be idempotent.

**Solution**: Layered file generation approach:

**Layer 1: Content Generation**
- Pure functions that generate content strings
- Input: ProjectConfig + agent selection
- Output: String content
- Testable without filesystem

**Layer 2: File Operations**
- Detect existing file state
- Apply managed block logic
- Atomic writes (temp file + rename)

**Layer 3: Orchestration**
- Iterate sources and prompts
- Call generators
- Collect errors
- Report summary

**Idempotency Requirements**:
- Running `update-agents` multiple times produces same result
- User content outside managed blocks is preserved
- File permissions and timestamps updated, but content unchanged

**Testing Strategy**:
- Unit test Layer 1 (content generation)
- Unit test Layer 2 (file operations with temp directories)
- Integration test Layer 3 (full workflow)

### 6. Agent Selection Persistence

**Problem**: Users shouldn't have to re-select agents every time they run `update-agents`.

**Solution**: Persist selection in `dnaspec.yaml`:

```yaml
version: 1
agents:
  - "claude-code"
  - "github-copilot"
sources: [...]
```

**Modes**:
1. **Interactive** (default): Show current selection, allow changes, save to config
2. **Non-interactive** (`--no-ask`): Use saved selection, error if not set

**Use Cases**:
- **First time**: Interactive mode, user selects agents
- **After adding guidelines**: Non-interactive mode (`update-agents --no-ask`)
- **CI/CD**: Non-interactive mode (reproducible)
- **Changing agents**: Interactive mode (re-run without `--no-ask`)

### 7. Error Handling Philosophy

**Principle**: Collect all errors, continue processing, report at end.

**Rationale**:
- User wants to know about ALL failures, not just the first one
- Partial success is better than no success
- Example: If 9/10 prompt files generate successfully, generate those 9

**Implementation**:
```go
func GenerateAgentFiles(config ProjectConfig, agents []string) (summary Summary, errors []error) {
    // Generate AGENTS.md
    if err := GenerateAgentsMD(config); err != nil {
        errors = append(errors, err)
        // Continue anyway
    }

    // Generate per-agent files
    for _, agent := range agents {
        // ... generate files, collect errors
    }

    return summary, errors
}
```

**Error Display**:
```
✓ Generated AGENTS.md
✓ Generated CLAUDE.md
✓ Generated 8/10 Claude commands
  ✗ Failed: company-dna-review (file permission denied)
  ✗ Failed: team-patterns-lint (prompt file missing)

2 errors occurred. Review and fix before using commands.
```

## Future Extensibility

### Adding New Agents

To add a new agent (e.g., Windsurf):

1. Add agent definition to `registry.go`
2. Create `internal/core/agents/windsurf_*.go` with generator functions
3. Add case to orchestration in `generate.go`
4. Add tests for new generator
5. Update documentation

**No changes required**:
- Command structure
- Managed block system
- Configuration format (agents array already exists)

### Agent Cleanup

**Future Enhancement**: Detect when agent is deselected and offer to remove files:

```
You deselected "claude-code". Remove generated files? [y/N]
  - CLAUDE.md
  - .claude/commands/dnaspec/ (5 files)
```

**Implementation Notes**:
- Track generated files in config (future: `.dnaspec-generated.json`)
- Compare current agent selection with previous
- Offer removal on deselection

### Managed Block Versioning

**Future Enhancement**: Version managed blocks for migrations:

```markdown
<!-- DNASPEC:START:v2 -->
... content ...
<!-- DNASPEC:END:v2 -->
```

**Use Case**: Format changes requiring migration (e.g., path format changes)

## Security Considerations

### Prompt Content

**Risk**: Prompts are user-defined markdown files that could contain malicious content.

**Mitigation**:
- DNASpec treats prompts as opaque content (doesn't execute)
- AI agents are responsible for safe prompt handling
- Manifest validation prevents path traversal in prompt file references

### File Overwrites

**Risk**: Managed block replacement could destroy user content if markers are placed incorrectly.

**Mitigation**:
- Clear documentation about managed blocks
- Marker detection is robust (exact string match)
- Content outside markers is always preserved
- Atomic writes prevent partial corruptions

### Directory Creation

**Risk**: Creating directories in arbitrary locations (e.g., `.claude/commands/dnaspec/`).

**Mitigation**:
- All paths are relative to project root
- No user-supplied path components in generated file locations
- Use `filepath.Join` to prevent path traversal

## Performance Considerations

### File Generation Scale

**Scenario**: Project with 5 sources × 10 prompts × 2 agents = 100 files + 2 markdown files

**Performance**:
- File writes are I/O bound, not CPU bound
- Parallel generation possible (future optimization)
- Current: Sequential is fine for typical scale (<1 second for 100 files)

**Future Optimization**: Use goroutines for parallel file generation:
```go
var wg sync.WaitGroup
for _, source := range sources {
    for _, prompt := range prompts {
        wg.Add(1)
        go func(s, p) {
            defer wg.Done()
            GenerateClaudeCommand(s, p)
        }(source, prompt)
    }
}
wg.Wait()
```

### Config Loading

**Current**: Load config once at command start

**Future**: If config becomes large (many sources), consider lazy loading or partial loading

## Testing Strategy

### Unit Tests

1. **Content Generation** (no I/O):
   - Test AGENTS.md content format
   - Test applicable_scenarios rendering
   - Test frontmatter generation

2. **Managed Block Logic** (in-memory strings):
   - Test detection with various content
   - Test replacement preserving outside content
   - Test append and create modes

3. **Orchestration Logic** (mocked generators):
   - Test agent selection
   - Test file counts
   - Test error collection

### Integration Tests

1. **Full Workflow** (temp directory):
   - Create test config with sources and guidelines
   - Run update-agents
   - Verify all files created
   - Verify content correctness

2. **Idempotency**:
   - Run update-agents twice
   - Verify identical results

3. **Preservation**:
   - Create AGENTS.md with user content outside markers
   - Run update-agents
   - Verify user content preserved

### Manual Tests

1. **Real Agent Validation**:
   - Generate files for Claude Code
   - Open project in Claude Code
   - Verify commands appear and work
   - Verify AGENTS.md instructions are followed

2. **Real Copilot Validation**:
   - Generate files for GitHub Copilot
   - Open project in VSCode with Copilot
   - Verify prompts are discovered
   - Test prompt invocation

## Summary

This design prioritizes:

1. **User Control**: Managed blocks preserve user content
2. **Extensibility**: Easy to add new agents without core changes
3. **Safety**: Idempotent operations, atomic writes, error collection
4. **Clarity**: Clear separation between DNASpec-managed and user content
5. **Testing**: Layered design enables comprehensive testing

The implementation uses straightforward Go patterns (separate packages, pure functions, error collection) and avoids premature abstraction while maintaining extensibility for future agents.
