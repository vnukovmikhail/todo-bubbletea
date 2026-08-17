package main

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Todo struct {
	Title     string
	Completed bool
}

type model struct {
	todos     []Todo
	cursor    int
	inputMode bool
	input     string
	quitting  bool
}

func initialModel() model {
	m := model{
		todos:     []Todo{},
		cursor:    0,
		inputMode: false,
		input:     "",
		quitting:  false,
	}
	if err := m.Load(); err != nil {
		os.Exit(1)
	}
	return m
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		if m.inputMode {
			return m.handleInputMode(msg)
		}
		defer m.Save()
		switch msg.String() {
		case "up":
			if len(m.todos) > 0 {
				m.cursor = (m.cursor + len(m.todos) - 1) % len(m.todos)
			}
		case "down":
			if len(m.todos) > 0 {
				m.cursor = (m.cursor + 1) % len(m.todos)
			}
		case " ":
			if len(m.todos) > 0 {
				m.todos[m.cursor].Completed = !m.todos[m.cursor].Completed
			}
		case "d":
			if len(m.todos) > 0 {
				m.todos = append(m.todos[:m.cursor], m.todos[m.cursor+1:]...)
				if m.cursor == len(m.todos) && m.cursor > 0 {
					m.cursor--
				}
			}
		case "n":
			m.inputMode = true
		}
	}
	return m, nil
}

func (m model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "alt+esc":
		m.inputMode = false
		m.input = ""
	case "enter":
		m.todos = append(m.todos, Todo{Title: m.input, Completed: false})
		m.input = ""
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		m.input += msg.String()
	}
	return m, nil
}

const todosFile = "todos.json"

func (m *model) Load() error {
	f, err := os.Open(todosFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.NewDecoder(f).Decode(&m.todos)
}

func (m model) Save() error {
	f, err := os.OpenFile(todosFile, os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		return err
	}
	return json.NewEncoder(f).Encode(m.todos)
}

func (m model) View() string {
	if m.quitting {
		return "Quitting the program, Bye Bye\n"
	}
	s := "Welcome to Todo List" + "\n"
	for i, todo := range m.todos {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		checked := " "
		if m.todos[i].Completed {
			checked = "x"
		}
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, todo.Title)
	}
	if m.inputMode {
		s += "\n" + "> New todo: " + m.input + "_\n"
		s += "Press enter to add, esc to cancel"
	} else {
		s += "Press space to toggle, d to delete, n to add"
	}
	return s
}

func (m model) Init() tea.Cmd {
	return nil
}
