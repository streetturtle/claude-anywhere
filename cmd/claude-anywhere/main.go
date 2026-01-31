package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/streetturtle/claude-anywhere/pkg/session"
)

var rootCmd = &cobra.Command{
	Use:   "claude-anywhere",
	Short: "Resume any Claude Code session from anywhere",
	Long: `claude-anywhere - Resume any Claude Code session from anywhere.

Never lose track of your Claude sessions again. Launch claude-anywhere from any directory to:
  • See all your active and recent Claude Code sessions
  • Jump back into any session with one keystroke
  • View costs, context usage, and session details at a glance`,
	RunE: listCommand,
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
	// Load all sessions
	sessions, err := session.LoadSessions()
	if err != nil {
		return fmt.Errorf("failed to load sessions: %w", err)
	}

	if len(sessions) == 0 {
		fmt.Println("No Claude sessions found.")
		fmt.Println("")
		fmt.Println("Start a Claude Code session and it will appear here automatically.")
		return nil
	}

	// Start the TUI
	return session.RunTUI(sessions)
}

func main() {
	rootCmd.AddCommand(listCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
