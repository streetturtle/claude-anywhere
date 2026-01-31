package session

import "testing"

func TestFormatTokens(t *testing.T) {
	tests := []struct {
		name   string
		tokens int
		want   string
	}{
		{"zero", 0, "0"},
		{"small number", 100, "100"},
		{"just under 1K", 999, "999"},
		{"exactly 1K", 1000, "1.0K"},
		{"1.5K", 1500, "1.5K"},
		{"10K", 10000, "10.0K"},
		{"999K", 999000, "999.0K"},
		{"exactly 1M", 1000000, "1.0M"},
		{"1.5M", 1500000, "1.5M"},
		{"10M", 10000000, "10.0M"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatTokens(tt.tokens); got != tt.want {
				t.Errorf("formatTokens(%d) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

func TestItemFilterValue(t *testing.T) {
	s := &Session{
		ProjectName: "my-project",
		Model:       "Claude Sonnet",
	}
	i := item{session: s}

	got := i.FilterValue()
	want := "my-project Claude Sonnet"

	if got != want {
		t.Errorf("FilterValue() = %q, want %q", got, want)
	}
}

func TestItemTitle(t *testing.T) {
	tests := []struct {
		name        string
		status      SessionStatus
		projectName string
		wantContain string
	}{
		{"active session", StatusActive, "test-project", "test-project"},
		{"idle session", StatusIdle, "my-app", "my-app"},
		{"closed session", StatusClosed, "old-project", "old-project"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				Status:      tt.status,
				ProjectName: tt.projectName,
			}
			i := item{session: s}

			got := i.Title()
			if !contains(got, tt.wantContain) {
				t.Errorf("Title() = %q, want to contain %q", got, tt.wantContain)
			}
		})
	}
}

func TestItemDescription(t *testing.T) {
	s := &Session{
		Model:          "Claude Opus",
		CostUSD:        0.0150,
		ContextUsedPct: 75.5,
	}
	i := item{session: s}

	desc := i.Description()

	// Check that description contains expected parts
	expectedParts := []string{
		"Model: Claude Opus",
		"Cost: $0.0150",
		"Context: 76%",
	}

	for _, part := range expectedParts {
		if !contains(desc, part) {
			t.Errorf("Description() = %q, want to contain %q", desc, part)
		}
	}
}
