package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	choices := []string{"Pizza", "Burger", "Salad", "Pasta"}
	menu := NewMenu(choices)
	p := tea.NewProgram(menu)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running the program: %s\n", err)
		os.Exit(1)
	}
}

func (m Menu) Init() tea.Cmd { return nil }
func (m Menu) View() string {
	s := "What would you like to eat?\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s
}

func (m Menu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	l := len(m.choices)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+q", "q":
			return m, tea.Quit
		case "up", "k":
			m.cursor = (m.cursor - 1 + l) % l
		case "down", "j":
			m.cursor = (m.cursor + 1) % l
		case "enter", "return":
			fmt.Printf("You chose: %s\n\n", m.choices[m.cursor])
			return m, tea.Quit
		}
	}
	return m, nil
}
