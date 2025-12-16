# Development Roadmap

## Design Principles & Coding Standards

> **Reference:** All design principles, coding standards, and implementation guidelines are defined in [`.cursor/rules/rules.mdc`](../.cursor/rules/rules.mdc).

### How To Apply These Rules

Automatically loads rules from the `.cursor/rules/` directory. The `rules.mdc` file includes `alwaysApply: true` in its frontmatter, which ensures:

- **Automatic Application:** Rules are always active during coding sessions
- **Context Awareness:** Understands project-specific patterns (Vim navigation, TUI-first UX, Go conventions)
- **Consistency:** All code suggestions follow the defined principles without manual reminders

## Bug Fix Protocol

1. **Global Fix:** Search codebase (`rg`/`fd`) for similar patterns/implementations. Fix **all** occurrences, not just the reported one.
2. **Documentation:**
    - Update "Known Bugs" table (Status: Fixed).
    - Update coding standards in `.cursor/rules/rules.mdc` if the bug reflects a common anti-pattern.
3. **Testing:** Verify edge cases: Interactive, Piped (`|`), Redirected (`<`), and Non-interactive modes.
> **Reference:** Bug Fix Protocol are defined in [`.cursor/rules/rules.mdc`](../.cursor/rules/rules.mdc).

## Philosophy

- Practice and revision are the best way to learn languages (Anki, SuperMemo SM-2 algorithm)
- Interactive CLI tool enables AI-powered tutoring with ELI5 explanations
- Users can customize prompts (presets) to adapt to their preferences
- Supports learning any language

## Core Features

- **Vocabulary:** Multi-faceted approach (typing, guessing, phrases, revision)
- **Reading:** Active engagement with real content (blogs/videos)
- **Writing:** Progressive difficulty (fill-in → rewrite → freeform conversation)

---

# v0.1.0 - Foundation + Vocabulary Module

**Status:** In Progress

**Features:**

### Foundation
- [x] Cobra CLI framework
- [x] Interactive TUI with Bubbletea
- [x] AI API client (Gemini) with streaming
- [x] Config: `~/.config/langtut/config.yaml`
- [x] Preset system for custom prompts
### Vocabulary Module
- [x] Word library: Interactive TUI for CRUD operations
  - `langtut vocab add <word>` - AI generates meaning/examples, user edits in modal
  - `langtut vocab` or `langtut vocab list` - Split-pane TUI (word list + details)
  - Keyboard shortcuts: `e` edit, `d` delete, `a` add, `/` search

---

# v0.1.1 - Nested meanings

**Status:** Planned
**Focus:** Nested meanings, Part of Speech (POS), context-aware definitions, and the "Rich Add" workflow.

#### 1. Core Data Refactor (Vertical)
- [ ] **Struct Migration:** Refactor `Vocab` to replace flat definition with a slice.
    - `Definition string` $\rightarrow$ `Meanings []Meaning`.
    - `Meaning` struct includes:
        - `Type` (String: "verb", "noun", "idiom", "phrasal_verb", etc.).
        - `Definition` (String).
        - `Context` (String: Tag or Sentence).
        - `Examples` ([]String: List of usage sentences).

#### 2. Advanced CRUD (TUI)
- [ ] **Rich "Add" Modal (Fields 1-5):**
    - **Input Form:**
        1.  `Term`: The word (Locked if appending).
        2.  `Type`: **Selection List** (Verb, Noun, Adjective, Phrasal Verb, Idiom, Collocation). (*Optional*, AI will decide on context)
        3.  `Context`: A tag (e.g., "Medical") **OR** a raw sentence containing the word (*Optional*).
    - **AI Logic:** Update `GenerateDefinition`:
        - Input: `Term` + `Context` (*optional*).
        - Output: Returns structured JSON including inferred `Type`, `Definition`, and relevant `Examples`.
- [ ] **Smart Duplicate Handling:**
    - **Detection:** If `Term` exists, block "Create New".
    - **Action:** Prompt "Append Meaning?".
    - **UI:** Open Add Modal with `Term` locked, focus starts on `Type` selection.

