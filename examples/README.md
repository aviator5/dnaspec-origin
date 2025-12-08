# DNASpec Examples

This directory contains example manifest and project configuration files to help you get started.

## Manifest Examples

These are examples of `dnaspec-manifest.yaml` files for DNA repository maintainers.

###

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

## Project Configuration Examples

These are examples of `dnaspec.yaml` files for project developers using DNA guidelines.

### `project-config-simple.yaml`

A basic project configuration with:
- Single DNA source from a git repository
- One guideline and one prompt
- Minimal setup

**Use this when:**
- Starting with DNASpec in your project
- Using guidelines from a single source
- You want a simple, clean configuration

### `project-config-multi-source.yaml`

A comprehensive project configuration with:
- Multiple DNA sources (git repository and local directory)
- Multiple guidelines and prompts
- Agent configuration for Claude and GitHub Copilot

**Use this when:**
- Using DNA guidelines from multiple sources
- Combining company-wide and team-specific guidelines
- You need a more complex setup

### `project-config-local-only.yaml`

A project configuration using only local sources:
- DNA guidelines from a local directory
- Useful for development and testing

**Use this when:**
- Developing DNA guidelines locally
- Testing guidelines before publishing
- Working offline

## Using Manifest Examples

1. Copy an example to your DNA repository:
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

## Using Project Configuration Examples

1. Initialize your project (creates empty `dnaspec.yaml`):
   ```bash
   dnaspec init
   ```

2. Add DNA sources as needed:
   ```bash
   # From git repository
   dnaspec add --git-repo https://github.com/company/dna-guidelines
   
   # From local directory
   dnaspec add /path/to/local/dna
   ```

The `dnaspec add` command automatically updates your `dnaspec.yaml` file. You can also manually edit the file by referencing these examples.

## Note

**For manifests:** These example manifests reference guideline and prompt files that don't exist. You'll need to create those files before the manifest will validate successfully. The file paths in these examples follow the required conventions:

**For project configs:** The example configurations reference sources and files that may not exist on your system. When you use `dnaspec add`, the tool automatically populates the configuration with the correct paths and metadata.

- Guideline files must be in the `guidelines/` directory
- Prompt files must be in the `prompts/` directory
- All paths are relative to the manifest file location
