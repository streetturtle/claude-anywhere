package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type StatusLineConfig struct {
	Type    string `json:"type"`
	Command string `json:"command"`
	Padding int    `json:"padding"`
}

type Settings struct {
	StatusLine *StatusLineConfig `json:"statusLine,omitempty"`
}

// Initialize sets up the Claude Code statusline configuration
func Initialize(scriptContent []byte) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	claudeDir := filepath.Join(homeDir, ".claude")
	settingsPath := filepath.Join(claudeDir, "settings.json")
	defaultScriptPath := filepath.Join(claudeDir, "statusline.sh")

	// Create .claude directory if it doesn't exist
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	// Read existing settings if file exists
	var settings Settings

	if _, err := os.Stat(settingsPath); err == nil {
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			return fmt.Errorf("failed to read existing settings: %w", err)
		}

		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("failed to parse existing settings: %w", err)
		}
	}

	// Check if statusline is already configured
	if settings.StatusLine != nil {
		fmt.Println("⚠️  Statusline already configured")
		fmt.Println("")
		fmt.Printf("Your current statusline script: %s\n", settings.StatusLine.Command)
		fmt.Println("")
		fmt.Println("To use clse, you need to add monitoring code to your existing statusline script.")
		fmt.Println("")
		fmt.Println("📖 See detailed instructions:")
		fmt.Println("   https://github.com/streetturtle/clse#existing-statusline")
		fmt.Println("")
		fmt.Println("Quick start - ask Claude Code:")
		fmt.Println("")
		fmt.Println("  \"Please add this code to my statusline script to enable clse monitoring:\"")
		fmt.Println("")
		fmt.Println("  # Extract session ID and CWD")
		fmt.Println("  session_id=$(echo \"$input\" | jq -r '.session_id // \"unknown\"')")
		fmt.Println("  cwd=$(echo \"$input\" | jq -r '.cwd // \"\"')")
		fmt.Println("")
		fmt.Println("  # Create a sanitized filename based on CWD")
		fmt.Println("  if [ -n \"$cwd\" ]; then")
		fmt.Println("    sanitized_cwd=$(echo \"$cwd\" | sed 's/^\\//-/' | sed 's/\\//-/g')")
		fmt.Println("    status_file=\"/tmp/claude-status${sanitized_cwd}.json\"")
		fmt.Println("  else")
		fmt.Println("    status_file=\"/tmp/claude-status-${session_id}.json\"")
		fmt.Println("  fi")
		fmt.Println("")
		fmt.Println("  # Write the status data to temp file with timestamp")
		fmt.Println("  echo \"$input\" | jq \". + {\\\"_statusline_update_time\\\": $(date +%s)000}\" > \"$status_file\"")
		fmt.Println("")
		return nil
	}

	// No statusline configured - install our full script
	fmt.Println("No statusline configured. Installing clse statusline...")
	fmt.Println("")

	// Write our full statusline script
	if err := os.WriteFile(defaultScriptPath, scriptContent, 0755); err != nil {
		// Check if it's a permission error
		if os.IsPermission(err) {
			fmt.Printf("⚠️  Permission denied writing to %s\n", defaultScriptPath)
			fmt.Println("")
			fmt.Println("Please run this command to install manually:")
			fmt.Printf("  sudo install -m 755 /dev/stdin %s <<'EOF'\n", defaultScriptPath)
			fmt.Print(string(scriptContent))
			fmt.Println("EOF")
			return fmt.Errorf("permission denied")
		}
		return fmt.Errorf("failed to write statusline script: %w", err)
	}
	fmt.Printf("✓ Installed statusline script to %s\n", defaultScriptPath)

	// Add statusline configuration to settings
	settings.StatusLine = &StatusLineConfig{
		Type:    "command",
		Command: defaultScriptPath,
		Padding: 0,
	}

	// Write updated settings
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	fmt.Printf("✓ Updated settings at %s\n", settingsPath)
	fmt.Println("")
	fmt.Println("✅ Setup complete!")
	fmt.Println("")
	fmt.Println("Next steps:")
	fmt.Println("  1. Restart any running Claude Code sessions")
	fmt.Println("  2. Run 'clse' to view your sessions")
	fmt.Println("")

	return nil
}
