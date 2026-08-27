package main

import (
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var inputStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#FFF"))

type todoItem struct {
	done bool
	name string
}

func (t todoItem) FilterValue() string {
	return t.name
}

func (t todoItem) Title() string {
	return t.name
}

func (t todoItem) Description() string {
	return t.name
}

type model struct {
	input    textinput.Model
	list     list.Model
	delegate *todoDelegate
	width    int
	height   int
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "New placeholder"
	ti.Focus()

	delegate := newTodoDelegate()
	li := list.New([]list.Item{}, delegate, 0, 0)
	li.SetShowStatusBar(false)
	li.DisableQuitKeybindings()

	return model{
		input:    ti,
		list:     li,
		delegate: delegate,
	}
}

type (
	loadedTodosMsg []list.Item
	loadingTodos   struct{}
)

func (m model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg { return loadingTodos{} },
		func() tea.Msg {
			<-time.After(1 * time.Second)
			return loadedTodosMsg{
				todoItem{done: false, name: "Item 1"},
				todoItem{done: false, name: "Item 2"},
				todoItem{done: false, name: "Item 3"},
			}
		},
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			todo := todoItem{name: m.input.Value()}
			m.input.SetValue("")
			cmd := m.list.InsertItem(0, todo)
			return m, cmd
		case tea.KeyDown:
			if m.input.Focused() {
				m.input.Blur()
				m.delegate.focused = true
				return m, nil
			}
		case tea.KeyUp:
			if m.input.Focused() == false && m.list.Cursor() == 0 {
				m.input.Focus()
				m.delegate.focused = false
				return m, nil
			}
		}
	case loadingTodos:
		return m, m.list.ToggleSpinner()
	case loadedTodosMsg:
		return m, tea.Batch(
			m.list.ToggleSpinner(),
			m.list.SetItems(msg),
		)
	case tea.WindowSizeMsg:
		m.input.Width = msg.Width - 2
		m.list.SetSize(msg.Width, msg.Height-3)
		m.width = msg.Width
		m.height = msg.Height
	case todoToggleMsg:
		todo := m.list.Items()[msg.index].(todoItem)
		todo.done = !todo.done
		return m, m.list.SetItem(msg.index, todo)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	var inputCmd, listCmd tea.Cmd
	if m.input.Focused() {
		m.input, inputCmd = m.input.Update(msg)
	} else {
		m.list, listCmd = m.list.Update(msg)
	}

	return m, tea.Batch(inputCmd, listCmd)
}

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top,
		inputStyle.Width(m.width-2).Height(1).Render(m.input.View()),
		m.list.View(),
	)
}

func main() {
	m := newModel()
	tea.NewProgram(m, tea.WithAltScreen()).Run()
}
