# Assets

This directory contains the session tracking script for `claude-anywhere`.

## File

**`write-status.sh`** - Session writer script that saves Claude Code session data to `~/.claude_sessions/`

## How It Works

When you configure Claude Code's statusline to pipe through `write-status.sh`, it:

1. Receives session data from Claude Code via stdin (JSON format)
2. Adds a timestamp to the data
3. Saves it to `~/.claude_sessions/claude-status-<project>.json`
4. Passes the data to stdout (so it can be piped to your existing statusline script)

This allows `claude-anywhere` to read your session information while keeping your existing statusline display intact.

## Installation

**Step 1: Set up Claude Code statusline (if not already enabled)**

Run this command in Claude Code to set up a statusline:

```
/status
```

**Step 2: Install the session writer**

```bash
cp write-status.sh ~/.claude/
chmod +x ~/.claude/write-status.sh
```

**Step 3: Update your statusline command**

Edit `~/.claude/settings.json` and modify the `statusLine.command` to pipe through `write-status.sh`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/write-status.sh | ~/.claude/your-existing-script.sh",
    "padding": 0
  }
}
```

For example, if your current command is `~/.claude/print-statusline.sh`, change it to:

```json
"command": "~/.claude/write-status.sh | ~/.claude/print-statusline.sh"
```

**Step 4: Restart Claude Code**

Your sessions will now be tracked in `~/.claude_sessions/` and you can use `claude-anywhere` to browse and resume them.

## Requirements

- `jq` - JSON processor
  - macOS: `brew install jq`
  - Linux: `apt-get install jq` or `yum install jq`
