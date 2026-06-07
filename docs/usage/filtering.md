# Filtering

![filter and hidden files demo](../media/hide_files-and-filter_demo.gif)

Press `/` to enter filter mode. Type anything and the list narrows to matching entries in real time. The filter input appears at the bottom of the panel, above the hint bar — like vim's command line.

| Key       | Action                                  |
| --------- | --------------------------------------- |
| `enter`   | confirm and exit filter mode (filter stays active) |
| `esc`     | clear the filter entirely               |
| `↑` / `↓` | still navigate the list while you type  |

## Hidden files

Press `.` to toggle hidden files (dotfiles) on and off. The hint bar always shows the current state — `. show hidden` or `. hide hidden`.

Both filters compose: you can search by name with dotfiles visible *or* hidden at the same time.
