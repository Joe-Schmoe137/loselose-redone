package main

import (
	"fmt"
	"os"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

///===================================
// Globals
//===================================/ 
var xmax int
var ymax int
var enemySlice = []enemyModel{}

type model struct {
	xpos   int
}

type enemyModel struct {
	xpos	int
	ypos	int
}

//===================================
// Helper Functions
//===================================

func spawnEnemy() enemyModel {
	return enemyModel{
		xpos: rand.Intn(xmax - 1),
		ypos: 0,
	}
}

//===================================
// Bubble tea functions
//===================================

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

		case "e":
			enemySlice = append(enemySlice, spawnEnemy())
		//case "enter", " ":
		//	_, ok := m.selected[m.xpos]
		//	if ok {
		//		delete(m.selected, m.xpos)
		//	} else {
		//		m.selected[m.xpos] = struct{}{}
		//	}
		}
	// Used for window resizing
	case tea.WindowSizeMsg:
		xmax = msg.Width
		ymax = msg.Height

		// Check for out of bounds characters
		// Player
		if m.xpos >= xmax {
			m.xpos = xmax - 1
		}
	}
	return m, nil
}

func (m model) View() string {
	var headerOffset int = 0
	// The header
	display := "What matters most?\n\n"
	headerOffset += 2
	// Debug info
	display += fmt.Sprintf("DEBUG:\txpos = %v\tymax = %v\txmax = %v\n", m.xpos, ymax, xmax)
	headerOffset++
	// Print "space"
	for i:=1; i < (ymax - headerOffset); i++ {
		// check for enemies on each line
		for enemyIndex:=0; enemyIndex < len(enemySlice); enemyIndex++ {
			enemy := enemySlice[enemyIndex]
			// Later we will grab all info on line and put it into a Slice before checking this
			if enemy.ypos == i - 1 {
				display += fmt.Sprintf("%*s", enemy.xpos, "@")
			}
		}
		display += fmt.Sprintf("\n")
	}
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
