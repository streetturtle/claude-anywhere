# claude-anywhere
<p align="center">
https://github.com/user-attachments/assets/f7cf75e6-0f76-4c69-8c79-ae2b76373396
</p>



**Resume any Claude Code session from anywhere.**

Never lose track of your Claude sessions again. `claude-anywhere` lets you see all your Claude Code sessions across all projects and instantly jump back into any of them - no matter which directory you're in.

## Why claude-anywhere?

Working on multiple projects with Claude Code? Switching between sessions? Lost track of where you started that important conversation?

**Just run `claude-anywhere` from anywhere:**
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

## Prerequisites

`claude-anywhere` requires Claude Code's statusline feature to be enabled. This feature writes session data to `~/.claude_sessions/` as JSON files.

### Setup Instructions

**Step 1: Enable Claude Code statusline (if not already enabled)**

If you don't have a statusline configured yet, run this in Claude Code:

```
/status
```

Follow the prompts to set up your statusline. This will create a statusline script in `~/.claude/` and configure it in your `settings.json`.

**Step 2: Install the session writer script**

Copy the session tracking script to your Claude directory:

```bash
# Clone this repo (or download write-status.sh from assets/)
git clone https://github.com/streetturtle/claude-anywhere.git
cd claude-anywhere

# Copy the script
cp assets/write-status.sh ~/.claude/
chmod +x ~/.claude/write-status.sh
```

**Step 3: Update your statusline command**

Edit `~/.claude/settings.json` and modify your existing `statusLine.command` to pipe through `write-status.sh`:

```json
{
  "statusLine": {
    "type": "command",
    "command": "~/.claude/write-status.sh | ~/.claude/your-existing-statusline-script.sh",
    "padding": 0
  }
}
```

For example, if your current statusline command is `~/.claude/print-statusline.sh`, change it to:

```json
"command": "~/.claude/write-status.sh | ~/.claude/print-statusline.sh"
```

The `write-status.sh` script will:
- Receive session data from Claude Code
- Save it to `~/.claude_sessions/claude-status-<project>.json`
- Pass the data through to your existing statusline script (so your terminal display stays the same)

**Step 4: Restart Claude Code**

Restart Claude Code to apply the changes. Session data will now be written to `~/.claude_sessions/`.

> **Note**: The scripts require `jq` to be installed: `brew install jq` (macOS) or `apt-get install jq` (Linux)

## Installation

### Homebrew (Recommended)

```bash
brew install streetturtle/tap/claude-anywhere
```

### From Source

```bash
go install github.com/streetturtle/claude-anywhere/cmd/claude-anywhere@latest
```

**Make sure `$GOPATH/bin` is in your PATH:**

```bash
# Add to ~/.bashrc, ~/.zshrc, or similar:
export PATH="$PATH:$(go env GOPATH)/bin"
```

After installation, the `claude-anywhere` command will be available globally from any directory.

### From Source

```bash
git clone https://github.com/streetturtle/claude-anywhere.git
cd claude-anywhere
make build
sudo mv claude-anywhere /usr/local/bin/
# Or use: make install
```

## Usage

Simply run `claude-anywhere` from any directory:

```bash
claude-anywhere
# Or explicitly: claude-anywhere list
```

The tool reads session data from `~/.claude_sessions/` (populated by the statusline feature) and displays all your sessions in an interactive terminal UI.

### Keyboard Controls

- **↑/↓** or **j/k** - Navigate through sessions
- **/** - Search/filter sessions by project name or model
- **Enter** - Resume selected session
- **q** or **Ctrl+C** - Quit
- **Esc** - Clear search filter

## Commands

```bash
claude-anywhere              # Launch interactive session browser (default)
claude-anywhere list         # Same as above
claude-anywhere completion   # Generate shell completion script
claude-anywhere help         # Show help
```

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

## Troubleshooting

### No sessions showing up

**Problem**: Running `claude-anywhere` shows an empty list or "No sessions found".

**Solutions**:
1. **Verify statusline is configured**: Check that `~/.claude/settings.json` contains a `statusLine` configuration
2. **Check the session writer script**: Ensure `~/.claude/write-status.sh` exists and is executable (`chmod +x ~/.claude/write-status.sh`)
3. **Verify the sessions directory**: Check that `~/.claude_sessions/` exists and contains `claude-status-*.json` files
4. **Start a Claude session**: The statusline only writes data when a Claude Code session is active. Start a session and interact with it, then try again

### "jq: command not found" error

**Problem**: The statusline script fails because `jq` is not installed.

**Solution**: Install `jq`:
```bash
# macOS
brew install jq

# Ubuntu/Debian
sudo apt-get install jq

# Fedora
sudo dnf install jq
```

### Permission denied errors

**Problem**: Scripts fail with permission errors.

**Solution**: Make the scripts executable:
```bash
chmod +x ~/.claude/write-status.sh
chmod +x ~/.claude/print-statusline.sh  # if you have this
```

### Sessions show as "Closed" immediately

**Problem**: All sessions appear as closed even when Claude is running.

**Solutions**:
1. **Check script pipeline**: Verify your `settings.json` has the correct command pipeline with `write-status.sh` first
2. **Restart Claude Code**: Changes to `settings.json` require restarting Claude Code
3. **Check file timestamps**: Run `ls -la ~/.claude_sessions/` to verify files are being updated

### Old sessions cluttering the list

**Problem**: Sessions from weeks ago are still showing up.

**Note**: By default, `claude-anywhere` shows sessions from the last 7 days. Older session files in `~/.claude_sessions/` can be safely deleted:
```bash
# Remove session files older than 7 days
find ~/.claude_sessions -name "claude-status-*.json" -mtime +7 -delete
```
