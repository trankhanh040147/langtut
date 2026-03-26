# Current Base - Utilities & Flows

## Overview

This document describes the current state of the codebase: utilities, services, commands, TUI components, and flows. Use this as a reference when planning new features.

## Commands & Subcommands

### Root Command
- **Command:** `prepf` (binary name: `langtut`)
- **Description:** AI assistant for software development
- **Flags:**
  - `-c, --cwd`: Current working directory
  - `-D, --data-dir`: Custom data directory
  - `-d, --debug`: Debug mode
  - `-h, --help`: Help
  - `-y, --yolo`: Auto-accept all permissions (dangerous)

### Subcommands

#### `prepf run [prompt...]`
- **Purpose:** Run single non-interactive prompt
- **Flags:**
  - `-q, --quiet`: Hide spinner
  - `-m, --model`: Model to use (format: `model` or `provider/model`)
  - `--small-model`: Small model override
- **Flow:** Reads prompt from args or stdin, executes via LLM, outputs to stdout

#### `prepf login [platform]`
- **Purpose:** Authenticate with platform
- **Platforms:** `hyper`, `copilot`, `github`, `github-copilot`
- **Flow:** OAuth device flow, stores tokens in config

#### `prepf dirs [config|data]`
- **Purpose:** Print directories used by Prepf
- **Subcommands:**
  - `prepf dirs config`: Config directory
  - `prepf dirs data`: Data directory

#### `prepf projects`
- **Purpose:** List project directories
- **Flags:**
  - `--json`: Output as JSON
- **Flow:** Reads from projects registry, displays table or JSON

#### `prepf models`
- **Purpose:** Manage LLM models

#### `prepf logs`
- **Purpose:** View logs

#### `prepf schema`
- **Purpose:** Schema operations

#### `prepf stats`
- **Purpose:** Statistics

#### `prepf update-providers`
- **Purpose:** Update provider configurations

## Core Services

### Database Service (`internal/db`)
- **Type:** SQLite (modernc.org/sqlite or ncruces/go-sqlite3)
- **Location:** `{dataDir}/prepf.db`
- **Migrations:** Goose-based migrations in `internal/db/migrations/`
- **Tables:**
  - `sessions`: Session metadata (id, title, message_count, tokens, cost, timestamps)
  - `messages`: Messages per session (id, session_id, role, parts, model, timestamps)
  - `files`: File versions per session (id, session_id, path, content, version, timestamps)

### Session Service (`internal/session`)
- **Interface:** `Service`
- **Methods:**
  - `Create(ctx, title) (Session, error)`
  - `Get(ctx, id) (Session, error)`
  - `List(ctx) ([]Session, error)`
  - `Update(ctx, id, title) error`
  - `Delete(ctx, id) error`
  - `UpdateStats(ctx, id, stats) error`

### Message Service (`internal/message`)
- **Interface:** `Service`
- **Methods:**
  - `Create(ctx, sessionID, role, parts, model) (Message, error)`
  - `Get(ctx, id) (Message, error)`
  - `ListBySession(ctx, sessionID) ([]Message, error)`
  - `Update(ctx, id, parts) error`
  - `Delete(ctx, id) error`

### File History Service (`internal/history`)
- **Interface:** `Service`
- **Methods:**
  - `Create(ctx, sessionID, path, content) (File, error)`
  - `CreateVersion(ctx, sessionID, path, content) (File, error)`
  - `Get(ctx, id) (File, error)`
  - `GetByPathAndSession(ctx, path, sessionID) (File, error)`
  - `ListBySession(ctx, sessionID) ([]File, error)`
  - `ListLatestSessionFiles(ctx, sessionID) ([]File, error)`
  - `Delete(ctx, id) error`
  - `DeleteSessionFiles(ctx, sessionID) error`
- **PubSub:** Implements `pubsub.Subscriber[File]` for file change events

### Permission Service (`internal/permission`)
- **Interface:** `Service`
- **Purpose:** Manages tool execution permissions
- **Methods:**
  - `RequestPermission(ctx, tool, args) (bool, error)`
  - `IsAllowed(tool) bool`
- **Features:**
  - Skip requests if `yolo` mode enabled
  - Session-scoped permissions
  - Allowed tools whitelist

## LLM Providers & Configuration

### Supported Providers
- **OpenAI** (`openai`)
- **Anthropic** (`anthropic`)
- **Google/Gemini** (`gemini`)
- **Azure** (`azure`)
- **OpenRouter** (`openrouter`)
- **OpenAI Compatible** (`openai-compat`)
- **AWS Bedrock** (`bedrock`)
- **Google Vertex** (`vertexai`)
- **Hyper** (`hyper`) - OAuth-based

### Provider Configuration (`internal/config`)
- **Location:** `~/.config/langtut/config.yaml`
- **Structure:**
  - `providers[]`: List of provider configs
  - Each provider: `id`, `name`, `type`, `base_url`, `api_key`, `oauth`, `extra_headers`, `extra_body`
- **Model Selection:** Per-agent model configuration
- **Agent Types:** `coder`, `task` (extensible)

