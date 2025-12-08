package config

import (
	"os"
)

// ExampleManifestYAML returns an example manifest YAML content with helpful comments
const ExampleManifestYAML = `# DNASpec Manifest
# This file defines the guidelines and prompts available in this DNA repository.

version: 1

guidelines:
  # Example guideline entry
  - name: go-style
    file: guidelines/go-style.md
    description: Go language style conventions and best practices
    applicable_scenarios:
      - "Writing Go code"
      - "Reviewing Go code"
      - "Setting up Go projects"
    prompts:
      - code-review
      - implementation

  # Add more guidelines here
  # - name: rest-api
  #   file: guidelines/rest-api.md
  #   description: RESTful API design guidelines
  #   applicable_scenarios:
  #     - "Designing REST APIs"
  #     - "Implementing API endpoints"

prompts:
  # Example prompt entry
  - name: code-review
    file: prompts/code-review.md
    description: Prompt for conducting thorough code reviews

  - name: implementation
    file: prompts/implementation.md
    description: Prompt for implementing new features

  # Add more prompts here
  # - name: debugging
  #   file: prompts/debugging.md
  #   description: Prompt for systematic debugging
`

// CreateExampleManifest creates an example manifest file at the given path
func CreateExampleManifest(path string) error {
	return os.WriteFile(path, []byte(ExampleManifestYAML), 0o644)
}

// ExampleProjectYAML returns an example project config YAML content with helpful comments
const ExampleProjectYAML = `# DNASpec Project Configuration
# This file configures which DNA guidelines are active in your project.

version: 1

# AI agents that should use these guidelines
# agents:
#   - "claude-code"
#   - "github-copilot"

# DNA sources - guidelines from repositories or local directories
# sources:
#   - name: company-dna
#     type: git-repo
#     url: https://github.com/company/dna
#     ref: v1.2.0
#     commit: abc123def456789...
#     guidelines:
#       - name: go-style
#         file: guidelines/go-style.md
#         description: Go code style conventions
#         applicable_scenarios:
#           - "writing new Go code"
#         prompts:
#           - code-review
#     prompts:
#       - name: code-review
#         file: prompts/code-review.md
#         description: Review Go code
`

// CreateExampleProjectConfig creates an example project config file at the given path
func CreateExampleProjectConfig(path string) error {
	return os.WriteFile(path, []byte(ExampleProjectYAML), 0o644)
}
