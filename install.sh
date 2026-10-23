#!/usr/bin/env bash
set -euo pipefail

REPO="yuin/goldmark-jsonrenderer"
BINARY="markdown2json"
INSTALL_DIR="/usr/local/bin"
VERSION=""

usage() {
  cat <<'EOF'
Usage: install.sh [-b INSTALL_DIR] [-v VERSION]

Options:
  -b DIR       Installation directory (default: /usr/local/bin)
  -v VERSION   Version to install (e.g. v1.2.3)
  -h           Show this help
EOF
}

while getopts ":b:v:h" opt; do
  case "$opt" in
    b) INSTALL_DIR="$OPTARG" ;;
    v) VERSION="$OPTARG" ;;
    h) usage; exit 0 ;;
    :)
      echo "Error: -$OPTARG requires an argument" >&2
      exit 2
      ;;
    \?)
      echo "Error: unknown option -$OPTARG" >&2
      usage >&2
      exit 2
      ;;
  esac
done

for cmd in curl tar install uname; do
  command -v "$cmd" >/dev/null 2>&1 || {
    echo "Error: required command not found: $cmd" >&2
    exit 1
  }
done

OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Linux)  GOOS="linux" ;;
  Darwin) GOOS="darwin" ;;
  *)
    echo "Error: unsupported OS: $OS" >&2
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64 | amd64) GOARCH="amd64" ;;
  aarch64 | arm64) GOARCH="arm64" ;;
  *)
    echo "Error: unsupported architecture: $ARCH" >&2
    exit 1
    ;;
esac

if [[ -z "$VERSION" ]]; then
  VERSION="$(
    curl -fsSL \
      "https://api.github.com/repos/${REPO}/releases/latest" |
      sed -n 's/.*"tag_name": "\([^"]*\)".*/\1/p' |
      head -n 1
  )"
fi

if [[ -z "$VERSION" ]]; then
  echo "Error: could not determine release version" >&2
  exit 1
fi

ARCHIVE="${BINARY}_${VERSION}_${GOOS}_${GOARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading ${BINARY} ${VERSION} for ${GOOS}/${GOARCH}..."
curl -fL --retry 3 "$URL" -o "${TMP_DIR}/${ARCHIVE}"

tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "$TMP_DIR"

if [[ ! -f "${TMP_DIR}/${BINARY}" ]]; then
  echo "Error: binary not found in archive: ${BINARY}" >&2
  exit 1
fi

mkdir -p "$INSTALL_DIR"
install -m 0755 "${TMP_DIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"

echo "Installed ${BINARY} ${VERSION} to ${INSTALL_DIR}/${BINARY}"
