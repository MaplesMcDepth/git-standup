# git-standup

![CI](https://github.com/MaplesMcDepth/git-standup/actions/workflows/ci.yml/badge.svg)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)


What did you do yesterday? A simple git commit summary tool.

## Install

```bash
go install github.com/MaplesMcDepth/git-standup/cmd/git-standup@latest
```

## Commands

```bash
git-standup              # Yesterday's commits
git-standup -d 3         # Last 3 days
git-standup -a "Bark"    # Only commits by author
git-standup -s           # Short format
git-standup -r           # Include all branches
```

## Options

| Flag | Description |
|------|-------------|
| `-d int` | Days back (default 1) |
| `-a string` | Author filter |
| `-b string` | Branch (default current) |
| `-r` | Include remote branches |
| `-s` | Short format |