#### 3. Enhanced Detail View
- [ ] **List Rendering:** Update `list_model` to show `Term` + `Meanings.Type` (short badge) + `Meanings.Definition` (preview).
- [ ] **Detail Pane:** Iterate over `Meanings[]` with visual hierarchy:
    - **Header:** `[Type] Definition` (e.g., `[Verb] To move fast`).
    - **Badge:** Context tag (e.g., `[Business]`).
    - **Body:** Examples rendered in italics below each definition.

#### Issues
- [ ] IS01: Wrong "Add Workflow" --> Do not type context
- [x] IS02: Legacy Type Removal --> Remove `WordInfoGenerator`, `generateWordInfo()`, and `wordInfoGeneratedMsg` from `add_model.go` once API client fully migrates to `MeaningInfoGenerator`

---
# v0.1.2 - Word Graph

**Status:** Planned
**Focus:** Word relationships, family trees, and navigation.

#### 1. Core Data Refactor (Horizontal)
- [ ] **Struct Update:** Add relational fields to `Vocab`.
    - `RootWord`: String (Parent pointer).
    - `Synonyms`: `[]string`.
    - `Antonyms`: `[]string`.
    - `Acronyms`: `[]string`.

#### 2. Advanced CRUD (TUI)
- [ ] **Extended "Add" Modal:**
    - **New Field:** Add `Synonyms / Acronyms` input (Field 5).
    - **Logic:** Accept comma-separated strings.
- [ ] **Family Linking:**
    - **Action:** "Link to Root" command in Edit Mode.
    - **UI:** Simple input to type the root word (e.g., set "happy" as root for "unhappiness").

#### 3. Networked UI
- [ ] **Wiki-Links:**
    - **Interaction:** Make Synonyms/Antonyms clickable/selectable in Detail View.
    - **Navigation:** Pressing `Enter` on a synonym jumps to that word if it exists.
- [ ] **Family Tree View:**
    - **Query:** `FindDerived(rootWord)` function.
    - **Display:** In Detail View, list all words where `RootWord == CurrentTerm`.

---

# v0.1.3 - Review

**Status:** In Progress

**Features:**

#### 1. Review Workflow

- [ ] **Context Hint:** Flashcards show `Term` + `[Context]` (e.g. "Run [Business]") so the user knows which definition to recall.

### ?
- [ ] Vocab guessing/typing: 
	- **`review_workflow`**: user types meaning → AI reviews (hint if guess wrong, user can guess again or choose to give up) --> Add word to library
	- [ ] Show word (collocation/PV, etc.) → `review_workflow`
	- [ ] Basic revise (no spaced repetition yet):  Show word from library → `review_workflow`
	
# v0.2 - Writing Module

**Status:** Planned

**Features:**

### Writing Module
- [ ] Fill-in phrases: Incomplete phrase → user fills words
- [ ] Rewrite sentences: Fragments → user creates complete sentence
- [ ] Interactive conversations: AI helps fix mistakes, improve writing

---

# v(?) 

**Status:** Planned, new features/cut off from current release

**Features:**

### Vocabulary Module

**Vocab Templates**
- User choose a template when adding a new vocab:
	- vocab add --> choose a template (select default)
- Can set default template
- Can CRUD template

**Phrase-based learning:** 
- [ ] Show phrases → prompt user for strange words → explain meaning
- [ ] User enter a phrase + strange words (optional, empty = AI break down all) --> explain meaning --> add to library

**Topic Based learning**
- [ ] User choose a topic (empty for random by AI) --> AI generate words + phrases to learn

- [ ] TTS for word pronunciation
- [ ] Export/sync library


### Reading Module
- [ ] Read blog/article with AI assistance
- [ ] Watch video with AI (URL input)
- [ ] Interactive annotation (hover for definitions, pronunciation, examples)
- [ ] Add words to library from content

---


# v0.4 - Spaced Repetition & Grammar

**Status:** Planned

**Features:**

### Spaced Repetition System
- [ ] FSRS algorithm (Free Spaced Repetition Scheduler)
- [ ] Confidence-based learning: "Forgot", "Hard", "Good", "Easy" ratings
- [ ] Review scheduling with retention tracking
- [ ] Stats command: `langtut stats --week`