### Agent Coordinator (`internal/agent`)
- **Purpose:** Orchestrates LLM interactions, tool execution, session management
- **Features:**
  - Multi-model support (large + small)
  - Tool execution with permissions
  - Auto-summarization for long contexts
  - Session title generation
  - Message queuing
  - Token/cost tracking

## Prompt System

### Prompt Templates (`internal/agent/templates/`)
- **Format:** Markdown templates with Go template syntax
- **Templates:**
  - `coder.md.tpl`: Main coding assistant prompt
  - `task.md.tpl`: Task-oriented prompt
  - `initialize.md.tpl`: Project initialization prompt
  - `title.md`: Session title generation
  - `summary.md`: Context summarization
  - `agentic_fetch_prompt.md.tpl`: Web content analysis agent
- **Template Data:**
  - `Provider`, `Model`, `Config`
  - `WorkingDir`, `IsGitRepo`, `Platform`, `Date`
  - `GitStatus`, `ContextFiles[]`, `AvailSkillXML`

### Prompt Builder (`internal/agent/prompt`)
- **Type:** `Prompt`
- **Methods:**
  - `NewPrompt(name, template, opts...) (*Prompt, error)`
  - `Build(ctx, provider, model, cfg) (string, error)`
- **Options:**
  - `WithTimeFunc(fn)`: Custom time function
  - `WithPlatform(platform)`: Platform override
  - `WithWorkingDir(dir)`: Working directory override

## TUI Components

### Architecture
- **Framework:** Bubbletea v2
- **UI Modes:**
  - Legacy: `internal/tui` (old UI)
  - New: `internal/ui` (new UI, enabled via `PREPF_NEW_UI=true`)

### Global Keymaps

#### Root Keymap (`internal/tui/keys.go`, `internal/ui/model/keys.go`)
- `ctrl+c`: Quit
- `ctrl+g`: Help/More
- `ctrl+p`: Commands dialog
- `ctrl+z`: Suspend
- `ctrl+l` / `ctrl+m`: Models dialog
- `ctrl+s`: Sessions dialog
- `tab`: Change focus

### Chat Page Keymaps (`internal/tui/page/chat/keys.go`, `internal/ui/model/keys.go`)

#### Editor
- `enter`: Send message
- `shift+enter` / `ctrl+j`: Newline
- `ctrl+o`: Open external editor
- `ctrl+f`: Add attachment/image
- `/`: Add file
- `@`: Mention file
- `ctrl+r`: Delete attachment mode
- `esc` / `alt+esc`: Cancel/Escape

#### Chat Navigation
- `ctrl+n`: New session
- `ctrl+d`: Toggle details
- `ctrl+space`: Toggle pills/tasks
- `←` / `→`: Switch section
- `j` / `k`: Navigate messages (vim-style)
- `g` / `G`: First/Last message
- `/`: Search
- `ctrl+c`: Copy
- `home` / `end`: Scroll to top/bottom
- `pageup` / `pagedown`: Page navigation

### Dialog Keymaps

#### Commands Dialog (`internal/tui/components/dialogs/commands/keys.go`)
- `ctrl+p`: Open commands
- `ctrl+g`: Toggle help
- `ctrl+c`: Quit
- `ctrl+o`: Open external editor (if `$EDITOR` set)
- `ctrl+f`: File picker (if session active + model supports images)

#### Models Dialog (`internal/tui/components/dialogs/models/keys.go`)
- Navigation keys for model selection

#### Sessions Dialog (`internal/ui/dialog/sessions.go`)
- List navigation for session selection

#### Permissions Dialog (`internal/tui/components/dialogs/permissions/keys.go`)
- `a` / `A` / `ctrl+a`: Allow
- `s` / `S` / `ctrl+s`: Allow session
- `d` / `D` / `esc`: Deny
- `enter` / `ctrl+y`: Confirm
- `t`: Toggle diff mode
- `←` / `→` / `h` / `l`: Navigate
- `shift+↓` / `shift+↑` / `J` / `K`: Scroll
- `shift+←` / `shift+→` / `H` / `L`: Horizontal scroll

#### File Picker (`internal/tui/components/dialogs/filepicker/keys.go`)
- Navigation and selection keys

### TUI Components Structure

#### Core Components (`internal/tui/components/core/`)
- **Layout:** Main layout manager
- **Status:** Status bar component

#### Chat Components (`internal/tui/components/chat/`)
- **Editor:** Text input with attachments
- **Messages:** Message rendering and caching
- **Header:** Chat header
- **Sidebar:** Session sidebar
- **Splash:** Landing screen
- **Todos:** Task management

#### Dialog Components (`internal/tui/components/dialogs/`)
- **Commands:** Command palette
- **Models:** Model selection
- **Sessions:** Session management
- **Permissions:** Permission requests
- **File Picker:** File selection
- **Quit:** Quit confirmation
- **Reasoning:** Reasoning display

#### List Components (`internal/tui/exp/list/`)
- **Filterable:** Filterable list with search
- **Grouped:** Grouped list items
- **Items:** List item interfaces

