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

## AI Agent Features

### JSON Output (`-j`)
All tools support structured JSON output for programmatic consumption:

```bash
git-standup -j              # Machine-readable commit history
dupes -j /path              # Structured duplicate report
watch -j '*.go' go test     # JSON events with output + exit codes
```

### Quiet Mode (`-q`)
Suppress human-readable output. Useful in automated workflows:

```bash
git-standup -jq             # JSON only, no headers
dupes -jq /path             # JSON only, no progress
```

### Environment Variables
- `STANDUP_DAYS` — Default days back for git-standup

### Webhook Support (watch)
POST events to a URL when files change:

```bash
watch -w http://localhost:8080/hook '*.go' go build
```

### Exit Codes
- `0` — Success / no issues found
- `1` — Error or duplicates found (dupes)
