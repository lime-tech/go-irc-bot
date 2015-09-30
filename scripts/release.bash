#!/usr/bin/env bash
set -e

version="$2"

git clean -Xdff
git clean -xdff

if [ "$1" = "version" ] && [ ! -z "$version" ]; then
    source "$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/bootstrap"
    
    git tag "$version"

    perl -p -i -e "s|<<version>>|$version|" version/version.go
    perl -p -i -e "s|<<rev>>|$(git rev-parse HEAD)|" version/version.go

    mkdir -p irc-workflow-"$version"/src/irc-workflow
    
    rsync -avzr --delete \
	  --filter='- irc-workflow-*' \
	  --filter='- .*' \
	  . irc-workflow-"$version"/src/irc-workflow
    
    tar czf irc-workflow-"$version".tgz irc-workflow-"$version"
    git checkout version/version.go
else
    echo "Usage: $0 version <num>" 1>&2
    exit 1
fi
