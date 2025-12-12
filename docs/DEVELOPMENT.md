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

# v0.1 - Vocabulary Module

**Status:** In Progress

**Features:**

### Vocabulary Module
- [x] Word library: Interactive TUI for CRUD operations
  - `langtut vocab add <word>` - AI generates meaning/examples, user edits in modal
  - `langtut vocab` or `langtut vocab list` - Split-pane TUI (word list + details)
  - Keyboard shortcuts: `e` edit, `d` delete, `a` add, `/` search
- [ ] Vocab guessing/typing: Show word → user types meaning → AI reviews
- [ ] Phrase-based learning: Show phrases → prompt user for strange words → explain meaning
- [ ] Export/sync library
- [ ] TTS for word pronunciation
- [ ] Topic-based learning
- [ ] Daily random vocab

### Foundation
- [x] Cobra CLI framework
- [x] Interactive TUI with Bubbletea
- [x] AI API client (Gemini) with streaming
- [x] Config: `~/.config/langtut/config.yaml`
- [x] Preset system for custom prompts

---

# v0.2 - Reading Module

**Status:** Planned

**Features:**

### Reading Module
- [ ] Read blog/article with AI assistance
- [ ] Watch video with AI (URL input)
- [ ] Interactive annotation (hover for definitions, pronunciation, examples)
- [ ] Add words to library from content

---

# v0.3 - Writing Module

**Status:** Planned

**Features:**

### Writing Module
- [ ] Fill-in phrases: Incomplete phrase → user fills words
- [ ] Rewrite sentences: Fragments → user creates complete sentence
- [ ] Interactive conversations: AI helps fix mistakes, improve writing

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

| Bug | Status | Notes |
|-----|--------|-------|
| StreamChat treats iterator.Done as error | Fixed | Fixed in `internal/api/gemini.go` - now checks `if err == iterator.Done` before treating as exception |
| Cannot take address of map index expression | Fixed | Fixed in `internal/preset/preset.go:36` - assign to variable before taking address |
| #bug01: Cannot save words after Enter until last step | Fixed | Fixed in `internal/ui/vocab/add_model.go:195-210` - Enter now saves immediately if word is valid and not editing |
| #bug02: Cannot use Tab to navigate fields | Fixed | Fixed in `internal/ui/vocab/add_model.go:218-225` - Tab now navigates between fields when not editing |
| #bug03: Multi-word add modal only shows last word | Fixed | Fixed in `internal/cli/vocab.go:114-143` - Library reloaded at start of each iteration, ensuring fresh state for each modal |

---
> **Reminder**: Contents written in this file need to be condensed. Remove fluff, preserve meaning, maintain clarity for machine processing.