### Grammar Mastery
- [ ] Grammar lesson decoder: `langtut grammar "present perfect continuous" --level beginner`
- [ ] Three difficulty levels: ELI5, Intermediate, Advanced
- [ ] Interactive grammar exercises:
  - Sentence construction (AI gives words → user arranges)
  - Error detection (user identifies & fixes errors)
  - Progressive complexity (simple → compound → subjunctive)
- [ ] Grammar by context: Learn rules specific to input text/topic

### Daily Engagement
- [ ] Streak counter
- [ ] Reminder system: `langtut remind --time 09:00 --duration 15m`
- [ ] Daily review notifications

---

# v0.5 - Conversation & Gamification

**Status:** Planned

**Features:**

### Conversational Fluency
- [ ] AI conversation partner: Multi-turn dialogue in target language
- [ ] Difficulty levels: casual chat → debate → technical discussion
- [ ] AI adapts to user vocabulary level
- [ ] Real-time feedback: Corrections marked in chat, full explanation after
- [ ] Role-play scenarios:
  - Pre-built: "Order food", "Job interview", "Travel booking"
  - Custom: User describes situation → AI generates dialogue
- [ ] Progressive scoring: Grammar → Vocabulary → Fluency → Naturalness
- [ ] Correction philosophy options:
  - Real-time interruption (perfectionists)
  - Post-conversation review (fluency-first)
  - Custom correction types

### Gamification
- [ ] Achievement system: Badges for milestones
- [ ] Language badges: "A1 Spanish", "B1 French", "C2 English"
- [ ] Skill badges: "Conversationalist", "Reader", "Grammar Master"
- [ ] Challenge badges: "Speed Demon", "Consistency"
- [ ] Profile dashboard: `langtut profile`
- [ ] XP system and progress tracking
- [ ] Optional leaderboards (privacy-respecting)

---

# v0.6 - Learning Paths & Content Integration

**Status:** Planned

**Features:**

### Structured Learning Paths
- [ ] Pre-built curricula: A1 Beginner → B1 Intermediate → C1 Advanced
- [ ] Adaptive paths: `langtut path --create "Travel Spanish"`
- [ ] Topic-based bundles: "Business English", "Medical Spanish", "Tech Japanese"
- [ ] Milestone tracking and projections

### Enhanced Content Integration
- [ ] Article difficulty filter: Auto-classify CEFR levels A1–C2
- [ ] `langtut read --level B1 --topic technology`
- [ ] News digest by language level: Same article in multiple CEFR versions
- [ ] YouTube/Podcast integration:
  - Subtitle extraction + interactive glossary
  - Quiz generation after video
  - Transcript with auto-pause for difficult phrases
- [ ] Flashcard generation from content: `langtut generate-flashcards <url> --count 20`

---

# v0.7 - Personalization & Analytics

**Status:** Planned

**Features:**

### AI Tutor Customization
- [ ] Teaching style presets: `langtut config tutor --style`
  - Socratic (questions, guides discovery)
  - Direct (rules, then examples)
  - Storytelling (narratives)
  - Gamified (challenges/quests)
  - Minimalist (brief facts)
- [ ] Personality tuning: `--personality strict|encouraging|funny`
- [ ] Difficulty curve: `--progression conservative|aggressive`
- [ ] Correction philosophy: `--error-handling strict|encouraging|selective`

### Analytics Dashboard
- [ ] `langtut analytics` command
- [ ] Learning velocity tracking
- [ ] Weakness analysis (listening, grammar, vocabulary)
- [ ] Time efficiency: Best/worst learning times
- [ ] Content preferences: Favorite/weakest topics
- [ ] Recommendations based on data

---

# v1.0 - Ecosystem & Advanced Features

**Status:** Future

**Features:**

### Ecosystem Integration
- [ ] Export to Anki: Generate .apkg files
- [ ] Cloud sync: iCloud, Google Drive, Dropbox
- [ ] Browser extension: Right-click words → add to library
- [ ] Mobile app companion: Lightweight review app
- [ ] Webhook integration: Discord, Slack, Telegram reminders

