#!/bin/bash

set -e

RELEASE="v$1"
CURRENT="$(git tag | tail -n1)"

if [ -z "$1" ]; then
  echo "You must pass a version. Example release.sh 1.0"
  exit 1
fi

echo "Installing needed tools..."
go get github.com/aktau/github-release

echo "Creating Release"
mkdir -p release
go build podfetcher.go
mv podfetcher release

mv release podfetcher

tar cfav podfetcher-linux-amd64.tar.bz2 podfetcher
rm -rf podfetcher

LOG="$(git log --pretty=oneline --abbrev-commit "$CURRENT"..HEAD)"
git tag "$RELEASE"
git push origin "$RELEASE"
github-release release \
--user gregf \
--repo podfetcher \
--tag "$RELEASE" \
--name "podfetcher" \
--description "$LOG" \

github-release upload \
--user gregf \
--repo podfetcher \
--tag "$RELEASE" \
--name "podfetcher-linux-amd64.tar.bz2" \
--file "podfetcher-linux-amd64.tar.bz2"

rm -rf release
rm -f podfetcher-linux-amd64.tar.bz2
