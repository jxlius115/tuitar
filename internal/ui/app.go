// internal/ui/app.go
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Cod-e-Codes/tuitar/internal/models"
	"github.com/Cod-e-Codes/tuitar/internal/storage"
	"github.com/Cod-e-Codes/tuitar/internal/ui/components"
)

type Model struct {
	state      models.SessionState
	storage    storage.Storage
	tabs       []models.Tab

	// Components
	tabEditor  components.TabEditorModel
	tabBrowser components.TabBrowserModel
	statusBar  components.StatusBarModel
	help       help.Model

	// UI State
	windowSize tea.WindowSizeMsg
	showHelp   bool
	keys       KeyMap
}

type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Enter     key.Binding
	Quit      key.Binding
	Help      key.Binding
	Save      key.Binding
	New       key.Binding
	Play      key.Binding
	Insert    key.Binding
	Normal    key.Binding
	Browser   key.Binding
	Delete    key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.Enter, k.Quit, k.Help},
		{k.Save, k.New, k.Play},
		{k.Insert, k.Normal, k.Browser, k.Delete},
	}
}

func NewKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("h", "left"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys("l", "right"),
			key.WithHelp("→/l", "move right"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select/edit"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Save: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "save tab"),
		),
		New: key.NewBinding(
			key.WithKeys("ctrl+n"),
			key.WithHelp("ctrl+n", "new tab"),
		),
		Play: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "play/pause"),
		),
		Insert: key.NewBinding(
			key.WithKeys("i"),
			key.WithHelp("i", "insert mode"),
		),
		Normal: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "normal mode"),
		),
		Browser: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "toggle browser"),
		),
		Delete: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "delete fret"),
		),
	}
}

func NewModel(storage storage.Storage) Model {
	tabs, _ := storage.LoadAllTabs()

	m := Model{
		storage:    storage,
		tabs:       tabs,
		keys:       NewKeyMap(),
		help:       help.New(),
		tabBrowser: components.NewTabBrowser(tabs),
		statusBar:  components.NewStatusBar(),
	}

	m.state.ViewMode = models.ViewBrowser
	m.state.EditMode = models.EditNormal

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.SetWindowTitle("Guitar Tab TUI")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg
		m.tabEditor.SetSize(msg.Width, msg.Height-3) // Reserve space for status bar
		m.tabBrowser.SetSize(msg.Width, msg.Height-3)

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			return m, nil

		case key.Matches(msg, m.keys.Browser) && m.state.ViewMode == models.ViewEditor:
			// Only allow browser switch from editor, not from browser to avoid conflicts
			m.state.ViewMode = models.ViewBrowser
			return m, nil

		case key.Matches(msg, m.keys.New):
			newTab := models.NewEmptyTab("New Tab")
			m.state.CurrentTab = newTab
			m.tabEditor = components.NewTabEditor(newTab)
			m.tabEditor.SetEditMode(models.EditNormal)
			m.state.ViewMode = models.ViewEditor
			m.state.EditMode = models.EditNormal
			return m, nil

		case key.Matches(msg, m.keys.Save):
			if m.state.CurrentTab != nil {
				err := m.storage.SaveTab(m.state.CurrentTab)
				if err != nil {
					m.statusBar.SetStatus("Error saving tab: " + err.Error())
				} else {
					m.statusBar.SetStatus("Tab saved successfully")
					// Refresh tabs list
					if tabs, err := m.storage.LoadAllTabs(); err == nil {
						m.tabs = tabs
						m.tabBrowser.SetTabs(tabs)
					}
				}
			}
			return m, nil
		}

		// Handle view-specific key presses
		switch m.state.ViewMode {
		case models.ViewBrowser:
			return m.updateBrowser(msg)
		case models.ViewEditor:
			return m.updateEditor(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) updateBrowser(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Enter):
		if len(m.tabs) > 0 && m.tabBrowser.Cursor() < len(m.tabs) {
			selectedTab := &m.tabs[m.tabBrowser.Cursor()]
			// Create a copy of the tab to avoid modifying the original
			tabCopy := *selectedTab
			m.state.CurrentTab = &tabCopy
			m.tabEditor = components.NewTabEditor(&tabCopy)
			m.tabEditor.SetEditMode(models.EditNormal)
			m.state.ViewMode = models.ViewEditor
			m.state.EditMode = models.EditNormal
		}
		return m, nil
	}

	m.tabBrowser, cmd = m.tabBrowser.Update(msg)
	return m, cmd
}

func (m Model) updateEditor(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(msg, m.keys.Insert):
		m.state.EditMode = models.EditInsert
		m.tabEditor.SetEditMode(models.EditInsert)
		m.statusBar.SetStatus("-- INSERT MODE --")
		return m, nil

	case key.Matches(msg, m.keys.Normal):
		m.state.EditMode = models.EditNormal
		m.tabEditor.SetEditMode(models.EditNormal)
		m.statusBar.SetStatus("-- NORMAL MODE --")
		return m, nil
	}

	// Pass the message to the tab editor for handling
	m.tabEditor, cmd = m.tabEditor.Update(msg)
	
	// Update the current tab if it has changed
	if m.tabEditor.HasChanged() {
		m.state.CurrentTab = m.tabEditor.GetTab()
	}

	return m, cmd
}

func (m Model) View() string {
	if m.showHelp {
		return m.help.View(m.keys)
	}

	var content string

	switch m.state.ViewMode {
	case models.ViewBrowser:
		content = m.renderBrowser()
	case models.ViewEditor:
		content = m.renderEditor()
	}

	statusBar := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left, content, statusBar)
}

func (m Model) renderBrowser() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render("Guitar Tab Browser")

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Render("Enter: Edit • Ctrl+N: New • ?: Help • Q: Quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		m.tabBrowser.View(),
		"",
		help,
	)
}

func (m Model) renderEditor() string {
	if m.state.CurrentTab == nil {
		return "No tab selected"
	}

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		Render(fmt.Sprintf("Editing: %s", m.state.CurrentTab.Name))

	mode := "NORMAL"
	modeColor := lipgloss.Color("12")
	if m.state.EditMode == models.EditInsert {
		mode = "INSERT"
		modeColor = lipgloss.Color("11")
	}

	modeIndicator := lipgloss.NewStyle().
		Bold(true).
		Foreground(modeColor).
		Render(fmt.Sprintf("-- %s --", mode))

	var help string
	if m.state.EditMode == models.EditInsert {
		help = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render("0-9: Insert fret • -: Rest • Esc: Normal mode • Arrows: Navigate")
	} else {
		help = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Render("I: Insert • X: Delete • Ctrl+S: Save • Tab: Browser • Space: Play • Arrows: Navigate")
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		m.tabEditor.View(),
		"",
		modeIndicator,
		help,
	)
}