### Language Exchange
- [ ] AI matchmaking: Connect learners (A learning Spanish ↔ B learning English)
- [ ] Structured exchange prompts: Topic + conversation goals
- [ ] Asynchronous exchange: Voice messages (timezone-friendly)
- [ ] AI scoring: Grammar/fluency/engagement metrics

### Accessibility
- [ ] Keyboard-only navigation
- [ ] Voice input for answers
- [ ] Unicode support: Arabic, Chinese, Devanagari, etc.
- [ ] High-contrast themes
- [ ] Dyslexia-friendly fonts
- [ ] Adjustable text size
- [ ] Screen reader compatibility

---

# v2.0 - Future Vision

**Status:** Ideas

**Features:**

### Advanced Learning
- [ ] Multi-language support (learn multiple languages simultaneously)
- [ ] Advanced SRS algorithms and research integration
- [ ] Community features: Shared word lists, study groups
- [ ] Offline mode with local models

### Platform Expansion
- [ ] Web interface (optional)
- [ ] API for third-party integrations
- [ ] Plugin system for custom analyzers
- [ ] Self-hosted option

---

# Ideas Backlog

> Raw ideas for future consideration

**Vocabulary**
- Common misunderstanding words
- Multi-meaning words
- Enhance list:
	- Sort by date added
	- Can change sort order 

**Uncategorized**
- Compare two languages side-by-side
- Pronunciation practice with voice recognition
- Cultural context explanations
- Idiom and phrase learning
- Test preparation modes (TOEFL, IELTS, etc.)
- Custom learning schedules
- Integration with language exchange platforms
- Spaced repetition for grammar rules
- Context-aware vocabulary suggestions
- Learning streak challenges and competitions

---

# Known Bugs

> Track and fix these issues.

