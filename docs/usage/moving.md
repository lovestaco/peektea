# Moving files

peektea uses a cut → navigate → paste flow so you can move any file or directory without typing paths.

| Key   | Action                                      |
| ----- | ------------------------------------------- |
| `x`   | cut — stage the highlighted entry for move  |
| `v`   | paste — move the staged entry here          |
| `esc` | cancel the staged move                      |

## Flow

1. Highlight a file or directory and press `x`. The hint bar confirms what's staged:
   `moving: report.md  (v to drop here, esc to cancel)`
2. Navigate to the destination using the normal keys (`↑`/`↓`, `→`/`←`, `H`, `/`).
3. Press `v`. The entry is moved and the list refreshes with it selected.

## Edge cases

| Situation | What happens |
| --------- | ------------ |
| Name already exists at destination | Inline prompt: `'name' exists here — overwrite? (y/n)`. `y` overwrites, `n` keeps the entry staged so you can navigate elsewhere. `esc` cancels the move entirely. |
| Paste in the same directory | No-op — shows "already here" and clears the staged move. |
| Moving a directory into itself or a descendant | Refused with a clear message; the source is untouched. |
| Source and destination on different filesystems | `os.Rename` fails with `EXDEV`; peektea falls back to a recursive copy + delete automatically. |
| Permission denied | Error shown in the hint bar; the source is left untouched. |
