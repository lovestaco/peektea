"""Rewrite IMG:media/file.gif placeholders in release notes to relative paths.

release-notes are shared with scripts/release_gh.py, which rewrites the same
IMG:media/file.gif markers to raw GitHub URLs when publishing a GitHub release.
For the docs site we instead point them at docs/media (symlinked to ../media)
so the same source files render correctly in both places.
"""

import re

IMG_RE = re.compile(r"IMG:(media/[^\s)\"]+)")


def on_page_markdown(markdown, page, config, files):
    if not page.file.src_uri.startswith("releases/"):
        return markdown

    depth = page.file.src_uri.count("/")
    prefix = "../" * depth
    return IMG_RE.sub(lambda m: f"{prefix}{m.group(1)}", markdown)
