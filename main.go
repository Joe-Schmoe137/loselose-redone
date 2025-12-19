package main

import (
	"fmt"
	"os"
	"math/rand"

	tea "github.com/charmbracelet/bubbletea"
)

//===================================
// Globals
//=================================== 
var xmax int
var ymax int
var enemySlice = []entityModel{}
var laserSlice = []entityModel{}

const (
	enemyType int	= 0
	laserType		= 1
)

type model struct {
	xpos   int
}

type entityModel struct {
	entityType	int
	xpos		int
	ypos		int
	symbol		string
}

type enemyModel struct {
	xpos	int
	ypos	int
}

type laserModel struct {
	xpos	int
	ypos	int
}

//===================================
// Helper Functions
//===================================

func spawnEnemy() entityModel{
	return entityModel{
		entityType: enemyType,
		xpos: rand.Intn(xmax - 1),
		ypos: 3, // TODO: change this since it's hard coded
		symbol: "@",
	}
}

func spawnLaser(playerXpos int) entityModel{
	return entityModel{
		entityType: laserType,
		xpos: playerXpos,
		ypos: ymax - 2,
		symbol: "|",
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
		case "enter", " ":
			laserSlice = append(laserSlice, spawnLaser(m.xpos))
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
	var headerOffset int = 1
	// The header
	display := "What matters most?\n\n"
	headerOffset += 2
	// Debug info
	if len(laserSlice) >= 1 {
		laser := laserSlice[len(laserSlice) - 1]
		display += fmt.Sprintf("DEBUG: laser.xpos = %v\t laser.ypos = %v\n", laser.xpos, laser.ypos)
	} else {
		display += fmt.Sprintf("DEBUG:\txpos = %v\tymax = %v\n", m.xpos, ymax)
	}
	headerOffset++
	// Print "space"
	for i:=headerOffset; i < ymax; i++ {
		// check for enemies on each line
		for enemyIndex:=0; enemyIndex < len(enemySlice); enemyIndex++ {
			enemy := enemySlice[enemyIndex]
			// Later we will grab all info on line and put it into a Slice before checking this
			if enemy.ypos == i - 1 {
				display += fmt.Sprintf("%*s", enemy.xpos, enemy.symbol)
			}
		}
		// check for lasers on each line
		for laserIndex:=0; laserIndex < len(laserSlice); laserIndex++ {
			laser := laserSlice[laserIndex]
			
			if laser.ypos == i - 1 {
				display += fmt.Sprintf("%*s", laser.xpos + 1, laser.symbol)
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
