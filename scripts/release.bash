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

    mkdir -p go-irc-bot-"$version"/src/go-irc-bot
    
    rsync -avzr --delete \
	  --filter='- go-irc-bot-*' \
	  --filter='- .*' \
	  . go-irc-bot-"$version"/src/go-irc-bot
    
    tar czf go-irc-bot-"$version".tgz go-irc-bot-"$version"
    git checkout version/version.go
else
    echo "Usage: $0 version <num>" 1>&2
    exit 1
fi
