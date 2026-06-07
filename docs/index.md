---
title: peektea
description: A minimal terminal file browser built with Bubble Tea
---

# peektea

![cover](media/cover.png)

A minimal terminal file browser built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Peek through your filesystem with arrow keys (or vim keys), then pour each file straight into the app you've configured for it.

## A quick peek before you steep

![demo](media/demo.gif)

## Why peektea

- **Keyboard-first** — arrow keys or vim keybindings (`hjkl`), no mouse required
- **Smart previews** — text with syntax highlighting via [bat](https://github.com/sharkdp/bat), images via [chafa](https://hpjansson.org/chafa/), directory listings, and binary detection
- **Real-time filtering** — narrow the list as you type, toggle hidden files on demand
- **Configurable openers** — route every file extension to the program of your choice
- **WSL aware** — automatically routes file opens through `wslview` or `explorer.exe` with path conversion

## Get started

<div class="grid cards" markdown>

-   :material-download:{ .lg .middle } **Install**

    ---

    One-liner script, prebuilt binaries, `go install`, or build from source.

    [:octicons-arrow-right-24: Installation](installation.md)

-   :material-keyboard:{ .lg .middle } **Learn the keys**

    ---

    Navigate, filter, preview, and sort — all from the keyboard.

    [:octicons-arrow-right-24: Usage](usage/index.md)

-   :material-cog:{ .lg .middle } **Configure openers**

    ---

    Run `peektea init` and tell peektea which app opens which file type.

    [:octicons-arrow-right-24: Configuration](configuration.md)

-   :material-microsoft-windows:{ .lg .middle } **Using WSL?**

    ---

    peektea detects WSL automatically and opens files in your Windows apps.

    [:octicons-arrow-right-24: WSL support](wsl.md)

</div>

## Stack

peektea is built on the [Charm](https://charm.sh/) stack:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework (Elm Architecture)
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components (`textinput` for filtering)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling

## Source

peektea is open source on [GitHub](https://github.com/lovestaco/peektea). Issues, stars, and contributions are welcome.
