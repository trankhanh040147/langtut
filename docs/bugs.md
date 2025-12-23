
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

