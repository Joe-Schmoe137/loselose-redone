package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// globals
var xmax int
var ymax int

type model struct {
	xpos   int
}

func initalModel() model {
	return model{
		xpos:   0,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		//what was the actual key pressed?
		switch msg.String() {

		//exit program
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "a":
			if m.xpos > 0 {
				m.xpos--
			}

		case "right", "d":
			if m.xpos < xmax - 1 {
				m.xpos++
			}

		//case "enter", " ":
		//	_, ok := m.selected[m.xpos]
		//	if ok {
		//		delete(m.selected, m.xpos)
		//	} else {
		//		m.selected[m.xpos] = struct{}{}
		//	}
		}
	case tea.WindowSizeMsg:
		xmax = msg.Width
		ymax = msg.Height
	}
	return m, nil
}

func (m model) View() string {
	// The header
	display := "What matters most?\n\n"
	// Debug info
	display += fmt.Sprintf("DEBUG:\txpos = %v\tymax = %v\n", m.xpos, ymax)

	// Display the ship
	display += fmt.Sprintf("%*s", m.xpos + 1, "^")
	return display
}

func main() {
	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