#### New UI Components (`internal/ui/`)
- **Model:** Main UI model (`model/ui.go`)
- **Chat:** Chat logic (`model/chat.go`)
- **Dialog:** Dialog system (`dialog/`)
- **List:** Generic list component (`list/`)
- **Common:** Shared utilities (`common/`)
- **Styles:** Style definitions (`styles/`)

## Application Flow

### Interactive Mode (`prepf` without args)
1. **Initialization:**
   - Load config from `~/.config/langtut/config.yaml`
   - Connect to SQLite DB at `{dataDir}/prepf.db`
   - Run migrations if needed
   - Initialize services (Session, Message, History, Permission)
   - Initialize Agent Coordinator
   - Initialize LSP clients (background)
   - Initialize MCP clients (background)

2. **TUI Startup:**
   - Create Bubbletea program
   - Choose UI mode (legacy vs new based on `PREPF_NEW_UI`)
   - Subscribe app events to TUI
   - Run TUI program

3. **Session Flow:**
   - User creates/selects session (`ctrl+n` or sessions dialog)
   - User types message in editor
   - On `enter`, message sent to Agent Coordinator
   - Agent Coordinator:
     - Creates message in DB
     - Builds prompt from template + context
     - Calls LLM provider
     - Streams response back to TUI
     - Executes tools if needed (with permission prompts)
     - Updates session stats (tokens, cost)
   - Response displayed in chat
   - Auto-summarization if context too large

4. **Tool Execution:**
   - Agent requests tool execution
   - Permission Service checks if allowed
   - If not allowed, shows permission dialog
   - User approves/denies
   - Tool executes
   - Results returned to agent

### Non-Interactive Mode (`prepf run`)
1. **Setup:** Same as interactive (config, DB, services)
2. **Execution:**
   - Read prompt from args or stdin
   - Create temporary session
   - Call Agent Coordinator with prompt
   - Stream response to stdout
   - Exit

## Configuration System

### Config File Structure
- **Location:** `~/.config/langtut/config.yaml`
- **Schema:** JSON schema in `schema.json`
- **Key Sections:**
  - `providers[]`: LLM provider configurations
  - `agents{}`: Agent-specific configs (coder, task)
  - `options{}`: Global options (data_dir, debug, etc.)
  - `permissions{}`: Permission settings

### Config Utilities (`internal/config`)
- **Load:** `config.Init(cwd, dataDir, debug) (*Config, error)`
- **Resolve:** Variable resolution (env vars, etc.)
- **Provider Management:** Get, list, update providers
- **Model Selection:** Resolve models by name/type

## Utilities & Helpers

### File Operations (`internal/fsext`)
- File system extensions

### String Utilities (`internal/stringext`)
- String manipulation helpers

### Shell Utilities (`internal/shell`)
- Shell command execution
- Environment variable resolution

### Format Utilities (`internal/format`)
- Text formatting
- Spinner for non-interactive mode

### Logging (`internal/log`)
- Structured logging with slog
- HTTP request logging for debug mode

### Event System (`internal/event`)
- Application lifecycle events
- Metrics/analytics events

### PubSub (`internal/pubsub`)
- Event broker for service communication
- Used by File History Service

### Projects (`internal/projects`)
- Project directory tracking
- Registry of known projects

### Update Check (`internal/update`)
- Version checking
- Update notifications

## Data Directory Structure

```
~/.config/langtut/
├── config.yaml          # Main configuration
└── data/
    ├── prepf.db         # SQLite database
    └── projects/        # Project-specific data
```

## Key Design Patterns

### Service Pattern
- Services implement interfaces
- Services use database queries (`db.Queries`)
- Services can publish/subscribe to events

### TUI Pattern
- Components are "dumb" renderers
- Main model handles state and routing
- Commands for side effects (IO)
- Messages for state updates

### Agent Pattern
- Coordinator manages agent lifecycle
- Agents use prompt templates
- Tools are registered and executed via permissions
- Streaming responses to TUI

### Configuration Pattern
- YAML config with JSON schema validation
- Environment variable resolution
- Provider abstraction over LLM APIs

## Notes for v0.1.0 Planning

### Reusable Components
- **Session Service:** Can store writing drill sessions
- **Message Service:** Can store Q&A pairs
- **Agent Coordinator:** Can be adapted for examiner persona
- **Prompt Templates:** Can create new templates for writing modes
- **TUI Components:**
  - Editor component for text input
  - Dialog system for modals
  - List components for vocab selection
  - Message rendering for feedback display

### Gaps to Address
- **No vocabulary library service yet** (planned in v0.3.0)
- **No session report generation** (need to add)
- **No structured feedback format** (need to define)
- **No band score storage** (need schema extension)
- **No "mode" selection UI** (need dialog/command)

### Potential Extensions
- Add `writing_drill` agent type to config
- Create prompt templates for examiner persona
- Extend session schema for drill metadata
- Create vocab harvesting service (or extend existing)
- Add report generation utility
