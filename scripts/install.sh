#!/usr/bin/env sh
set -e

REPO="lovestaco/peektea"
BINARY="peektea"

# detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
    linux)  OS="linux"  ;;
    darwin) OS="darwin" ;;
    *)      echo "error: unsupported OS: $OS"; exit 1 ;;
esac

# detect arch
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)  ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *)             echo "error: unsupported arch: $ARCH"; exit 1 ;;
esac

# resolve latest version (e.g. v0.2.1)
VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | head -1 | cut -d'"' -f4)

if [ -z "$VERSION" ]; then
    echo "error: could not determine latest version (check your internet connection)"
    exit 1
fi

# GoReleaser strips the leading "v" from artifact names: tag v0.2.1 -> 0.2.1
VERSION_NUM=${VERSION#v}

FILENAME="${BINARY}_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

echo "installing peektea $VERSION ($OS/$ARCH)…"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

curl -fsSL "$URL" -o "$TMP/$FILENAME"
tar -xzf "$TMP/$FILENAME" -C "$TMP"

# pick install destination
if echo ":$PATH:" | grep -q ":$HOME/.local/bin:"; then
    DEST="$HOME/.local/bin"
elif [ -w "/usr/local/bin" ]; then
    DEST="/usr/local/bin"
else
    DEST="$HOME/.local/bin"
    mkdir -p "$DEST"
    echo "note: add $DEST to your PATH if peektea isn't found after install"
fi

install -m 755 "$TMP/$BINARY" "$DEST/$BINARY"
echo "installed → $DEST/$BINARY"
echo "run 'peektea -h' to get started"
