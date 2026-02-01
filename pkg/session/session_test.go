package session

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseSessionFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		content     string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, s *Session)
	}{
		{
			name:        "empty file",
			content:     "",
			wantErr:     true,
			errContains: "empty file",
		},
		{
			name:        "invalid JSON",
			content:     "not valid json{",
			wantErr:     true,
			errContains: "failed to parse JSON",
		},
		{
			name:        "missing update time",
			content:     `{"session_id": "abc123", "cwd": "/tmp"}`,
			wantErr:     true,
			errContains: "no update time",
		},
		{
			name: "valid session data",
			content: func() string {
				data := StatusLineData{
					SessionID:  "session-123",
					CWD:        "/home/user/project",
					UpdateTime: time.Now().UnixMilli(),
				}
				data.Model.DisplayName = "Claude Sonnet"
				data.Model.ID = "claude-sonnet-4"
				data.Workspace.ProjectDir = "/home/user/project"
				data.Cost.TotalCostUSD = 0.0123
				data.Cost.TotalDurationMs = 120000
				data.Cost.TotalLinesAdded = 50
				data.Cost.TotalLinesRemoved = 10
				data.ContextWindow.UsedPercentage = 45.5
				data.ContextWindow.TotalInputTokens = 1000
				data.ContextWindow.TotalOutputTokens = 500
				b, _ := json.Marshal(data)
				return string(b)
			}(),
			wantErr: false,
			validate: func(t *testing.T, s *Session) {
				if s.SessionID != "session-123" {
					t.Errorf("SessionID = %q, want %q", s.SessionID, "session-123")
				}
				if s.CWD != "/home/user/project" {
					t.Errorf("CWD = %q, want %q", s.CWD, "/home/user/project")
				}
				if s.Model != "Claude Sonnet" {
					t.Errorf("Model = %q, want %q", s.Model, "Claude Sonnet")
				}
				if s.CostUSD != 0.0123 {
					t.Errorf("CostUSD = %f, want %f", s.CostUSD, 0.0123)
				}
				if s.ContextUsedPct != 45.5 {
					t.Errorf("ContextUsedPct = %f, want %f", s.ContextUsedPct, 45.5)
				}
				if s.InputTokens != 1000 {
					t.Errorf("InputTokens = %d, want %d", s.InputTokens, 1000)
				}
				if s.OutputTokens != 500 {
					t.Errorf("OutputTokens = %d, want %d", s.OutputTokens, 500)
				}
				if s.TotalTokens != 1500 {
					t.Errorf("TotalTokens = %d, want %d", s.TotalTokens, 1500)
				}
				if s.LinesAdded != 50 {
					t.Errorf("LinesAdded = %d, want %d", s.LinesAdded, 50)
				}
				if s.LinesRemoved != 10 {
					t.Errorf("LinesRemoved = %d, want %d", s.LinesRemoved, 10)
				}
				if s.Status != StatusActive {
					t.Errorf("Status = %q, want %q", s.Status, StatusActive)
				}
			},
		},
		{
			name: "model falls back to ID when display name empty",
			content: func() string {
				data := StatusLineData{
					SessionID:  "session-456",
					CWD:        "/tmp/test",
					UpdateTime: time.Now().UnixMilli(),
				}
				data.Model.ID = "claude-opus-4"
				data.Model.DisplayName = ""
				b, _ := json.Marshal(data)
				return string(b)
			}(),
			wantErr: false,
			validate: func(t *testing.T, s *Session) {
				if s.Model != "claude-opus-4" {
					t.Errorf("Model = %q, want %q", s.Model, "claude-opus-4")
				}
			},
		},
		{
			name: "model defaults to Unknown when both empty",
			content: func() string {
				data := StatusLineData{
					SessionID:  "session-789",
					CWD:        "/tmp/test",
					UpdateTime: time.Now().UnixMilli(),
				}
				b, _ := json.Marshal(data)
				return string(b)
			}(),
			wantErr: false,
			validate: func(t *testing.T, s *Session) {
				if s.Model != "Unknown" {
					t.Errorf("Model = %q, want %q", s.Model, "Unknown")
				}
			},
		},
		{
			name: "project name extracted from CWD",
			content: func() string {
				data := StatusLineData{
					SessionID:  "session-abc",
					CWD:        "/home/user/projects/my-app",
					UpdateTime: time.Now().UnixMilli(),
				}
				b, _ := json.Marshal(data)
				return string(b)
			}(),
			wantErr: false,
			validate: func(t *testing.T, s *Session) {
				if s.ProjectName != "my-app" {
					t.Errorf("ProjectName = %q, want %q", s.ProjectName, "my-app")
				}
			},
		},
		{
			name: "project dir falls back to CWD",
			content: func() string {
				data := StatusLineData{
					SessionID:  "session-def",
					CWD:        "/home/user/code",
					UpdateTime: time.Now().UnixMilli(),
				}
				// Workspace.ProjectDir intentionally empty
				b, _ := json.Marshal(data)
				return string(b)
			}(),
			wantErr: false,
			validate: func(t *testing.T, s *Session) {
				if s.ProjectDir != "/home/user/code" {
					t.Errorf("ProjectDir = %q, want %q", s.ProjectDir, "/home/user/code")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			filePath := filepath.Join(tmpDir, "test-"+tt.name+".json")
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			session, err := parseSessionFile(filePath)

			if tt.wantErr {
				if err == nil {
					t.Errorf("parseSessionFile() expected error, got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("parseSessionFile() error = %q, want error containing %q", err.Error(), tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("parseSessionFile() unexpected error: %v", err)
				return
			}

			if tt.validate != nil {
				tt.validate(t, session)
			}
		})
	}
}

func TestParseSessionFile_NonExistentFile(t *testing.T) {
	_, err := parseSessionFile("/nonexistent/path/file.json")
	if err == nil {
		t.Error("parseSessionFile() expected error for non-existent file")
	}
}

func TestSessionStatus(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name       string
		updateTime time.Time
		wantStatus SessionStatus
	}{
		{
			name:       "active session (updated just now)",
			updateTime: time.Now(),
			wantStatus: StatusActive,
		},
		{
			name:       "idle session (updated 10 minutes ago)",
			updateTime: time.Now().Add(-10 * time.Minute),
			wantStatus: StatusIdle,
		},
		{
			name:       "closed session (updated 2 hours ago)",
			updateTime: time.Now().Add(-2 * time.Hour),
			wantStatus: StatusClosed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := StatusLineData{
				SessionID:  "test-session",
				CWD:        "/tmp",
				UpdateTime: tt.updateTime.UnixMilli(),
			}
			content, _ := json.Marshal(data)

			filePath := filepath.Join(tmpDir, "status-test.json")
			if err := os.WriteFile(filePath, content, 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			session, err := parseSessionFile(filePath)
			if err != nil {
				t.Fatalf("parseSessionFile() unexpected error: %v", err)
			}

			if session.Status != tt.wantStatus {
				t.Errorf("Status = %q, want %q", session.Status, tt.wantStatus)
			}
		})
	}
}

func TestGetStatusIcon(t *testing.T) {
	tests := []struct {
		status SessionStatus
		want   string
	}{
		{StatusActive, "●"},
		{StatusIdle, "◐"},
		{StatusClosed, "○"},
		{SessionStatus("unknown"), "?"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			s := &Session{Status: tt.status}
			if got := s.GetStatusIcon(); got != tt.want {
				t.Errorf("GetStatusIcon() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetStatusText(t *testing.T) {
	tests := []struct {
		status SessionStatus
		want   string
	}{
		{StatusActive, "Active"},
		{StatusIdle, "Idle"},
		{StatusClosed, "Closed"},
		{SessionStatus("unknown"), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			s := &Session{Status: tt.status}
			if got := s.GetStatusText(); got != tt.want {
				t.Errorf("GetStatusText() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetDurationString(t *testing.T) {
	tests := []struct {
		name       string
		durationMs int64
		want       string
	}{
		{"zero minutes", 0, "0m"},
		{"30 seconds", 30000, "0m"},
		{"5 minutes", 300000, "5m"},
		{"59 minutes", 3540000, "59m"},
		{"1 hour", 3600000, "1h 0m"},
		{"1 hour 30 minutes", 5400000, "1h 30m"},
		{"2 hours 15 minutes", 8100000, "2h 15m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{DurationMs: tt.durationMs}
			if got := s.GetDurationString(); got != tt.want {
				t.Errorf("GetDurationString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetLastActivityString(t *testing.T) {
	tests := []struct {
		name         string
		lastActivity time.Time
		want         string
	}{
		{"just now", time.Now().Add(-30 * time.Second), "just now"},
		{"5 minutes ago", time.Now().Add(-5 * time.Minute), "5m ago"},
		{"45 minutes ago", time.Now().Add(-45 * time.Minute), "45m ago"},
		{"2 hours ago", time.Now().Add(-2 * time.Hour), "2h ago"},
		{"23 hours ago", time.Now().Add(-23 * time.Hour), "23h ago"},
		{"1 day ago", time.Now().Add(-25 * time.Hour), "1d ago"},
		{"3 days ago", time.Now().Add(-72 * time.Hour), "3d ago"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{LastActivity: tt.lastActivity}
			if got := s.GetLastActivityString(); got != tt.want {
				t.Errorf("GetLastActivityString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetResumeCommand(t *testing.T) {
	s := &Session{
		SessionID: "abc-123-def",
		CWD:       "/home/user/my project",
	}

	want := `cd "/home/user/my project" && claude -r abc-123-def`
	if got := s.GetResumeCommand(); got != want {
		t.Errorf("GetResumeCommand() = %q, want %q", got, want)
	}
}

func TestFilterRecentSessions(t *testing.T) {
	now := time.Now()

	sessions := []*Session{
		{SessionID: "recent-1", LastActivity: now.Add(-1 * time.Hour)},
		{SessionID: "recent-2", LastActivity: now.Add(-24 * time.Hour)},
		{SessionID: "recent-3", LastActivity: now.Add(-6 * 24 * time.Hour)}, // 6 days old
		{SessionID: "old-1", LastActivity: now.Add(-8 * 24 * time.Hour)},    // 8 days old
		{SessionID: "old-2", LastActivity: now.Add(-30 * 24 * time.Hour)},   // 30 days old
	}

	filtered := FilterRecentSessions(sessions)

	if len(filtered) != 3 {
		t.Errorf("FilterRecentSessions() returned %d sessions, want 3", len(filtered))
	}

	// Verify the correct sessions are kept
	wantIDs := map[string]bool{"recent-1": true, "recent-2": true, "recent-3": true}
	for _, s := range filtered {
		if !wantIDs[s.SessionID] {
			t.Errorf("FilterRecentSessions() kept unexpected session %q", s.SessionID)
		}
	}
}

func TestFilterRecentSessions_Empty(t *testing.T) {
	filtered := FilterRecentSessions([]*Session{})
	if len(filtered) != 0 {
		t.Errorf("FilterRecentSessions() returned %d sessions for empty input, want 0", len(filtered))
	}
}

func TestFilterRecentSessions_AllOld(t *testing.T) {
	now := time.Now()
	sessions := []*Session{
		{SessionID: "old-1", LastActivity: now.Add(-10 * 24 * time.Hour)},
		{SessionID: "old-2", LastActivity: now.Add(-20 * 24 * time.Hour)},
	}

	filtered := FilterRecentSessions(sessions)
	if len(filtered) != 0 {
		t.Errorf("FilterRecentSessions() returned %d sessions, want 0", len(filtered))
	}
}

func TestLoadSessions_NonExistentDirectory(t *testing.T) {
	// Temporarily change HOME to a non-existent directory
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)

	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)
	// Don't create .claude_sessions directory

	sessions, err := LoadSessions()
	if err != nil {
		t.Errorf("LoadSessions() unexpected error for non-existent dir: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("LoadSessions() returned %d sessions, want 0", len(sessions))
	}
}

func TestLoadSessions_WithValidFiles(t *testing.T) {
	// Create a temporary HOME directory
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	// Create .claude_sessions directory
	sessionsDir := filepath.Join(tmpDir, ".claude_sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("Failed to create sessions dir: %v", err)
	}

	// Create test session files
	now := time.Now()
	sessions := []struct {
		filename   string
		sessionID  string
		updateTime time.Time
	}{
		{"claude-status-project1.json", "session-1", now.Add(-1 * time.Minute)},
		{"claude-status-project2.json", "session-2", now.Add(-5 * time.Minute)},
		{"claude-status-project3.json", "session-3", now.Add(-10 * time.Minute)},
	}

	for _, s := range sessions {
		data := StatusLineData{
			SessionID:  s.sessionID,
			CWD:        "/tmp/test",
			UpdateTime: s.updateTime.UnixMilli(),
		}
		content, _ := json.Marshal(data)
		filePath := filepath.Join(sessionsDir, s.filename)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
	}

	// Also create some files that should be ignored
	os.WriteFile(filepath.Join(sessionsDir, "other-file.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(sessionsDir, "claude-status-invalid.txt"), []byte("{}"), 0644)
	os.MkdirAll(filepath.Join(sessionsDir, "claude-status-dir.json"), 0755) // Directory, not file

	loaded, err := LoadSessions()
	if err != nil {
		t.Fatalf("LoadSessions() unexpected error: %v", err)
	}

	if len(loaded) != 3 {
		t.Errorf("LoadSessions() returned %d sessions, want 3", len(loaded))
	}

	// Verify sessions are sorted by last activity (most recent first)
	if len(loaded) >= 2 {
		if loaded[0].SessionID != "session-1" {
			t.Errorf("First session = %q, want %q (most recent)", loaded[0].SessionID, "session-1")
		}
		if loaded[len(loaded)-1].SessionID != "session-3" {
			t.Errorf("Last session = %q, want %q (oldest)", loaded[len(loaded)-1].SessionID, "session-3")
		}
	}
}

func TestLoadSessions_SkipsInvalidFiles(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer os.Setenv("HOME", oldHome)
	os.Setenv("HOME", tmpDir)

	sessionsDir := filepath.Join(tmpDir, ".claude_sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("Failed to create sessions dir: %v", err)
	}

	// Create one valid and one invalid session file
	validData := StatusLineData{
		SessionID:  "valid-session",
		CWD:        "/tmp",
		UpdateTime: time.Now().UnixMilli(),
	}
	validContent, _ := json.Marshal(validData)
	os.WriteFile(filepath.Join(sessionsDir, "claude-status-valid.json"), validContent, 0644)

	// Invalid: empty file
	os.WriteFile(filepath.Join(sessionsDir, "claude-status-empty.json"), []byte(""), 0644)

	// Invalid: bad JSON
	os.WriteFile(filepath.Join(sessionsDir, "claude-status-bad.json"), []byte("not json"), 0644)

	// Invalid: no update time
	os.WriteFile(filepath.Join(sessionsDir, "claude-status-notime.json"), []byte(`{"session_id":"x"}`), 0644)

	loaded, err := LoadSessions()
	if err != nil {
		t.Fatalf("LoadSessions() unexpected error: %v", err)
	}

	// Should only have the valid session
	if len(loaded) != 1 {
		t.Errorf("LoadSessions() returned %d sessions, want 1", len(loaded))
	}
	if len(loaded) > 0 && loaded[0].SessionID != "valid-session" {
		t.Errorf("LoadSessions() got session %q, want %q", loaded[0].SessionID, "valid-session")
	}
}

// Helper function for checking if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
