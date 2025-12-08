# project-management Specification Delta

## ADDED Requirements

### Requirement: List Project Configuration Command

The system SHALL provide a `dnaspec list` command that displays all configured DNA sources, guidelines, prompts, and AI agents from the project configuration.

#### Scenario: Display full configuration successfully

- **WHEN** user runs `dnaspec list` in a project with valid `dnaspec.yaml`
- **THEN** display "Configured Agents" section showing all agent IDs or "None configured" if empty
- **AND** display "Sources:" header
- **AND** for each source, display source name with type in parentheses
- **AND** for git-repo sources, display URL, Ref, and Commit fields
- **AND** for local-dir sources, display Path field
- **AND** for each source, display "Guidelines:" section with indented list of guidelines (name: description)
- **AND** for each source, display "Prompts:" section with indented list of prompts (name: description)
- **AND** exit with code 0

#### Scenario: Display configuration with no agents

- **WHEN** user runs `dnaspec list` and configuration has empty or missing agents array
- **THEN** display "Configured Agents" section with message indicating no agents configured
- **AND** continue to display sources normally
- **AND** exit with code 0

#### Scenario: Display configuration with no sources

- **WHEN** user runs `dnaspec list` and configuration has empty sources array
- **THEN** display configured agents if any
- **AND** display "Sources:" header with message indicating no sources configured
- **AND** exit with code 0

#### Scenario: Display source with no guidelines

- **WHEN** source has empty guidelines array
- **THEN** display source metadata normally
- **AND** display "Guidelines:" section with message indicating no guidelines or empty list
- **AND** continue to display prompts if any

#### Scenario: Display source with no prompts

- **WHEN** source has empty prompts array
- **THEN** display source metadata and guidelines normally
- **AND** display "Prompts:" section with message indicating no prompts or empty list

#### Scenario: Configuration file not found

- **WHEN** user runs `dnaspec list` in directory without `dnaspec.yaml`
- **THEN** display error message indicating configuration file not found
- **AND** suggest running `dnaspec init` to create configuration
- **AND** exit with non-zero code

#### Scenario: Malformed configuration file

- **WHEN** user runs `dnaspec list` and `dnaspec.yaml` contains invalid YAML syntax
- **THEN** display error message with YAML parsing error details
- **AND** exit with non-zero code

#### Scenario: Display multiple sources of different types

- **WHEN** configuration contains both git-repo and local-dir sources
- **THEN** display each source with appropriate type-specific fields
- **AND** maintain consistent formatting across all sources
- **AND** display sources in order they appear in configuration
