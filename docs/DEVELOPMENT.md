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

# v0.1.0 - Writing Drill Mode

**Status:** In Planning

**Features:**

- [ ] **Session Flow:** 
  - AI acts as Examiner (Part 1/2/3 style questions).
  - User answers in text.
- [ ] **"The Red Pen" Feedback Engine:** - **Immediate Linting:** Analyzes each response for:
    - *Lexical Range:* Highlights repeated words (e.g., "good", "nice") and suggests C1 synonyms.
    - *Grammar:* Detects tense inconsistencies.
    - *Coherence:* Checks for linking words (However, Therefore, Consequently).
  - **Band Score Estimator:** Real-time projected score (e.g., "Current Level: 6.5").
- [ ] **"Stop-and-Fix" Mechanic:** - If a specific sentence contains a critical error, the AI prompts a "Quick Fix" before proceeding to the next question.
- [ ] **Knowledge Injection (End-of-Session):**
  - [ ] **Auto-Summarizer:** Generates `session_report_<date>.md` with mistakes vs. corrections.
  - [ ] **Vocab Harvester:** Parses "Golden Phrases" (native-like expressions suggested by AI) and presents an interactive checklist to `[Add to Library]`. (deferred to v0.?)

---

# v0.2.0 - Writing Trainer Modes

**Status:** In Planning

**Features:**

### SOTA Writing Trainer Modes:
- [ ] **The "Refiner" (Diff View):** - User responds to a prompt → AI generates "Native Ideal" → TUI shows `git diff` style comparison (Red for errors, Green for improvements).
- [ ] **Semantic Validator:** - AI validates answers based on *meaning*, not string matching (allows for synonyms/variations).
- [ ] **Style Transfer:** - Drills focusing on register changes (Formal ↔ Casual).

#### A. The "Reverse Translation" Loop (The Interpolator)

Instead of asking the user to translate "Hello" -> "Hola", give them a target *context*.

1. **AI:** Presents a sentence in the User's Native Language (e.g., "I'm not sure if I can make it to the meeting.").
2. **User:** Attempts to write it in Target Language.
3. **AI Analysis:** Does not look for an exact string match. It evaluates **Semantic Equivalence**.
4. **Feedback:** If the meaning is close but the grammar is off, the AI highlights *only* the grammatical friction points.

#### B. "The Refiner" (Git Diff Learning)

This is the "Killer Feature" for a CLI tool.

1. **Prompt:** "Write a complaint about cold coffee."
2. **User Input:** "Coffee cold. I no like."
3. **AI SOTA Response:** It generates the "Native Speaker Version" ("The coffee is cold and I'm not happy about it.") and displays a **colored character-level Diff** between the user's input and the ideal output.
* *Why it works:* It's instant visual feedback on syntax and preposition usage.

#### C. "Style Transfer" Drills

1. **AI:** "Here is a formal sentence: 'I require assistance.' Rewrite this as if you are texting a close friend."
2. **User:** "Help me."
3. **AI:** Evaluates **Tone** and **Register**, not just grammar.

---

# v0.3.0 - Vocabulary Module

**Status:** In Planning

**Features:**

### Vocabulary Module
- [ ] Word library: Interactive TUI for CRUD operations
  - `langtut vocab add <word>` - AI generates meaning/examples, user edits in modal
  - `langtut vocab` or `langtut vocab list` - Split-pane TUI (word list + details)
  - Keyboard shortcuts: `e` edit, `d` delete, `a` add, `/` search

---

# v0.4.0 - Reading Module

**Status:** In Progress

**Features:**

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

### IELTS Writing Module (The "Essay Compiler")
- [ ] **Blueprint Mode:** Interactive form to outline Thesis and Body paragraphs before writing.
- [ ] **Live Linter:** TUI overlay showing word count, sentence variety score, and cohesion markers in real-time.
- [ ] **Lexical Upgrader:** Post-write analysis that identifies "weak" verbs/adjectives and suggests C1/C2 alternatives (e.g., "bad" -> "detrimental").
- [ ] **Task 1 Data Generator:** AI generates a text-based description of a graph (e.g., "A bar chart showing..."), user must write the report.

---

# Known Bugs

> Track and fix these issues.
@docs/BUGS.md

---
> **Reminder**: Contents written in this file need to be condensed. Remove fluff, preserve meaning, maintain clarity for machine processing.
