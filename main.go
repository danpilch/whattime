package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Margin(1, 0)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Margin(0, 0, 1, 0)

	tableStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Margin(1, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			Margin(1, 0)

	searchStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(0, 1).
			Margin(0, 0, 1, 0)
)

type model struct {
	table         table.Model
	spinner       spinner.Model
	search        textinput.Model
	loading       bool
	error         error
	userTimezones []UserTimezone
	allUsers      []UserTimezone
	slackClient   *SlackClient
	searchFocused bool
}

type tickMsg time.Time
type loadCompleteMsg []UserTimezone
type loadErrorMsg error

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.loadUsers,
		tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			return m, tea.Quit
		case key.Matches(msg, key.NewBinding(key.WithKeys("r"))):
			m.loading = true
			m.error = nil
			cmds = append(cmds, m.loadUsers)
		case key.Matches(msg, key.NewBinding(key.WithKeys("/"))):
			if !m.searchFocused {
				m.searchFocused = true
				m.search.Focus()
				m.table.Blur()
				return m, tea.Batch(cmds...)
			}
		case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
			m.searchFocused = false
			m.search.Blur()
			m.table.Focus()
			m.search.SetValue("")
			m.filterUsers()
		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.searchFocused {
				m.searchFocused = false
				m.search.Blur()
				m.table.Focus()
			}
		}

		if m.searchFocused {
			m.search, cmd = m.search.Update(msg)
			cmds = append(cmds, cmd)
			m.filterUsers()
		}

	case tickMsg:
		m.updateTable()
		cmds = append(cmds, tickCmd())

	case loadCompleteMsg:
		m.loading = false
		m.allUsers = []UserTimezone(msg)
		m.userTimezones = m.allUsers
		m.updateTable()

	case loadErrorMsg:
		m.loading = false
		m.error = error(msg)

	case spinner.TickMsg:
		if m.loading {
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if !m.loading && !m.searchFocused {
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	var s strings.Builder

	s.WriteString(titleStyle.Render("⏰ Coworker Timezones"))
	s.WriteString("\n")
	s.WriteString(subtitleStyle.Render("Real-time timezone information from Slack"))
	s.WriteString("\n")

	if m.error != nil {
		s.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.error)))
		s.WriteString("\n")
		s.WriteString(helpStyle.Render("Press 'r' to retry, 'q' to quit"))
		return s.String()
	}

	s.WriteString(searchStyle.Render(fmt.Sprintf("🔍 Search: %s", m.search.View())))
	s.WriteString("\n")

	if m.loading {
		s.WriteString(fmt.Sprintf("%s Loading coworker timezones...", m.spinner.View()))
		s.WriteString("\n")
	} else {
		s.WriteString(tableStyle.Render(m.table.View()))
		s.WriteString("\n")
	}

	if m.searchFocused {
		s.WriteString(helpStyle.Render("Press 'enter' to apply search, 'esc' to cancel"))
	} else {
		s.WriteString(helpStyle.Render("Press '/' to search, 'r' to refresh, 'q' to quit"))
	}
	return s.String()
}

func (m *model) updateTable() {
	if len(m.userTimezones) == 0 {
		return
	}

	rows := make([]table.Row, len(m.userTimezones))
	for i, user := range m.userTimezones {
		currentTime := user.GetCurrentTime()
		timeStr := currentTime.Format("3:04 PM")
		dateStr := currentTime.Format("Jan 2")

		statusEmoji := getStatusEmoji(currentTime)

		rows[i] = table.Row{
			user.Name,
			user.Username,
			user.Timezone,
			timeStr,
			dateStr,
			statusEmoji,
			user.GetTimezoneOffset(),
		}
	}

	m.table.SetRows(rows)
}

func (m *model) filterUsers() {
	searchTerm := strings.ToLower(m.search.Value())
	if searchTerm == "" {
		m.userTimezones = m.allUsers
		m.updateTable()
		return
	}

	var filtered []UserTimezone
	for _, user := range m.allUsers {
		if strings.Contains(strings.ToLower(user.Name), searchTerm) ||
			strings.Contains(strings.ToLower(user.Username), searchTerm) ||
			strings.Contains(strings.ToLower(user.Timezone), searchTerm) {
			filtered = append(filtered, user)
		}
	}

	m.userTimezones = filtered
	m.updateTable()
}

func (m model) loadUsers() tea.Msg {
	users, err := m.slackClient.GetUserTimezones()
	if err != nil {
		return loadErrorMsg(err)
	}
	return loadCompleteMsg(users)
}

func getStatusEmoji(t time.Time) string {
	hour := t.Hour()
	switch {
	case hour >= 9 && hour < 12:
		return "🌅 Morning"
	case hour >= 12 && hour < 17:
		return "☀️ Afternoon"
	case hour >= 17 && hour < 21:
		return "🌆 Evening"
	case hour >= 21 || hour < 6:
		return "🌙 Night"
	default:
		return "🌄 Early"
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialModel() model {
	token := getSlackToken()
	if token == "" {
		return model{
			error: fmt.Errorf("SLACK_BOT_TOKEN or SLACK_USER_TOKEN environment variable is required"),
		}
	}

	slackClient := NewSlackClient(token)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search by name, username, or timezone..."
	searchInput.CharLimit = 156
	searchInput.Width = 50

	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Username", Width: 15},
		{Title: "Timezone", Width: 25},
		{Title: "Current Time", Width: 12},
		{Title: "Date", Width: 10},
		{Title: "Status", Width: 15},
		{Title: "Offset", Width: 8},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(15),
	)

	tableStyles := table.DefaultStyles()
	tableStyles.Header = tableStyles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tableStyles.Selected = tableStyles.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(tableStyles)

	return model{
		table:       t,
		spinner:     s,
		search:      searchInput,
		loading:     true,
		slackClient: slackClient,
	}
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
