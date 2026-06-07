# Development

```bash
git clone https://github.com/lovestaco/peektea
cd peektea
```

| Command           | What it does                                              |
| ----------------- | --------------------------------------------------------- |
| `make build`      | build `./peektea`                                         |
| `make install`    | build and install to `~/go/bin` (and wire up `$PATH`)     |
| `make start`      | live reload via [air](https://github.com/air-verse/air) — rebuilds on every `.go` save |
| `make snapshot`   | build release archives locally without publishing         |
| `make release`    | tag and publish a GitHub release via [goreleaser](https://goreleaser.com) |

## Requirements

- [air](https://github.com/air-verse/air) for `make start`:
  ```bash
  go install github.com/air-verse/air@latest
  ```
- [goreleaser](https://goreleaser.com) for `make release` / `make snapshot`:
  ```bash
  go install github.com/goreleaser/goreleaser/v2@latest
  ```

## Stack

peektea is built on the [Charm](https://charm.sh/) stack:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework (Elm Architecture)
- [Bubbles](https://github.com/charmbracelet/bubbles) — TUI components (`textinput` for the filter)
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- [BurntSushi/toml](https://github.com/BurntSushi/toml) — config parsing

## Releasing

1. `make release-notes-init VERSION=v0.x.x` — scaffold `release-notes/v0.x.x.md` from the template
2. Fill in the release notes (use `IMG:media/file.gif` for images — rewritten to raw GitHub URLs at publish time)
3. `make release-notes-preview VERSION=v0.x.x` — preview the rendered notes without publishing
4. `make bump VERSION=v0.x.x` — tag and push
5. `make release VERSION=v0.x.x` — build with goreleaser and publish the GitHub release with rendered notes

## Contributing

Issues and pull requests are welcome on [GitHub](https://github.com/lovestaco/peektea).
