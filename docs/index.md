---
title: peektea
description: A minimal terminal file browser built with Bubble Tea
---

# peektea

 

A minimal terminal file browser built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

Peek through your filesystem with arrow keys (or vim keys), then pour each file straight into the app you've configured for it.

## A quick peek before you steep

<div class="video-wrapper">
  <iframe src="https://www.youtube.com/embed/yDNN5x9Y9Ok" title="peektea demo" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" allowfullscreen></iframe>
</div>

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

    [:octicons-arrow-right-24: Installation](installation/index.md)

-   :material-keyboard:{ .lg .middle } **Learn the keys**

    ---

    Navigate, filter, preview, and sort — all from the keyboard.

    [:octicons-arrow-right-24: Usage](usage/index.md)

-   :material-cog:{ .lg .middle } **Configure openers**

    ---

    Run `peektea init` and tell peektea which app opens which file type.

    [:octicons-arrow-right-24: Configuration](installation/configuration.md)

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

<img src="https://u8views.com/api/v1/github/profiles/66487268/views/day-week-month-total-count.svg" width="0" height="0" alt="">
