package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type SessionStatus string

const (
	StatusActive SessionStatus = "active"
	StatusIdle   SessionStatus = "idle"
	StatusClosed SessionStatus = "closed"
)

const (
	ActivityThresholdMs = 3000      // 3 seconds
	ClosedThresholdMs   = 3600000   // 1 hour
	MaxAgeMs            = 604800000 // 7 days - sessions older than this are filtered by default
)

// Session represents a Claude Code session with all its metadata
type Session struct {
	SessionID    string        `json:"session_id"`
	SessionName  string        `json:"session_name"`
	CWD          string        `json:"cwd"`
	ProjectDir   string        `json:"project_dir"`
	ProjectName  string        `json:"project_name"`
	Status       SessionStatus `json:"status"`
	LastActivity time.Time     `json:"last_activity"`

	// Model info
	Model string `json:"model"`

	// Cost and usage
	CostUSD    float64 `json:"cost_usd"`
	DurationMs int64   `json:"duration_ms"`

	// Context window
	ContextUsedPct float64 `json:"context_used_percent"`
	InputTokens    int     `json:"input_tokens"`
	OutputTokens   int     `json:"output_tokens"`
	TotalTokens    int     `json:"total_tokens"`

	// Code changes
	LinesAdded   int `json:"lines_added"`
	LinesRemoved int `json:"lines_removed"`

	// Internal
	StatusFilePath string `json:"-"`
}

