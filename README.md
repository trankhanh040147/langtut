# Langtut

[![Build](https://img.shields.io/github/actions/workflow/status/trankhanh040147/langtut/build.yml?branch=main&label=build)](https://github.com/trankhanh040147/langtut/actions/workflows/build.yml)
[![Release](https://img.shields.io/github/v/release/trankhanh040147/langtut)](https://github.com/trankhanh040147/langtut/releases)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](./LICENSE)

## Description

Langtut is a terminal AI assistant for development and language-learning workflows. It provides an interactive TUI, one-shot CLI execution, multiple working modes, provider-backed model access, and local session persistence.

## Features

- Interactive terminal UI for chat-driven workflows
- Non-interactive execution with `langtut run`
- Multiple modes, including coding assistance and writing tutor sessions
- Provider-backed model selection with support for API key and OAuth-based setup
- MCP integration for external tools and prompts
- LSP integration for language-aware project workflows
- Local session, message, and history persistence via SQLite
- Project tracking, per-project data directories, and log inspection commands
- JSON configuration with generated schema support

## Prerequisites

- Go `1.25.5` or newer
- At least one configured model provider credential, for example:
  - `OPENAI_API_KEY`
  - `ANTHROPIC_API_KEY`
  - `GOOGLE_API_KEY`

## Installation

### Recommended: install with `go install`

```bash
go install github.com/trankhanh040147/langtut@latest
```

Verify the installation:

```bash
langtut --version
```

### Alternative: clone and install locally

```bash
git clone https://github.com/trankhanh040147/langtut.git
cd langtut
go install .
```

### Last resort: build from source

```bash
git clone https://github.com/trankhanh040147/langtut.git
cd langtut
go build -o langtut .
./langtut --version
```

## Configuration

Langtut reads configuration from JSON files and environment variables.

### Minimal configuration

Create `langtut.json` in the project root, or use the global config path.

```json
{
  "$schema": "https://charm.land/langtut.json",
  "providers": {
    "openai": {
      "api_key": "$OPENAI_API_KEY"
    }
  }
}
```

### Configuration file locations

Langtut loads configuration from these locations:

1. Global config:
   - `~/.config/langtut/langtut.json`
   - or `$XDG_CONFIG_HOME/langtut/langtut.json`
2. Global data-backed override config:
   - `~/.local/share/langtut/langtut.json`
   - or `$XDG_DATA_HOME/langtut/langtut.json`
3. Project config discovered upward from the current working directory:
   - `langtut.json`
   - `.langtut.json`

### Runtime data directory

Per-project runtime data is stored in:

```text
.langtut/
```

### Useful environment variables

```bash
export OPENAI_API_KEY="..."
export ANTHROPIC_API_KEY="..."
export GOOGLE_API_KEY="..."
```

Optional overrides:

```bash
export XDG_CONFIG_HOME="$HOME/.config"
export XDG_DATA_HOME="$HOME/.local/share"
export CRUSH_DISABLE_PROVIDER_AUTO_UPDATE="1"
export CRUSH_DISABLE_METRICS="1"
export DO_NOT_TRACK="1"
```

### OAuth-based setup

You can authenticate interactively instead of placing raw tokens in config:

```bash
langtut login
langtut login copilot
```

## Usage

### Start the interactive interface

```bash
langtut
```

Expected result:

```text
An interactive terminal session starts.
```

### Run a one-shot prompt

```bash
langtut run --quiet "Summarize the current repository"
```

Example output:

```text
This repository contains a terminal AI assistant with an interactive UI, one-shot CLI mode, provider configuration, and local persistence.
```

### Review piped input

```bash
git diff --staged | langtut run "Review these staged changes"
```

Example output:

```text
The response is printed to stdout as plain text or markdown.
```

### Show config and data directories

```bash
langtut dirs
```

Example output:

```text
/home/user/.config/langtut
/home/user/.local/share/langtut
```

### List tracked projects

```bash
langtut projects --json
```

Example output:

```json
{
  "projects": [
    {
      "path": "/path/to/project",
      "data_dir": "/path/to/project/.langtut",
      "last_accessed": "2026-03-22T12:34:56Z"
    }
  ]
}
```

### List available models

```bash
langtut models
langtut models gpt
```

Example output:

```text
openai/gpt-5
openai/gpt-5-mini
```

### Inspect logs

```bash
langtut logs
langtut logs --follow
```

## Architecture

High-level package layout:

- `main.go`: process entry point
- `internal/cmd`: Cobra command tree
- `internal/app`: application wiring, lifecycle, permissions, MCP, LSP, updates
- `internal/agent`: model and tool orchestration
- `internal/tui`: Bubble Tea terminal UI
- `internal/config`: layered config loading and provider resolution
- `internal/db`: SQLite connection, migrations, and queries
- `internal/session`, `internal/message`, `internal/history`: persisted runtime state

Execution flow:

1. Resolve working directory and configuration.
2. Load and merge config files.
3. Resolve providers and selected models.
4. Initialize persistence, permissions, MCP, and LSP clients.
5. Start either the TUI or a one-shot non-interactive run.

## Contributing

- Use issues for bug reports and feature proposals
- Keep pull requests focused and reviewable
- Run checks before opening a PR:

```bash
go test ./...
```

If `task` is installed, the maintained workflow is:

```bash
task test
task lint
task fmt
```

For bug reports, include:

- exact command
- relevant configuration
- expected result
- actual result
- logs or reproduction steps

## License

MIT. See [`LICENSE`](./LICENSE).
