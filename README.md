go-irc-bot
--------------------------------------------------------------

## About
Just and IRC bot for convenient collective workflow on projects.
### Features
- HTTP API to post messages to channels
- Command-line like interface for IRC

## Extending
Treat this project as a boilerplate. Fork, create your branch, extend it.

## Clone
We assume that all your golang root located at `~/go` directory.
```bash
git clone git-server:u/aphex/go-irc-bot.git go/src/go-irc-bot
```

## How to build
Again, we assume that all your golang root located at `~/go` directory.
If your golang root is in different place, just setup `GOPATH`, `PATH` variables manualy before doing anything specified in this chapter.
This project have a script to setup environment and bootstrap the dependencies, all you need to do is just source it like this and make:
```bash
source go-irc-bot/bootstrap
make -C go-irc-bot
```

## Packaging information
All packaging things should be placed in a subdirectory of `packaging/`.
At this time there are:
- RPM spec

## Daemonizing
All daemonizing things should be placed in a subdirectory of `startup/`.
At this time there are:
- SystemD support(unit+tmpfiles.d)
