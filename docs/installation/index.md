# Installation for Linux/macOS/WSL

Pick whichever method fits your setup.

## One-liner 

```bash
curl -fsSL https://raw.githubusercontent.com/lovestaco/peektea/master/scripts/install.sh | sh
```

Downloads the right binary for your platform and places it on your `$PATH`.

## Install with Go

```bash
go install github.com/lovestaco/peektea@latest
```

## Download a binary

No Go toolchain required — grab the latest release for your platform from the [releases page](https://github.com/lovestaco/peektea/releases/latest):

| Platform            | File                              |
| ------------------- | --------------------------------- |
| Linux x86-64        | `peektea_*_linux_amd64.tar.gz`    |
| Linux arm64         | `peektea_*_linux_arm64.tar.gz`    |
| macOS x86-64        | `peektea_*_darwin_amd64.tar.gz`   |
| macOS Apple Silicon | `peektea_*_darwin_arm64.tar.gz`   |

Extract the archive and put the `peektea` binary anywhere on your `$PATH`.


## Build from source

```bash
git clone https://github.com/lovestaco/peektea
cd peektea
make install
```

`make install` builds the binary and puts it in `~/go/bin`, then figures out your `$PATH` for you:

1. **Already reachable** — nothing to do, you're set.
2. **`~/.local/bin` is on your `$PATH`** — symlinks the binary there, works immediately in the current shell.
3. **Neither** — appends `~/go/bin` to your `.bashrc`/`.zshrc` and tells you which file to `source`.

## Next steps

<div class="grid cards" markdown>

-   [:material-cog: **Configuration**](configuration.md)

    Run `peektea init` to wire up which apps open each file type, see it in action with a demo, and learn the config file it writes.

-   [:material-trash-can: **Uninstall**](uninstall.md)

    When you've had your fill — see what gets removed and how to confirm before it's gone.

</div>
