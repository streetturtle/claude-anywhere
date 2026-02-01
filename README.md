# claude-anywhere

<p align="center">
<img width="723" height="438" alt="Screenshot 2026-01-31 at 10 24 10 PM" src="https://github.com/user-attachments/assets/ad6c8b44-84ff-40f1-b422-54d83578de04" />
</p>

**Resume any Claude Code session from anywhere.**

Never lose track of your Claude sessions again. `claude-anywhere` lets you see all your Claude Code sessions across all projects and instantly jump back into any of them - no matter which directory you're in.

## Quick Start

> **Tip**: Just ask Claude Code to install this for you! Share this repo URL with Claude and say: "Install claude-anywhere from https://github.com/streetturtle/claude-anywhere"

### Installation

```bash
brew install streetturtle/tap/claude-anywhere
```

## Prerequisites

Before using `claude-anywhere`, you need to set up a statusline script that captures session data.

### 1. Install `jq` if you don't have it

```bash
brew install jq
```

### 2. Create the session writer script

Create `~/.claude/write-status.sh` with the following content:

```bash
#!/bin/bash

# Read the JSON input from stdin
input=$(cat)

# Extract session ID
session_id=$(echo "$input" | jq -r '.session_id // "unknown"')

# Create sessions directory if it doesn't exist
sessions_dir="$HOME/.claude_sessions"
mkdir -p "$sessions_dir"

# Use session_id as the filename to support multiple sessions from the same folder
status_file="$sessions_dir/claude-status-${session_id}.json"

# Write the status data to file with timestamp
# Use milliseconds - on macOS, date doesn't support %N, so we append 000 to convert seconds to ms
echo "$input" | jq ". + {\"_statusline_update_time\": $(date +%s)000}" > "$status_file"

# Pass the input to stdout for the next script in the pipe
echo "$input"
```

Make it executable:

```bash
chmod +x ~/.claude/write-status.sh
```

### 3. Configure Claude Code statusline

Edit `~/.claude/settings.json` and add or update the `statusLine` configuration:

```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/write-status.sh | ~/.claude/your-existing-statusline-script.sh",
    "padding": 0
  }
}
```

If you don't have an existing statusline script, you can use just the writer:

```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/write-status.sh",
    "padding": 0
  }
}
```

### 4. Restart Claude Code

Restart Claude Code to apply the changes. Session data will now be written to `~/.claude_sessions/`.

## Usage

Simply run from any directory:

```bash
claude-anywhere
```

### Keyboard Controls

- **↑/↓** or **j/k** - Navigate through sessions
- **/** - Search/filter sessions
- **Enter** - Resume selected session
- **q** or **Ctrl+C** - Quit
- **Esc** - Clear search filter

## Features

- ⚡ **Resume from anywhere** - Jump into any Claude session from any directory
- 📊 **Interactive TUI** - Beautiful terminal UI with keyboard navigation and search
- 💰 **Cost tracking** - See real-time costs, token usage, and context window usage
- 🎨 **Status indicators** - Quickly identify active (●), idle (◐), and closed (○) sessions
- 🔍 **Fast search** - Filter sessions by project name or model

## What's Displayed

Each session shows:
- **Status indicator** - Active (●), Idle (◐), or Closed (○)
- **Session name** - Custom name if set (via `/rename` in Claude Code), otherwise project folder name
- **Model** - Which Claude model is being used
- **Cost** - Total cost in USD for the session
- **Context usage** - Percentage of context window used
- **Last activity** - How long ago the session was active
- **Token counts** - Input/output tokens used
- **Code changes** - Lines added/removed in the session

The header shows aggregate stats: total sessions, active count, idle count, and total cost across all sessions.

> **Tip**: Use `/rename <your-name>` in Claude Code to give your sessions meaningful names that will show up in `claude-anywhere`!

> **Note**: Session names may take a moment to update due to how Claude Code updates its session index. This is a known limitation in Claude Code itself.

## Troubleshooting

### No sessions showing up

1. Check that `~/.claude/settings.json` has the `statusLine` configuration
2. Ensure `~/.claude/write-status.sh` exists and is executable: `chmod +x ~/.claude/write-status.sh`
3. Verify `~/.claude_sessions/` exists and contains `claude-status-*.json` files
4. Make sure you have an active Claude Code session

### "jq: command not found" error

Install `jq`:
```bash
brew install jq
```

### Sessions show as "Closed" immediately

1. Verify your `settings.json` has the correct command with `write-status.sh`
2. Restart Claude Code
3. Check file timestamps: `ls -la ~/.claude_sessions/` to verify files are being updated

### Old sessions cluttering the list

By default, `claude-anywhere` shows sessions from the last 7 days. You can clean up old session files:

```bash
find ~/.claude_sessions -name "claude-status-*.json" -mtime +7 -delete
```
