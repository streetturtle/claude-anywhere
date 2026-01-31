package session

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 25

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	paginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	activeStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("46"))  // Green
	idleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("226")) // Yellow
	closedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // Gray

	// Card styles
	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("238")).
			Padding(0, 1).
			MarginLeft(2).
			MarginBottom(0)

	selectedCardStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("170")).
				Padding(0, 1).
				MarginLeft(2).
				MarginBottom(0)

	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	valueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
)

type item struct {
	session *Session
}

func (i item) FilterValue() string {
	// Make searchable by project name and model only
	return fmt.Sprintf("%s %s",
		i.session.ProjectName,
		i.session.Model)
}

func (i item) Title() string {
	statusIcon := i.session.GetStatusIcon()

	var styledIcon string
	switch i.session.Status {
	case StatusActive:
		styledIcon = activeStyle.Render(statusIcon)
	case StatusIdle:
		styledIcon = idleStyle.Render(statusIcon)
	case StatusClosed:
		styledIcon = closedStyle.Render(statusIcon)
	}

	return fmt.Sprintf("%s %s", styledIcon, i.session.ProjectName)
}

func (i item) Description() string {
	s := i.session

	parts := []string{
		fmt.Sprintf("Model: %s", s.Model),
		fmt.Sprintf("Cost: $%.4f", s.CostUSD),
		fmt.Sprintf("Context: %.0f%%", s.ContextUsedPct),
		fmt.Sprintf("Updated: %s", s.GetLastActivityString()),
	}

	return strings.Join(parts, " • ")
}

type model struct {
	list     list.Model
	sessions []*Session
	choice   *Session
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "q":
			// Only quit with 'q' if not filtering
			if m.list.FilterState() != list.Filtering {
				m.quitting = true
				return m, tea.Quit
			}

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.session
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	return "\n" + m.list.View()
}

// RunTUI starts the interactive terminal UI for session selection
func RunTUI(sessions []*Session) error {
	items := make([]list.Item, len(sessions))
	totalCost := 0.0
	activeCount := 0
	idleCount := 0

	for i, s := range sessions {
		items[i] = item{session: s}
		totalCost += s.CostUSD
		if s.Status == StatusActive {
			activeCount++
		} else if s.Status == StatusIdle {
			idleCount++
		}
	}

	const defaultWidth = 120

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)

	// Create title with summary stats
	title := fmt.Sprintf("Claude Sessions  │  %d total  │  %d active  │  %d idle  │  $%.4f total",
		len(sessions), activeCount, idleCount, totalCost)
	l.Title = title

	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// Add additional help keys
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "resume"),
			),
			key.NewBinding(
				key.WithKeys("/"),
				key.WithHelp("/", "search"),
			),
		}
	}

	m := model{
		list:     l,
		sessions: sessions,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running TUI: %w", err)
	}

	// Check if user made a selection
	if m := finalModel.(model); m.choice != nil {
		return resumeSession(m.choice)
	}

	return nil
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 3 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	s := i.session

	// Status icon with color
	statusIcon := s.GetStatusIcon()
	var styledIcon string
	switch s.Status {
	case StatusActive:
		styledIcon = activeStyle.Render(statusIcon)
	case StatusIdle:
		styledIcon = idleStyle.Render(statusIcon)
	case StatusClosed:
		styledIcon = closedStyle.Render(statusIcon)
	}

	// Format tokens
	tokensIn := formatTokens(s.InputTokens)
	tokensOut := formatTokens(s.OutputTokens)

	// Format lines changed
	linesChanged := ""
	if s.LinesAdded > 0 || s.LinesRemoved > 0 {
		linesChanged = fmt.Sprintf("+%d -%d", s.LinesAdded, s.LinesRemoved)
	} else {
		linesChanged = "no changes"
	}

	// Build compact card - all info on 2 lines
	title := fmt.Sprintf("%s %s", styledIcon, valueStyle.Render(s.ProjectName))

	// Compact info line with separators
	details := fmt.Sprintf("%s │ %s │ %s │ %s │ %s in/%s out │ %s",
		valueStyle.Render(s.Model),
		valueStyle.Render(fmt.Sprintf("$%.4f", s.CostUSD)),
		valueStyle.Render(fmt.Sprintf("%.0f%%", s.ContextUsedPct)),
		labelStyle.Render(s.GetLastActivityString()),
		labelStyle.Render(tokensIn),
		labelStyle.Render(tokensOut),
		labelStyle.Render(linesChanged))

	cardContent := fmt.Sprintf("%s\n%s", title, details)

	// Apply card style based on selection
	var rendered string
	if index == m.Index() {
		rendered = selectedCardStyle.Render(cardContent)
	} else {
		rendered = cardStyle.Render(cardContent)
	}

	fmt.Fprint(w, rendered)
}

// formatTokens formats token count with K/M suffixes
func formatTokens(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	} else if tokens < 1000000 {
		return fmt.Sprintf("%.1fK", float64(tokens)/1000)
	}
	return fmt.Sprintf("%.1fM", float64(tokens)/1000000)
}

// resumeSession changes to the session's directory and runs claude with the session ID
func resumeSession(s *Session) error {
	fmt.Printf("\n📂 Changing to directory: %s\n", s.CWD)
	fmt.Printf("🔄 Resuming session: %s\n\n", s.SessionID)

	// Change directory
	if err := os.Chdir(s.CWD); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	// Execute claude command
	cmd := exec.Command("claude", "-r", s.SessionID)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
