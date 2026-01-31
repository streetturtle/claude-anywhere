package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/streetturtle/clse/pkg/config"
	"github.com/streetturtle/clse/pkg/session"
)

//go:embed assets
var assetsFS embed.FS

var rootCmd = &cobra.Command{
	Use:   "clse",
	Short: "Resume any Claude Code session from anywhere",
	Long: `clse (CLaude SEssions) - Resume any Claude Code session from anywhere.

Never lose track of your Claude sessions again. Launch clse from any directory to:
  • See all your active and recent Claude Code sessions
  • Jump back into any session with one keystroke
  • View costs, context usage, and session details at a glance`,
	RunE: listCommand,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize statusline configuration for Claude Code",
	Long: `Installs the statusline script to ~/.claude/statusline.sh and configures
Claude Code to use it by updating ~/.claude/settings.json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read the embedded statusline script
		scriptContent, err := assetsFS.ReadFile("assets/statusline.sh")
		if err != nil {
			return fmt.Errorf("failed to read embedded statusline script: %w", err)
		}

		// Initialize the configuration
		return config.Initialize(scriptContent)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Browse and resume Claude Code sessions",
	Long: `Opens an interactive TUI to browse and resume your Claude Code sessions.

Search by project name or model, view session details, and press Enter to resume.
Works from any directory - no need to remember where each session was started.`,
	RunE: listCommand,
}

func listCommand(cmd *cobra.Command, args []string) error {
	// Check if statusline is initialized
	homeDir, err := os.UserHomeDir()
	if err == nil {
		settingsPath := homeDir + "/.claude/settings.json"
		if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
			fmt.Println("⚠️  Claude settings not found.")
			fmt.Println("")
			fmt.Println("It looks like you haven't initialized the statusline yet.")
			fmt.Println("Run this command first:")
			fmt.Println("")
			fmt.Println("  clse init")
			fmt.Println("")
			return nil
		}
	}

	// Load all sessions
	sessions, err := session.LoadSessions()
	if err != nil {
		return fmt.Errorf("failed to load sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No Claude sessions found.")
		fmt.Println("")
		fmt.Println("Make sure:")
		fmt.Println("  1. You have run 'clse init' to set up the statusline")
		fmt.Println("  2. You have restarted any Claude Code sessions")
		fmt.Println("  3. Claude Code is running in at least one terminal")
		fmt.Println("")
		fmt.Println("To verify the statusline is working:")
		fmt.Println("  ls -lah /tmp/claude-status*.json")
		return nil
	}

	// Start the TUI
	return session.RunTUI(sessions)
}

func main() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
