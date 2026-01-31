#!/bin/bash

# Claude Code session writer for claude-anywhere
# 
# This script receives session data from Claude Code via stdin,
# saves it to ~/.claude_sessions/, and passes it through to stdout
# so it can be piped to your existing statusline script.
#
# Installation:
# 1. Copy this file to ~/.claude/write-status.sh
# 2. Make it executable: chmod +x ~/.claude/write-status.sh
# 3. Update ~/.claude/settings.json to pipe through this script:
#    "statusLine": {
#      "type": "command",
#      "command": "~/.claude/write-status.sh | ~/.claude/your-statusline-script.sh",
#      "padding": 0
#    }
#
# Requirements: jq (install with: brew install jq)

# Read the JSON input from stdin
input=$(cat)

# Extract session ID and CWD
session_id=$(echo "$input" | jq -r '.session_id // "unknown"')
cwd=$(echo "$input" | jq -r '.cwd // ""')

# Create sessions directory if it doesn't exist
sessions_dir="$HOME/.claude_sessions"
mkdir -p "$sessions_dir"

# Create a sanitized filename based on CWD
if [ -n "$cwd" ]; then
  # Replace / with - and remove leading slash
  sanitized_cwd=$(echo "$cwd" | sed 's/^\//-/' | sed 's/\//-/g')
  status_file="$sessions_dir/claude-status${sanitized_cwd}.json"
else
  status_file="$sessions_dir/claude-status-${session_id}.json"
fi

# Write the status data to file with timestamp
# Use milliseconds - on macOS, date doesn't support %N, so we append 000 to convert seconds to ms
echo "$input" | jq ". + {\"_statusline_update_time\": $(date +%s)000}" > "$status_file"

# Pass the input to stdout for the next script in the pipe
echo "$input"
