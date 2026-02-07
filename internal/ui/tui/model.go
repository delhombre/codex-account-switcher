// Package tui provides the interactive terminal user interface.
package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/delhombre/cxa/internal/account"
	"github.com/delhombre/cxa/internal/ui/styles"
)

// Repository interface for the TUI
type Repository interface {
	List() ([]*account.Account, error)
	Current() (string, error)
	Activate(name string) error
	Save(name string) (*account.Account, error)
}

// accountItem implements list.Item for accounts
type accountItem struct {
	account   *account.Account
	isCurrent bool
}

func (i accountItem) Title() string {
	if i.isCurrent {
		return styles.CurrentAccountStyle.Render(i.account.Name) + " " + styles.MutedStyle.Render("(current)")
	}
	return i.account.Name
}

func (i accountItem) Description() string {
	if i.account.Email != "" {
		return i.account.Email
	}
	return styles.MutedStyle.Render("Press enter to switch")
}

func (i accountItem) FilterValue() string {
	return i.account.Name
}

// Model is the main TUI model
type Model struct {
	list     list.Model
	repo     Repository
	current  string
	quitting bool
	message  string
	err      error
}

// NewModel creates a new TUI model
func NewModel(repo Repository) (*Model, error) {
	accounts, err := repo.List()
	if err != nil {
		return nil, err
	}

	current, _ := repo.Current()

	items := make([]list.Item, len(accounts))
	for i, acc := range accounts {
		items[i] = accountItem{
			account:   acc,
			isCurrent: acc.Name == current,
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().
		Foreground(styles.Primary).
		Bold(true).
		Padding(0, 0, 0, 2)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().
		Foreground(styles.TextDim).
		Padding(0, 0, 0, 2)

	l := list.New(items, delegate, 50, 14)
	l.Title = "Codex Accounts"
	l.Styles.Title = styles.HeaderStyle
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)

	return &Model{
		list:    l,
		repo:    repo,
		current: current,
	}, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("q", "ctrl+c"))):
			m.quitting = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if item, ok := m.list.SelectedItem().(accountItem); ok {
				if item.account.Name != m.current {
					if err := m.repo.Activate(item.account.Name); err != nil {
						m.err = err
						m.message = styles.RenderError(err.Error())
					} else {
						m.current = item.account.Name
						m.message = styles.RenderSuccess(fmt.Sprintf("Switched to %s", item.account.Name))
						// Refresh list
						m.refreshList()
					}
				}
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		h := msg.Height - 4
		if h < 5 {
			h = 5
		}
		m.list.SetHeight(h)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Model) refreshList() {
	accounts, _ := m.repo.List()
	items := make([]list.Item, len(accounts))
	for i, acc := range accounts {
		items[i] = accountItem{
			account:   acc,
			isCurrent: acc.Name == m.current,
		}
	}
	m.list.SetItems(items)
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	// Main list
	b.WriteString(m.list.View())

	// Message/error
	if m.message != "" {
		b.WriteString("\n\n")
		b.WriteString(m.message)
	}

	// Help
	b.WriteString("\n\n")
	b.WriteString(styles.MutedStyle.Render("  enter: switch  •  /: filter  •  q: quit"))

	return b.String()
}

// Run starts the TUI
func Run(repo Repository) error {
	model, err := NewModel(repo)
	if err != nil {
		return err
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
