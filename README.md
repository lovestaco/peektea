
![cover](media/cover.png)

A minimal terminal file browser built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). 

Peek through your filesystem with arrow keys (or vim keys), then pour each file straight into the app you've configured for it.

## Demo

A quick peek before you steep:

![demo](media/demo.gif)

## Install
 
Brew it from source:

```bash
git clone https://github.com/lovestaco/peektea
cd peektea
make install   # puts peektea in ~/go/bin
```

Make sure `~/go/bin` is on your `$PATH`.

## Usage

```bash
peektea
```

Starts in the current working directory.

## Keys

| Key | Action |
|-----|--------|
| `↑` / `k` | move up |
| `↓` / `j` | move down |
| `→` / `l` / `enter` | go inside directory |
| `←` / `h` / `backspace` | go to parent |
| `o` | open with configured program |
| `p` | toggle preview panel |
| `/` | filter entries as you type |
| `esc` | exit filter / clear active filter |
| `.` | toggle hidden files (dotfiles) |
| `q` / `ctrl+c` | quit |

## Filter

Press `/` to enter filter mode. Type anything and the list narrows to matching entries in real time. The filter input appears at the bottom of the panel above the hint bar — like vim's command line.

- `enter` — confirm and exit filter mode (filter stays active)
- `esc` — clear the filter entirely
- `↑` / `↓` still navigate while you type

Press `.` to toggle hidden files (dotfiles) on and off. The hint bar always shows the current state: `. show hidden` or `. hide hidden`. Both filters compose — you can search by name with dotfiles visible or hidden at the same time.

## Preview

Press `p` to open a side-by-side preview panel. Press `p` again to close it.

- **Text files** — first N lines rendered inline, truncated to panel width
- **Images** — rendered directly in the terminal via [chafa](https://hpjansson.org/chafa/)
- **Directories** — lists the contents of the folder
- **Binary files** — shows a `[binary file]` notice

The left panel auto-widens to fit the longest filename in the current directory. `peektea init` will tell you if chafa is installed and how to get it if not.

## Setup

Run `peektea init` to configure which apps open each file type. 

It peeks into your installed software and lets you pick your blend.

![peektea init](media/peek_tea_init.gif)

## Configuration

`init` writes `~/.peektea.toml`, but you can steep it by hand. Each key is derived straight from the file extension — dots become underscores, wrapped with `_` and `_config`:

```
file.md        → _md_config
archive.tar.gz → _tar_gz_config
hello.xd.dd    → _xd_dd_config
directory      → _dir_config
```

Example, brewed to taste:

```toml
_md_config      = "vim"
_png_config     = "feh"
_dir_config     = "nautilus"
_default_config = "less"

terminal_programs = ["vim", "nvim", "nano", "micro", "hx"]
```

Fallback order when opening a file — the bag never comes up empty:

1. The matching `_<ext>_config` key
2. `_default_config` for unknown extensions
3. `vim` if nothing is configured at all

**`terminal_programs`** tells peektea which programs need the full terminal. Anything in this set (vim, nvim, etc.) takes over the screen while open — the whole table, not just a saucer. Everything else (GUI apps like nautilus or feh) launches in the background and keeps the browser brewing.

## Help

```bash
peektea -h
```

Consider it the menu before you order:

![peektea -h](media/help.png)

## Uninstall

```bash
peektea uninstall
```

For when you've had your fill. 

It shows you exactly what it's about to delete, asks for confirmation, then asks separately whether you also want `~/.peektea.toml` removed. 

No silent nuke — no tea spilled without warning.

## Development

```bash
make build   # build ./peektea
make start   # live reload via air (rebuilds on every .go save)
make install # install to ~/go/bin
```

Requires [air](https://github.com/air-verse/air) for `make start` (`go install github.com/air-verse/air@latest`).

## Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework (Elm Architecture)
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components (textinput for filter)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling

---

## since you peeked this far — check out what I'm really brewing

I'm Maneshwar, and alongside peektea I'm building [git-lrc](https://github.com/HexmosTech/git-lrc) — a micro AI code reviewer that runs on every commit.

It hooks into `git commit`, reviews every diff before it lands, and catches the bugs AI agents introduce silently. 60-second setup.

AI agents write code fast.

They also silently remove logic, change behavior, and introduce bugs — without telling you.

You often find out in production.

git-lrc fixes this.

Any feedback or contributors are welcome.

⭐ [Star git-lrc on GitHub](https://github.com/HexmosTech/git-lrc) to help other devs discover it.
