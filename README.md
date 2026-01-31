# clse - CLaude SEssions

**Resume any Claude Code session from anywhere.**

Never lose track of your Claude sessions again. `clse` lets you see all your Claude Code sessions across all projects and instantly jump back into any of them - no matter which directory you're in.

## Why clse?

Working on multiple projects with Claude Code? Switching between sessions? Lost track of where you started that important conversation?

**Just run `clse` from anywhere:**
- See all your active and recent Claude sessions in one place
- Search by project name or model
- Press Enter to resume any session
- View costs, context usage, and session stats

## Features

- ⚡ **Resume from anywhere** - Jump into any Claude session from any directory
- 📊 **Interactive TUI** - Beautiful terminal UI with keyboard navigation and search
- 💰 **Cost tracking** - See real-time costs, token usage, and context window usage
- 🎨 **Status indicators** - Quickly identify active, idle, and closed sessions
- 🔍 **Fast search** - Filter sessions by project name or model

## Installation

### From Source

```bash
git clone <repository-url>
cd cs
go build -o cs ./cmd/cs
sudo mv cs /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/streetsidesoftware/cs/cmd/cs@latest
```

## Setup

Before using `cs`, you need to initialize the statusline configuration:

```bash
cs init
```

This command will:
1. Install the statusline script to `~/.claude/statusline.sh`
2. Update your `~/.claude/settings.json` with the statusline configuration
3. Make the script executable

After running `cs init`, **restart any running Claude Code sessions** for the changes to take effect.
