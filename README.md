# peektea

![cover](media/cover.png)

A minimal terminal file browser built with [Bubble Tea](https://github.com/charmbracelet/bubbletea). Navigate your filesystem from the terminal using arrow keys (or vim keys).

## Demo

![demo](media/demo.gif)

## Install

```bash
git clone https://github.com/lovestaco/peektea
cd peektea
make install   # puts peektea in ~/go/bin
```

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
| `q` / `ctrl+c` | quit |

## Setup

Run `peektea init` to configure which apps open each file type. It peeks into your installed software and lets you pick.

![peektea init](media/peek_tea_init.gif)

## Help

![peektea -h](media/help.png)

## Development

```bash
make build   # build ./peektea
make start   # live reload via air (rebuilds on every .go save)
make install # install to ~/go/bin
```

Requires [air](https://github.com/air-verse/air) for `make start` (`go install github.com/air-verse/air@latest`).

## Stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework (Elm Architecture)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling
