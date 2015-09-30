#!/usr/bin/env bash

if [ "$1" = "version" ] && [ ! -z "$2" ]; then
    perl -p -i -e "s|<<version>>|$2|" version/version.go
    perl -p -i -e "s|<<rev>>|$(git rev-parse HEAD)|" version/version.go
else
    echo "Usage: $0 version <num>" 1>&2
    exit 1
fi
