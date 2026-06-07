# Navigation & keys

| Key                     | Action                          |
| ----------------------- | ------------------------------- |
| `↑` / `k`               | move up                         |
| `↓` / `j`               | move down                       |
| `→` / `l` / `enter`     | go inside directory             |
| `←` / `h` / `backspace` | go to parent                    |
| `H`                     | go to home directory            |
| `o`                     | open with configured program    |
| `p`                     | toggle preview panel            |
| `[` / `]`               | scroll preview up / down        |
| `/`                     | filter entries as you type      |
| `esc`                   | exit filter / clear active filter |
| `.`                     | toggle hidden files (dotfiles)  |
| `s`                     | cycle sort: name → size → modified |
| `q` / `ctrl+c`          | quit                            |

## Opening files

Press `o` (or `enter` on a file) to open the selected entry with the program configured for its extension. See [Configuration](../configuration.md) for how to map extensions to programs — and what happens for extensions you haven't configured.

Programs listed in `terminal_programs` (like `vim` or `nvim`) take over the full screen while open. Everything else — GUI apps like `nautilus` or `feh` — launches in the background so the browser keeps running.