// StatusLineData represents the JSON data written by statusline.sh
type StatusLineData struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	UpdateTime     int64  `json:"_statusline_update_time"`

	Model struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"model"`

	Workspace struct {
		ProjectDir string `json:"project_dir"`
	} `json:"workspace"`

	Cost struct {
		TotalCostUSD      float64 `json:"total_cost_usd"`
		TotalDurationMs   int64   `json:"total_duration_ms"`
		TotalLinesAdded   int     `json:"total_lines_added"`
		TotalLinesRemoved int     `json:"total_lines_removed"`
	} `json:"cost"`

	ContextWindow struct {
		UsedPercentage    float64 `json:"used_percentage"`
		TotalInputTokens  int     `json:"total_input_tokens"`
		TotalOutputTokens int     `json:"total_output_tokens"`
	} `json:"context_window"`
}

// SessionsIndex represents the structure of sessions-index.json in Claude project directories
type SessionsIndex struct {
	Entries []struct {
		SessionID string `json:"sessionId"`
		Summary   string `json:"summary"`
	} `json:"entries"`
}

// LoadSessions reads all claude-status*.json files from ~/.claude_sessions and parses them
func LoadSessions() ([]*Session, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	sessionsDir := filepath.Join(homeDir, ".claude_sessions")

	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		// If directory doesn't exist, return empty sessions list
		if os.IsNotExist(err) {
			return []*Session{}, nil
		}
		return nil, fmt.Errorf("failed to read sessions directory: %w", err)
	}

	var sessions []*Session

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "claude-status") || !strings.HasSuffix(name, ".json") {
			continue
		}

		filePath := filepath.Join(sessionsDir, name)
		session, err := parseSessionFile(filePath)
		if err != nil {
			// Skip files that can't be parsed
			continue
		}

		sessions = append(sessions, session)
	}

	// Sort by last activity (most recent first)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].LastActivity.After(sessions[j].LastActivity)
	})

	return sessions, nil
}

func parseSessionFile(filePath string) (*Session, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("empty file")
	}

	var statusData StatusLineData
	if err := json.Unmarshal(data, &statusData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Skip if no update time
	if statusData.UpdateTime == 0 {
		return nil, fmt.Errorf("no update time")
	}

	// Determine project name from CWD
	projectName := filepath.Base(statusData.CWD)
	if projectName == "" || projectName == "/" {
		projectName = statusData.CWD
	}

	// Determine project directory
	projectDir := statusData.Workspace.ProjectDir
	if projectDir == "" {
		projectDir = statusData.CWD
	}

	// Calculate time since last update
	timeSinceUpdate := time.Since(time.UnixMilli(statusData.UpdateTime))

	// Determine status
	var status SessionStatus
	if timeSinceUpdate < time.Duration(ActivityThresholdMs)*time.Millisecond {
		status = StatusActive
	} else if timeSinceUpdate < time.Duration(ClosedThresholdMs)*time.Millisecond {
		status = StatusIdle
	} else {
		status = StatusClosed
	}

	// Get model name
	model := statusData.Model.DisplayName
	if model == "" {
		model = statusData.Model.ID
	}
	if model == "" {
		model = "Unknown"
	}

	session := &Session{
		SessionID:      statusData.SessionID,
		CWD:            statusData.CWD,
		ProjectDir:     projectDir,
		ProjectName:    projectName,
		Status:         status,
		LastActivity:   time.UnixMilli(statusData.UpdateTime),
		Model:          model,
		CostUSD:        statusData.Cost.TotalCostUSD,
		DurationMs:     statusData.Cost.TotalDurationMs,
		ContextUsedPct: statusData.ContextWindow.UsedPercentage,
		InputTokens:    statusData.ContextWindow.TotalInputTokens,
		OutputTokens:   statusData.ContextWindow.TotalOutputTokens,
		TotalTokens:    statusData.ContextWindow.TotalInputTokens + statusData.ContextWindow.TotalOutputTokens,
		LinesAdded:     statusData.Cost.TotalLinesAdded,
		LinesRemoved:   statusData.Cost.TotalLinesRemoved,
		StatusFilePath: filePath,
	}

	// Try to get session name from sessions-index.json
	session.SessionName = getSessionName(statusData.TranscriptPath, statusData.SessionID)

	return session, nil
}

// getSessionName looks up the session name from sessions-index.json
func getSessionName(transcriptPath, sessionID string) string {
	// If transcript path is empty, can't look up the name
	if transcriptPath == "" {
		return ""
	}

	// Get the project directory (parent of transcript file)
	projectDir := filepath.Dir(transcriptPath)
	indexPath := filepath.Join(projectDir, "sessions-index.json")

	// Read sessions-index.json
	data, err := os.ReadFile(indexPath)
	if err != nil {
		// Silently fail - file might not exist or be readable
		return ""
	}

	var index SessionsIndex
	if err := json.Unmarshal(data, &index); err != nil {
		// Silently fail - file might be malformed
		return ""
	}

	// Find matching session entry
	for _, entry := range index.Entries {
		if entry.SessionID == sessionID {
			return entry.Summary
		}
	}

	// Session not found in index
	return ""
}

// GetStatusIcon returns a colored status indicator for display
func (s *Session) GetStatusIcon() string {
	switch s.Status {
	case StatusActive:
		return "●" // Green circle
	case StatusIdle:
		return "◐" // Yellow half-circle
	case StatusClosed:
		return "○" // Gray circle
	default:
		return "?"
	}
}

// GetStatusText returns the status as a human-readable string
func (s *Session) GetStatusText() string {
	switch s.Status {
	case StatusActive:
		return "Active"
	case StatusIdle:
		return "Idle"
	case StatusClosed:
		return "Closed"
	default:
		return "Unknown"
	}
}

// GetDurationString returns a human-readable duration string
func (s *Session) GetDurationString() string {
	minutes := s.DurationMs / 1000 / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%dh %dm", hours, mins)
}

// GetLastActivityString returns a human-readable relative time
func (s *Session) GetLastActivityString() string {
	duration := time.Since(s.LastActivity)

	if duration < time.Minute {
		return "just now"
	} else if duration < time.Hour {
		mins := int(duration.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%dh ago", hours)
	} else {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}

// GetResumeCommand returns the command to resume this session
func (s *Session) GetResumeCommand() string {
	return fmt.Sprintf("cd \"%s\" && claude -r %s", s.CWD, s.SessionID)
}

// FilterRecentSessions filters out sessions older than 7 days
func FilterRecentSessions(sessions []*Session) []*Session {
	var recent []*Session
	cutoff := time.Now().Add(-time.Duration(MaxAgeMs) * time.Millisecond)

	for _, s := range sessions {
		if s.LastActivity.After(cutoff) {
			recent = append(recent, s)
		}
	}

	return recent
}
