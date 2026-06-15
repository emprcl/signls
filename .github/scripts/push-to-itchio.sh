#!/usr/bin/env bash
# Publishes the release archives to itch.io with butler, one channel per
# platform/arch. butler authenticates via the BUTLER_API_KEY env var. The
# archives are extracted first so itch hosts the runnable files (binary +
# LICENSE + README) rather than an opaque tarball; itch infers the OS from the
# "linux" / "windows" / "osx" substring in each channel name.
set -euo pipefail

VERSION="${VERSION#v}" # tag is vX.Y.Z; itch wants the bare version
DIST="${DIST:-dist}"
ITCH_TARGET="emprcl/signls"

# "channel archive" pairs; itch infers the OS from the channel name substring.
targets="
linux-amd64 signls_${VERSION}_linux_amd64.tar.gz
linux-arm64 signls_${VERSION}_linux_arm64.tar.gz
windows-amd64 signls_${VERSION}_windows_amd64.zip
osx-amd64 signls_${VERSION}_darwin_amd64.tar.gz
osx-arm64 signls_${VERSION}_darwin_arm64.tar.gz
"

# Install butler (itch.io's upload CLI).
curl -sL -o butler.zip https://broth.itch.zone/butler/linux-amd64/LATEST/archive/default
unzip -o butler.zip
chmod +x butler
./butler -V

while read -r channel filename; do
  [ -z "$channel" ] && continue
  archive="${DIST}/${filename}"
  if [ ! -f "$archive" ]; then
    echo "missing archive ${archive} for channel ${channel}" >&2
    exit 1
  fi

  outdir="itch/${channel}"
  rm -rf "$outdir"
  mkdir -p "$outdir"
  case "$archive" in
    *.zip) unzip -o "$archive" -d "$outdir" ;;
    *.tar.gz) tar -xzf "$archive" -C "$outdir" ;;
  esac

  ./butler push "$outdir" "${ITCH_TARGET}:${channel}" --userversion "${VERSION}"
done <<EOF
$targets
EOF
