# Session changes — 2026-06-06

Everything created/changed in this session, problem by problem.

## 1. `CLAUDE.md` (new)

Guidance file for Claude Code: build/run commands, architecture overview
(model update flow, config key derivation, preview loading, manual layout math).

## 2. `peektea: command not found` after `make install`

**Problem:** binary landed in `~/go/bin`, which wasn't on `$PATH`.

**Fix — `Makefile` (`install` target):** after building, it now checks whether
the binary actually resolves, and self-heals:

1. `command -v peektea` works → done.
2. `~/.local/bin` is on `$PATH` (default on most distros) → symlink the binary
   there. **Works immediately in the current shell** — no `source`, no new
   terminal. (`make` runs as a child process, so it can never modify the
   running shell's environment directly; the symlink trick sidesteps that.)
3. Last resort → append `export PATH="$HOME/go/bin:$PATH"` to `.bashrc`/`.zshrc`
   (idempotent, picks rc file by `$SHELL`) and print the one `source` command
   needed.

## 3. `open failed: exec: "xdg-open": executable file not found` (WSL)

**Problem:** WSL has no Linux GUI apps and no `xdg-open`; pressing `o` on
directories/images/PDFs failed.

**Fix — WSL support, three parts:**

- **`internal/config/wsl.go` (new):**
  - `IsWSL()` — detects WSL via `$WSL_DISTRO_NAME` or `/proc/version`.
  - `WSLOpener()` — best Windows-side opener: `wslview` (from `wslu`) if
    installed, else `explorer.exe` (PATH or `/mnt/c/Windows/explorer.exe`).
  - `WindowsPath()` — converts `/home/...` → `\\wsl.localhost\Ubuntu\...` via
    `wslpath -w`; Windows programs can't read Linux paths.
- **`main.go` (`o` handler):** when the configured program ends in `.exe`, the
  file path is auto-converted with `WindowsPath()` before launching.
  (explorer.exe always exits 1 — harmless, GUI launches already ignore exit codes.)
- **`internal/cmd/init.go`:** on WSL, the `xdg-open` fallbacks (file manager /
  image viewer / PDF viewer) become the Windows opener, with a
  "WSL detected" notice.
- **`~/.peektea.toml` patched in place:** all `xdg-open` entries →
  `/mnt/c/Windows/explorer.exe`, so `o` works without re-running init.

## 4. `peektea init` user-friendliness

**Problem:** chafa missing → init printed a wall of per-distro install
commands and left the work to the user. Also, declining the overwrite prompt
aborted init entirely.

**Fix — `internal/cmd/init.go`:**

- `offerInstall()` / `pkgInstallCmd()` (new): detects the system package
  manager (`apt` → `dnf` → `pacman` → `zypper` → `apk` → `brew`) and offers to
  install chafa on the spot — shows the exact command, bare enter = yes,
  re-checks afterwards and confirms previews are ready. Decline prints the
  command for later; no package manager found prints a graceful fallback.
- Declining `already exists. Overwrite? [y/N]` now keeps the existing config
  and **continues** to the chafa offer instead of exiting — so
  `peektea init` can be re-run just to install extras.

## Verification done

- `go vet ./...` clean; built and installed via `make install`.
- `make install` symlink branch tested live (branch 2 fired, `peektea -h`
  resolved immediately); rc-append branch tested against a throwaway `$HOME`.
- `explorer.exe` + `wslpath -w` opening a real file confirmed from WSL.
- Package-manager detection and the decline path of `offerInstall` covered by
  temporary Go tests (run, passed, removed).
- Full decline-overwrite → chafa-offer flow exercised end to end.

## Follow-ups (optional)

- `sudo apt install wslu` → gives `wslview`, slightly cleaner than explorer.exe;
  re-run `peektea init` afterwards to prefer it.
- Run `peektea init` and answer **Y** to actually install chafa for inline
  image previews (`p` key).
- README still says "Make sure `~/go/bin` is on your `$PATH`" — now stale, and
  the new WSL behavior + init auto-install aren't documented there yet.
