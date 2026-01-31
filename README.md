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

## Usage

### List Sessions

Run `cs` or `cs list` to open the interactive session browser:

```bash
cs
# or
cs list
```

#### Navigation

- **↑/↓ or j/k** - Navigate up/down
- **/** - Start searching/filtering
- **Enter** - Resume selected session (changes to the project directory and runs `claude -r <session-id>`)
- **Esc** - Exit search mode
- **q or Ctrl+C** - Quit

#### Session Information

Each session displays:
- **Status** - Active (●), Idle (◐), or Closed (○)
- **Project Name** - The name of the project/directory
- **Model** - The Claude model being used (e.g., Sonnet 4.5)
- **Cost** - Total cost in USD
- **Context Usage** - Percentage of context window used
- **Last Updated** - Relative time since last activity

### Example Output

```
Claude Sessions

● project-name                           Active
  Model: Sonnet 4.5 • Cost: $0.0234 • Context: 45% • Updated: 2m ago

◐ another-project                        Idle
  Model: Sonnet 4.5 • Cost: $0.0156 • Context: 23% • Updated: 15m ago

○ old-project                            Closed
  Model: Opus 4 • Cost: $0.1234 • Context: 67% • Updated: 2h ago
```

## How It Works

`cs` integrates with Claude Code's statusline feature:

1. **Statusline Script**: Claude Code pipes session status as JSON to `~/.claude/statusline.sh`
2. **Temp Files**: The script writes formatted data to `/tmp/claude-status-*.json` files
3. **cs CLI**: Reads these temp files to display current session information

The statusline updates every 300ms while Claude is active, providing real-time status.

## Architecture

```
cs/
├── cmd/cs/              # Main CLI entry point
│   ├── main.go         # Cobra commands and CLI structure
│   └── assets/         # Embedded statusline.sh script
├── pkg/
│   ├── config/         # Configuration and initialization
│   │   └── config.go
│   └── session/        # Session management and TUI
│       ├── session.go  # Session parsing and data structures
│       └── tui.go      # Bubbletea TUI implementation
└── assets/
    └── statusline.sh   # Statusline script (embedded in binary)
```

## Requirements

- Go 1.22+ (for building)
- Claude Code CLI
- macOS or Linux (uses `/tmp` for status files)

## Dependencies

- [cobra](https://github.com/spf13/cobra) - CLI framework
- [bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal styling

## Session Status

- **Active** (●) - Claude is currently processing (updated < 3 seconds ago)
- **Idle** (◐) - Claude session exists but is waiting (updated < 1 hour ago)
- **Closed** (○) - Session is old or inactive (updated > 1 hour ago)

### Why Sessions Might Disappear

Sessions may disappear from the list for these reasons:

1. **Empty status files** - If a Claude Code session ends abnormally, its status file in `/tmp` may become empty (0 bytes), causing it to be skipped during parsing.
2. **Status file deleted** - Status files in `/tmp` may be cleaned up by the system or manually deleted.
3. **Parsing errors** - Corrupted or malformed JSON in status files will cause them to be skipped.

To see all status files including empty/corrupted ones:
```bash
ls -lah /tmp/claude-status*.json
```

## Troubleshooting

### No sessions found

1. Verify statusline is installed: `ls -la ~/.claude/statusline.sh`
2. Check settings.json: `cat ~/.claude/settings.json`
3. Restart Claude Code sessions
4. Verify status files exist: `ls -la /tmp/claude-status*.json`

### Statusline not updating

1. Check if the script is executable: `chmod +x ~/.claude/statusline.sh`
2. Verify the script path in settings.json matches the installed location
3. Make sure you restarted Claude Code after running `cs init`

### Sessions not showing accurate data

The status files in `/tmp` may be stale. Claude Code updates these files every 300ms while active. If you don't see updates, the Claude session might be closed or the statusline isn't configured correctly.

## License

MIT

## Credits

Based on the [Claude Session Monitor Raycast Extension](https://github.com/raycast/extensions) by Pavel Makhov.
