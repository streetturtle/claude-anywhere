#!/bin/bash

# Claude Code statusline script for session monitoring
# This script receives real-time status updates from Claude Code via stdin
# and writes them to /tmp for the Raycast extension to read

# Read the JSON input from stdin
input=$(cat)

# Extract session ID and CWD
session_id=$(echo "$input" | jq -r '.session_id // "unknown"')
cwd=$(echo "$input" | jq -r '.cwd // ""')

# Create a sanitized filename based on CWD (so we can match by pane's working directory)
if [ -n "$cwd" ]; then
  # Replace / with - and remove leading slash
  sanitized_cwd=$(echo "$cwd" | sed 's/^\//-/' | sed 's/\//-/g')
  status_file="/tmp/claude-status${sanitized_cwd}.json"
else
  status_file="/tmp/claude-status-${session_id}.json"
fi

# Write the status data to temp file with timestamp
# Use milliseconds - on macOS, date doesn't support %N, so we append 000 to convert seconds to ms
echo "$input" | jq ". + {\"_statusline_update_time\": $(date +%s)000}" > "$status_file"

# Output a simple statusline for the terminal (optional)
# Extract useful info to display
model=$(echo "$input" | jq -r '.model.display_name // "Claude"')
context_used=$(echo "$input" | jq -r '.context_window.used_percentage // 0' | awk '{printf "%.0f", $1}')

echo "[$model] Context: ${context_used}%"
