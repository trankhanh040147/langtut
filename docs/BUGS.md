
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

