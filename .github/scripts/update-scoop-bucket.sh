#!/usr/bin/env bash
# Generates the Scoop manifest for signls and pushes it to the emprcl/scoop-bucket
# repository. signls assembles its release manually across several runners (see
# .goreleaser.yaml), so goreleaser's own scoops pipe can't run; this reproduces the
# manifest it would have generated, using the sha256 sum from the combined
# checksums.txt once every platform's archive is published.
set -euo pipefail

VERSION="${VERSION#v}" # tag is vX.Y.Z; the manifest wants X.Y.Z
DIST="${DIST:-dist}"

BUCKET_OWNER="emprcl"
BUCKET_REPO="scoop-bucket"
BASE_URL="https://github.com/emprcl/signls/releases/download/v${VERSION}"

windows_amd64="$(awk -v f="signls_${VERSION}_windows_amd64.zip" '$2 == f { print $1 }' "${DIST}/checksums.txt")"

if [ -z "${windows_amd64}" ]; then
  echo "missing checksum for windows_amd64 in ${DIST}/checksums.txt" >&2
  exit 1
fi

cat > signls.json <<EOF
{
  "version": "${VERSION}",
  "description": "Non-linear, generative MIDI sequencer for the terminal",
  "homepage": "https://empr.cl/signls/",
  "license": "MIT",
  "architecture": {
    "64bit": {
      "url": "${BASE_URL}/signls_${VERSION}_windows_amd64.zip",
      "hash": "${windows_amd64}",
      "bin": "signls.exe"
    }
  }
}
EOF

git clone --depth 1 \
  "https://x-access-token:${HOMEBREW_TAP_GITHUB_TOKEN}@github.com/${BUCKET_OWNER}/${BUCKET_REPO}.git" \
  bucket

cp signls.json bucket/signls.json

cd bucket
git config user.name "goreleaserbot"
git config user.email "bot@goreleaser.com"

git add signls.json
if git diff --cached --quiet; then
  echo "signls.json already up to date for v${VERSION}; nothing to push."
  exit 0
fi

git commit -m "Scoop manifest update for signls version v${VERSION}"
git push
