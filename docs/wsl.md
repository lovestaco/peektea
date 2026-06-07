# WSL support

![WSL support init](media/init_wsl_support.gif)

peektea works on Windows Subsystem for Linux out of the box. It detects WSL automatically and routes file opens through:

1. [`wslview`](https://wslutiliti.es/wslu/) (from `wslu`), if available
2. `explorer.exe`, otherwise

Linux paths are converted to Windows paths via `wslpath` so Windows apps can read them directly.

![WSL demo](media/wsl_demo.gif)

## Setup on WSL

`peektea init` on WSL skips the Linux GUI app categories — there's no point offering `nautilus` or `feh` when you're routing through Windows — and sets the Windows opener as the fallback instead. See [Configuration](configuration.md) for the full setup walkthrough.
