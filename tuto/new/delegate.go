package main

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var delegateStyle = list.NewDefaultItemStyles()

type todoToggleMsg struct {
	index int
}

type todoDelegate struct {
	focused bool
}

func newTodoDelegate() *todoDelegate {
	return &todoDelegate{}
}

func (d *todoDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	t, ok := item.(todoItem)
	if !ok {
		return
	}

	var content string
	if t.done {
		content += "[x]"
	} else {
		content += "[ ]"
	}
	content += " - " + t.name

	if index == m.Cursor() && d.focused {
		content = delegateStyle.SelectedTitle.Render(content)
	} else {
		content = delegateStyle.NormalTitle.Render(content)
	}

	fmt.Fprint(w, content)
}

func (d *todoDelegate) Height() int {
	return 1
}

func (d *todoDelegate) Spacing() int {
	return 0
}

func (d *todoDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeySpace:
			return func() tea.Msg {
				return todoToggleMsg{index: m.Cursor()}
			}
		}
	}
	return nil
}
