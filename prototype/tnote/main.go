package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type NoteModel struct {
	entryList list.Model
	editor    textarea.Model
}

func (m NoteModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m NoteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.entryList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmds []tea.Cmd
	var cmd tea.Cmd
	m.entryList, cmd = m.entryList.Update(msg)
	cmds = append(cmds, cmd)

	// idx := m.entryList.GlobalIndex()
	listItem := m.entryList.SelectedItem()
	entryItem, ok := listItem.(item)
	if ok {
		m.editor.Reset()
		m.editor.InsertString(entryItem.Description())
	}

	m.editor, cmd = m.editor.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m NoteModel) View() string {
	var s string
	s += lipgloss.JoinHorizontal(lipgloss.Top, docStyle.Render(m.entryList.View()), m.editor.View())
	return s
}

func newModel() NoteModel {

	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
		item{title: "Bitter melon", desc: "It cools you down"},
		item{title: "Nice socks", desc: "And by that I mean socks without holes"},
		item{title: "Eight hours of sleep", desc: "I had this once"},
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)

	editor := textarea.New()
	editor.Placeholder = "have a good day!"

	return NoteModel{
		entryList: list,
		editor:    editor,
	}
}

func main() {
	m := newModel()

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}
