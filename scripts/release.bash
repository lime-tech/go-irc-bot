#!/usr/bin/env bash
set -e

source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/bootstrap"

if [ "$1" = "version" ] && [ ! -z "$2" ]; then
    git tag "$2"
    perl -p -i -e "s|<<version>>|$2|" version/version.go
    perl -p -i -e "s|<<rev>>|$(git rev-parse HEAD)|" version/version.go
    tar czf "$root"/irc-workflow-"$2".tgz .
else
    echo "Usage: $0 version <num>" 1>&2
    exit 1
fi