| **Bug**                 | **Status** | **Resolution / Notes**                                                 |
| ----------------------- | ---------- | ---------------------------------------------------------------------- |
| **StreamChat Error**    | Fixed      | `api/gemini.go`: Ignore `iterator.Done` to prevent false errors.       |
| **Map Address**         | Fixed      | `preset/preset.go`: Assign map value to temp var before addressing.    |
| **#bug01: Enter Logic** | Fixed      | `add_model.go`: Enter saves immediately if input valid.                |
| **#bug02: Tab Nav**     | Fixed      | `add_model.go`: Enabled Tab navigation between fields.                 |
| **#bug03: Stale Modal** | Fixed      | `cli/vocab.go`: Force library reload per iteration for fresh state.    |
| **#bug04: Feedback**    | Fixed      | `add_model.go`: Added "✓ Saved!" msg; triggers auto-close.             |
| **#bug05: Width Panic** | Fixed      | `ui`: Guard `strings.Repeat` against negative; default width 80.       |
| **#bug06: Edit Keys**   | Fixed      | `add_model.go`: Remapped Enter to **Edit**, `Ctrl+S` to **Save**.      |
| **#bug07: Auto-close**  | Fixed      | `add_model.go`: `saveWord` returns `tea.Quit` to signal completion.    |
| **#bug08: Duplication** | Fixed      | `add_model.go`: Cache original text/ID for safe delete/update.         |
| **#bug09: List Help**   | Fixed      | `list_model.go`: Esc closes Help overlay before other handlers.        |
| **#bug10: Index Panic** | Fixed      | `list_model.go`: Added bounds checks; handle empty lists (`idx = -1`). |
| **#bug11: Add Help**    | Fixed      | `add_model.go`: Esc closes Help overlay before closing modal.          |
| **#bug12: Blocking I/O** | Fixed      | `add_model.go`, `list_model.go`: Wrapped `vocab.Save()` in `tea.Cmd` for async I/O. |
| **#bug13: Fragile Input** | Fixed      | `cli/vocab.go`: Replaced `fmt.Scanln` with `bufio.Reader` for robust input handling. |
| **#bug14: High Complexity** | Fixed      | `add_model.go`: Extracted field editing logic into `handleFieldEditUpdate()` helper. |
| **#bug15: JSON Sanitization** | Fixed      | `api/gemini.go`: Replaced brittle string manipulation with regex for JSON extraction. |
| **#bug16: Inefficient File I/O** | Fixed      | `cli/vocab.go`: Removed library reload from loop; update library reference from saved model state. |
| **#bug17: Swallowed JSON Error** | Fixed      | `api/gemini.go`: Added error logging before fallback to legacy parser. |
| **#bug18: Imperative Slice Padding** | Fixed      | `add_model.go`: Extracted repeated slice padding logic into `setExample()` helper method. |
| **#bug19: Monolithic File** | Fixed      | `add_model.go`: Split into `add_model.go`, `add_update.go`, `add_view.go`, `add_cmds.go` (<200 lines each). |
| **#bug20: Standard JSON** | Fixed      | `vocab.go`, `gemini.go`: Replaced `encoding/json` with `sonic` for performance. |
| **#bug21: Hardcoded Keys** | Fixed      | `add_model.go`, `list_model.go`: Created `KeyMap` struct, use `key.Matches()` instead of string comparisons. |
| **#bug22: Manual Layout** | Fixed      | `add_view.go`, `list_model.go`, `help.go`: Replaced manual padding with `lipgloss.Place()`. |
| **#bug23: Inefficient Slice** | Fixed      | `list_model.go`: Pre-allocate `filteredVocabs` capacity: `make([]*vocab.Vocab, 0, len(m.vocabs))`. |
| **#bug24: Duplicate Save/Advance** | Fixed      | `add_update.go`: Extracted `saveAndAdvance()` helper to eliminate duplication. |
| **#bug25: High Complexity** | Fixed      | `add_update.go`: Split `handleFieldEditUpdate` into `handleAutocompleteNav` and `handleDefaultEdit`. |
| **#bug26: Inefficient Regex** | Fixed      | `api/gemini.go`: Moved `jsonBlockRegex` and `exampleListRegex` to package-level variables. |
| **#bug27: Inconsistent Keys** | Fixed      | `add_model.go`, `add_update.go`, `list_model.go`: Extended `KeyMap` with `Down`, `Up`, `Enter`, `Esc`, `CtrlS`; replaced all `msg.String()` with `key.Matches()`. |
| **#bug28: Redundant Reader** | Fixed      | `cli/vocab.go`: Moved `bufio.NewReader(os.Stdin)` before loop. |
| **#bug29: Legacy Code** | Fixed      | `add_model.go`, `add_cmds.go`: Removed `WordInfoGenerator`, `wordInfoGeneratedMsg`, `generateWordInfo()`, `SetWordInfoClient()` per IS02. |
| **#bug30: Unstructured Logs** | Fixed      | `api/gemini.go`, `vocab/vocab.go`: Standardized warning logs to structured format `level=warn msg="..." err="..."`. |
| **#bug31: Inconsistent Key Binding** | Fixed      | `list_model.go`: Replaced `switch msg.String()` block with `key.Matches()` checks. Extended `KeyMap` with all missing bindings. |
| **#bug32: Inefficient TUI Instantiation** | Fixed      | `cli/vocab.go`: Created `batchAddModel` to manage word queue in single `tea.Program` instead of creating new program per word. |
| **#bug33: Manual Layout Calculation** | Fixed      | `list_model.go`, `add_view.go`: Replaced `strings.Repeat("─", count)` with `renderHorizontalRule()` helper using `lipgloss` borders. |
| **#bug34: Monolithic File Structure** | Fixed      | `list_model.go`: Split into `list_model.go` (151 lines), `list_update.go` (210 lines), `list_view.go` (227 lines). |
| **#bug35: Redundant Key Constants** | Fixed      | `constants/ui.go`: Removed legacy `Key*` string constants. All key handling now uses `KeyMap` struct with `key.Matches()`. |
| **#bug36: Mixed Key Handling Styles** | Fixed      | `add_model.go`, `base.go`: Replaced remaining `msg.String()` comparisons with `key.Matches()`. Standardized all key handling. |

---
> **Reminder**: Contents written in this file need to be condensed. Remove fluff, preserve meaning, maintain clarity for machine processing.
