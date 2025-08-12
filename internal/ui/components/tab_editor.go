// internal/ui/components/tab_editor.go
package components

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Cod-e-Codes/tuitar/internal/models"
)

type TabEditorModel struct {
	tab        *models.Tab
	cursor     models.Position
	viewport   viewport.Model
	width      int
	height     int
	changed    bool
	editMode   models.EditMode
}

func NewTabEditor(tab *models.Tab) TabEditorModel {
	vp := viewport.New(80, 20)
	
	return TabEditorModel{
		tab:      tab,
		viewport: vp,
		cursor:   models.Position{String: 0, Position: 0},
	}
}

func (m *TabEditorModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.viewport.Width = width
	m.viewport.Height = height - 4 // Reserve space for headers
}

func (m TabEditorModel) Update(msg tea.Msg) (TabEditorModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
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
		case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
			if m.editMode == models.EditInsert {
				m.insertCharAt(m.cursor, rune(msg.String()[0]))
				m.changed = true
			}
		case "-":
			if m.editMode == models.EditInsert {
				m.insertCharAt(m.cursor, '-')
				m.changed = true
			}
		case "x":
			m.deleteCharAt(m.cursor)
			m.changed = true
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
	
	// String labels
	stringLabels := []string{"e", "B", "G", "D", "A", "E"}
	
	for i, label := range stringLabels {
		line := lipgloss.NewStyle().
			Foreground(lipgloss.Color("14")).
			Render(label + "|")
		
		// Render tab content with cursor highlighting
		content := m.tab.Content[i]
		for pos, char := range content {
			style := lipgloss.NewStyle()
			
			// Highlight cursor position
			if m.cursor.String == i && m.cursor.Position == pos {
				if m.editMode == models.EditInsert {
					style = style.Background(lipgloss.Color("11")).Foreground(lipgloss.Color("0"))
				} else {
					style = style.Background(lipgloss.Color("12")).Foreground(lipgloss.Color("15"))
				}
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

func (m TabEditorModel) GetTab() *models.Tab {
	return m.tab
}

func (m *TabEditorModel) SetEditMode(mode models.EditMode) {
	m.editMode = mode
}
