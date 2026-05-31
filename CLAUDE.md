# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

**peektea** — a terminal file browser TUI written in Go, using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework. The binary starts in the current working directory and lets the user navigate the filesystem with arrow keys.

## Commands

```bash
go build -o peektea .   # build the binary
./peektea               # run it
go mod tidy             # sync go.mod/go.sum after adding/removing deps
```

## Architecture

Single-file app (`main.go`) following Bubble Tea's **Model / Update / View** pattern:

- **`model`** — holds `dir` (current path), `entries` (result of `os.ReadDir`), `cursor` (selected row index), and `err`.
- **`Update`** — handles all key input: `↑/k`, `↓/j` move the cursor; `→/l/enter` descends into the selected directory; `←/h/backspace` ascends to the parent and restores the cursor to the directory just left; `q/ctrl+c` quits.
- **`View`** — renders with Lipgloss styles; directories are shown with a trailing `/`. Uses `strings.Builder` to accumulate lines.

The program runs with `tea.WithAltScreen()` so it takes over the terminal and restores it cleanly on exit.

## Key bindings

| Key | Action |
|-----|--------|
| `↑` / `k` | move up |
| `↓` / `j` | move down |
| `→` / `l` / `enter` | open directory |
| `←` / `h` / `backspace` | go to parent |
| `q` / `ctrl+c` | quit |
