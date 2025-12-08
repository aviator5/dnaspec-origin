# Implementation Tasks: Add List Command

## 1. Command Structure

- [x] 1.1 Create internal/cli/project/list.go with NewListCmd() cobra command
- [x] 1.2 Define command structure with Use, Short, Long, and Example fields
- [x] 1.3 Add RunE function to execute list command logic
- [x] 1.4 Verify file compiles without errors

## 2. Configuration Loading

- [x] 2.1 Load project configuration using config.LoadProjectConfig("dnaspec.yaml")
- [x] 2.2 Handle missing configuration file with helpful error message
- [x] 2.3 Suggest running "dnaspec init" when config not found
- [x] 2.4 Handle YAML parsing errors with descriptive messages
- [x] 2.5 Exit with appropriate error codes for error scenarios

## 3. Agents Display

- [x] 3.1 Implement displayAgents() function to format agents section
- [x] 3.2 Display "Configured Agents (Phase 1):" header
- [x] 3.3 List each agent ID with "  - " prefix
- [x] 3.4 Handle empty agents array with "None configured" or similar message
- [x] 3.5 Verify output matches design doc format (docs/design.md lines 810-812)

## 4. Sources Display

- [x] 4.1 Implement displaySources() function to iterate through all sources
- [x] 4.2 Display "Sources:" header with blank line separator
- [x] 4.3 For each source, display name with type in parentheses (e.g., "my-company-dna (git-repo)")
- [x] 4.4 For git-repo sources, display URL, Ref, and Commit fields with proper indentation
- [x] 4.5 For local-dir sources, display Path field with proper indentation
- [x] 4.6 Handle empty sources array with appropriate message
- [x] 4.7 Verify output matches design doc format (docs/design.md lines 816-820)

## 5. Guidelines Display

- [x] 5.1 Implement displayGuidelines() function for a single source
- [x] 5.2 Display "  Guidelines:" header with proper indentation
- [x] 5.3 For each guideline, display "    - name: description" format
- [x] 5.4 Handle empty guidelines array gracefully (skip section or show empty message)
- [x] 5.5 Verify output matches design doc format (docs/design.md lines 822-824)

## 6. Prompts Display

- [x] 6.1 Implement displayPrompts() function for a single source
- [x] 6.2 Display "  Prompts:" header with proper indentation
- [x] 6.3 For each prompt, display "    - name: description" format
- [x] 6.4 Handle empty prompts array gracefully (skip section or show empty message)
- [x] 6.5 Verify output matches design doc format (docs/design.md lines 826-827)

## 7. Output Formatting and Styling

- [x] 7.1 Import internal/ui package for consistent styling
- [x] 7.2 Apply appropriate styles for headers and content
- [x] 7.3 Ensure proper spacing between sections
- [x] 7.4 Use lipgloss styles consistently with other commands

## 8. Command Registration

- [x] 8.1 Register NewListCmd() in cmd/dnaspec/main.go
- [x] 8.2 Verify command appears in "dnaspec --help" output
- [x] 8.3 Verify command can be invoked with "dnaspec list"

## 9. Unit Tests

- [x] 9.1 Create internal/cli/project/list_test.go
- [x] 9.2 Test full configuration with agents, multiple sources, guidelines, and prompts
- [x] 9.3 Test empty agents array
- [x] 9.4 Test empty sources array
- [x] 9.5 Test source with no guidelines
- [x] 9.6 Test source with no prompts
- [x] 9.7 Test git-repo source displays URL, Ref, Commit
- [x] 9.8 Test local-dir source displays Path
- [x] 9.9 Test mixed source types (git and local)
- [x] 9.10 Run tests with "go test ./internal/cli/project/" and verify all pass

## 10. Integration Tests

- [x] 10.1 Create test fixture with sample dnaspec.yaml file
- [x] 10.2 Test end-to-end: create config, run list command, verify output
- [x] 10.3 Test missing configuration file scenario
- [x] 10.4 Test malformed YAML scenario
- [x] 10.5 Verify integration tests pass

## 11. Manual Testing

- [x] 11.1 Build binary: go build -o dnaspec ./cmd/dnaspec
- [x] 11.2 Test with full configuration (multiple sources, agents, guidelines, prompts)
- [x] 11.3 Test with empty configuration (no sources or agents)
- [x] 11.4 Test with missing dnaspec.yaml file
- [x] 11.5 Test with malformed YAML file
- [x] 11.6 Verify output formatting and colors in terminal
- [x] 11.7 Verify all error messages are clear and helpful

## Notes

- Follow existing command patterns from internal/cli/project/init.go
- Use config.LoadProjectConfig() for configuration loading
- Use internal/ui styles for consistent terminal output
- Keep implementation simple - this is primarily formatting and display logic
- No complex business logic required
- Exit codes: 0 for success, non-zero for errors
