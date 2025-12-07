# DNASpec Manifest Examples

This directory contains example `dnaspec-manifest.yaml` files to help you get started.

## Examples

### `minimal-manifest.yaml`

The simplest valid manifest with:
- One guideline
- No prompts
- Minimal configuration

**Use this when:**
- Getting started with DNASpec
- You only need basic guideline management
- You want to understand the minimum required fields

### `complete-manifest.yaml`

A full-featured manifest demonstrating:
- Multiple guidelines with different configurations
- Guidelines with and without prompts
- Multiple applicable scenarios
- Prompt definitions and cross-references

**Use this when:**
- You need a comprehensive example
- You want to see all manifest features
- You're building a complex DNA repository

### `go-project-manifest.yaml`

A manifest tailored for Go projects with:
- Go-specific guidelines (style, testing, concurrency, etc.)
- Relevant prompts for Go development
- Scenarios specific to Go programming

**Use this when:**
- Working on a Go project
- You want inspiration for Go-specific guidelines
- You need a starting point for a Go DNA repository

## Using These Examples

1. Copy an example to your project:
   ```bash
   cp examples/minimal-manifest.yaml dnaspec-manifest.yaml
   ```

2. Edit the manifest to match your needs:
   - Update guideline names and descriptions
   - Add or remove guidelines and prompts
   - Modify applicable scenarios

3. Create the referenced files:
   ```bash
   mkdir -p guidelines prompts
   # Create your guideline and prompt files
   ```

4. Validate your manifest:
   ```bash
   dnaspec manifest validate
   ```

## Note

These example manifests reference guideline and prompt files that don't exist. You'll need to create those files before the manifest will validate successfully. The file paths in these examples follow the required conventions:

- Guideline files must be in the `guidelines/` directory
- Prompt files must be in the `prompts/` directory
- All paths are relative to the manifest file location
