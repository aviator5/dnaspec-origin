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
	return os.WriteFile(path, []byte(ExampleManifestYAML), 0644)
}
