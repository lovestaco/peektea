#!/usr/bin/env python3
"""
Render release notes and publish a GitHub release.

IMG:media/file.gif  →  https://raw.githubusercontent.com/lovestaco/peektea/refs/heads/master/media/file.gif

Usage:
    python3 scripts/release_gh.py <version>          # publish
    python3 scripts/release_gh.py <version> --render # just print rendered notes, don't publish
"""

import re
import sys
import subprocess
import tempfile
import os

REPO   = "lovestaco/peektea"
BRANCH = "master"
BASE   = f"https://raw.githubusercontent.com/{REPO}/refs/heads/{BRANCH}"

IMG_RE = re.compile(r"IMG:([^\s)\"]+)")


def render(version: str) -> str:
    path = f"docs/releases/{version}.md"
    if not os.path.exists(path):
        sys.exit(f"error: {path} not found — run 'make release-notes-init VERSION={version}' first")
    text = open(path).read()
    return IMG_RE.sub(lambda m: f"{BASE}/{m.group(1)}", text)


def publish(version: str, notes: str) -> None:
    # check tag exists
    result = subprocess.run(["git", "tag", "-l", version], capture_output=True, text=True)
    if version not in result.stdout.split():
        sys.exit(f"error: tag {version} not found — create it with 'make bump'")

    with tempfile.NamedTemporaryFile("w", suffix=".md", delete=False) as f:
        f.write(notes)
        tmp = f.name

    try:
        subprocess.run(
            ["gh", "release", "create", version,
             "--title", f"peektea {version}",
             "--notes-file", tmp],
            check=True,
        )
        print(f"released {version}")
    finally:
        os.unlink(tmp)


if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.exit("usage: release_gh.py <version> [--render]")

    version  = sys.argv[1]
    rendered = render(version)

    if "--render" in sys.argv:
        print(rendered)
    else:
        publish(version, rendered)
