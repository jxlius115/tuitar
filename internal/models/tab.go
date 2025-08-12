// internal/ui/components/tab_editor.go
package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Cod-e-Codes/tuitar/internal/models"
)

// HighlightUpdateMsg is sent to update the highlighted positions for playback
type HighlightUpdateMsg struct {
	Positions []models.Position
}

type TabEditorModel struct {
	tab             *models.Tab
	cursor          models.Position
	viewport        viewport.Model
	width           int
	height          int
	changed         bool
	editMode        models.EditMode
	highlightedPos  []models.Position // For playback highlighting
}

func NewTabEditor(tab *models.Tab) TabEditorModel {
	vp := viewport.New(80, 20)

	// Initialize tab content if it's empty
	if tab.Content[0] == "" {
		tab.Content = [6]string{
			"----------------",
			"----------------",
			"----------------",
			"----------------",
			"----------------",
			"----------------",
		}
	}

	return TabEditorModel{
		tab:      tab,
		viewport: vp,
		cursor:   models.Position{String: 0, Position: 0},
		editMode: models.EditNormal,
	}
}

func (m *TabEditorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height - 4 // Reserve space for headers
}

// SetHighlightedPositions sets playback highlights and marks model as changed
func (m *TabEditorModel) SetHighlightedPositions(positions []models.Position) {
	m.highlightedPos = positions
	m.changed = true
}

// Update processes messages including external highlight updates
func (m TabEditorModel) Update(msg tea.Msg) (TabEditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case HighlightUpdateMsg:
		m.highlightedPos = msg.Positions
		m.changed = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		// Navigation keys work in both modes
		case "h", "left":
			if m.cursor.Position > 0 {
				m.cursor.Position--
			}
		case "l", "right":
			maxPos := len(m.tab.Content[m.cursor.String]) - 1
			if m.cursor.Position < maxPos {
				m.cursor.Position++
			}
		case "k", "up":
			if m.cursor.String > 0 {
				m.cursor.String--
			}
		case "j", "down":
			if m.cursor.String < 5 {
				m.cursor.String++
			}
		case "home":
			m.cursor.Position = 0
		case "end":
			m.cursor.Position = len(m.tab.Content[m.cursor.String]) - 1

		// Insert mode specific keys
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.editMode == models.EditInsert {
				m.insertCharAt(m.cursor, rune(msg.String()[0]))
				m.changed = true
				if m.cursor.Position < len(m.tab.Content[m.cursor.String])-1 {
					m.cursor.Position++
				}
			}
		case "-":
			if m.editMode == models.EditInsert {
				m.insertCharAt(m.cursor, '-')
				m.changed = true
				if m.cursor.Position < len(m.tab.Content[m.cursor.String])-1 {
					m.cursor.Position++
				}
			}

		// Delete key works in normal mode
		case "x":
			if m.editMode == models.EditNormal {
				m.deleteCharAt(m.cursor)
				m.changed = true
			}

		// Backspace works in insert mode
		case "backspace", "ctrl+h":
			if m.editMode == models.EditInsert && m.cursor.Position > 0 {
				m.cursor.Position--
				m.deleteCharAt(m.cursor)
				m.changed = true
			}
		}
	}

	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m *TabEditorModel) insertCharAt(pos models.Position, char rune) {
	line := []rune(m.tab.Content[pos.String])
	if pos.Position < len(line) {
		line[pos.Position] = char
		m.tab.Content[pos.String] = string(line)
	}
}

func (m *TabEditorModel) deleteCharAt(pos models.Position) {
	line := []rune(m.tab.Content[pos.String])
	if pos.Position < len(line) {
		line[pos.Position] = '-'
		m.tab.Content[pos.String] = string(line)
	}
}

func (m TabEditorModel) View() string {
	var lines []string

	// String labels (high to low pitch, matching guitar orientation)
	stringLabels := []string{"e", "B", "G", "D", "A", "E"}

	// Helper to check if position is highlighted
	isHighlighted := func(str, pos int) bool {
		for _, hp := range m.highlightedPos {
			if hp.String == str && hp.Position == pos {
				return true
			}
		}
		return false
	}

	for i, label := range stringLabels {
		line := lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Render(label + "|")

		// Render tab content with cursor and playback highlighting
		content := m.tab.Content[i]
		for pos, char := range content {
			style := lipgloss.NewStyle()

			// Highlight cursor position (takes precedence)
			if m.cursor.String == i && m.cursor.Position == pos {
				if m.editMode == models.EditInsert {
					style = style.Background(lipgloss.Color("11")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Background(lipgloss.Color("12")).Foreground(lipgloss.Color("15"))
				}
			} else if isHighlighted(i, pos) {
				// Highlight playback positions with cyan background
				style = style.Background(lipgloss.Color("37")).Foreground(lipgloss.Color("0"))
			}

			line += style.Render(string(char))
		}

		line += lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Render("|")

		lines = append(lines, line)
	}

	content := strings.Join(lines, "\n")
	m.viewport.SetContent(content)

	return m.viewport.View()
}

func (m TabEditorModel) HasChanged() bool {
	return m.changed
}

func (m *TabEditorModel) ResetChanged() {
	m.changed = false
}

func (m TabEditorModel) GetTab() *models.Tab {
	return m.tab
}

func (m *TabEditorModel) SetEditMode(mode models.EditMode) {
	m.editMode = mode
}

func (m TabEditorModel) GetEditMode() models.EditMode {
	return m.editMode
}

func (m TabEditorModel) GetCursor() models.Position {
	return m.cursor
}
