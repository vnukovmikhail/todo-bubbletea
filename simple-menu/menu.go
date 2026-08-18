package main

import tea "github.com/charmbracelet/bubbletea"

type Menu struct {
	choices []string
	cursor  int
}

var _ tea.Model = Menu{}

func NewMenu(choices []string) Menu {
	return Menu{choices: choices}
}
