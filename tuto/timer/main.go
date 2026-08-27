package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	count int
}

type tickMsg struct{}

func (m model) Init() tea.Cmd {
	return tickCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.count++
		}
	case tickMsg:
		m.count++
		return m, tickCmd
	}
	return m, nil
}

func tickCmd() tea.Msg {
	<-time.After(1 * time.Second)
	return tickMsg{}
}

func (m model) View() string {
	return fmt.Sprintf("Count: %d", m.count)
}

func newModel() model {
	return model{}
}

func main() {
	m := newModel()
	tea.NewProgram(m).Run()
}
