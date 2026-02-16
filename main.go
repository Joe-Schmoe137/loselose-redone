package main

import (
	"fmt"
	"os"
	"math/rand"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

//===================================
// Globals
//=================================== 
var xmax int
var ymax int
var entitySlice = []entityModel{}
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

// Remove
type enemyModel struct {
	xpos	int
	ypos	int
}

// Remove
type laserModel struct {
	xpos	int
	ypos	int
}

//===================================
// Helper Functions
//===================================

// Returns an entityModel of type enemy
func createEnemy() entityModel{
	return entityModel{
		entityType: enemyType,
		xpos: rand.Intn(xmax - 1),
		ypos: 0,
		symbol: "@",
	}
}

// Returns an entityModel of type laser above the player's position
// playerXpos int - The current xpos of the player
func createLaser(playerXpos int) entityModel{
	return entityModel{
		entityType: laserType,
		xpos: playerXpos,
		ypos: ymax - 2, //above player
		symbol: "|",
	}
}

// Generates a sorted 2d slice of entities by line [y pos (asending)][entity (asending xpos order)]
// headerOffset int - Spaces to skip off the top of the terminal for other text
func generateEntitiesByLine(headerOffset int) [][]entityModel{
	// Error handling
	if (ymax - (headerOffset + 1)) < 0 {
		return make([][]entityModel, 0)
	}

	// generate empty slice
	retSlice := make([][]entityModel, ymax - (headerOffset + 1))
	for i := range retSlice {
		retSlice[i] = []entityModel{}
	}

	// enter in entity to correct slice
	for entityIndex := range entitySlice {
		var entity = entitySlice[entityIndex]
		// if slice is empty, it's sorted. Enter entity
		var currentLineSlice []entityModel = retSlice[entity.ypos]
		if len(currentLineSlice) == 0 {			
			currentLineSlice = append(retSlice[entity.ypos], entity)
		} else {
			// if current entities' xpos is smaller than comparison's, enter before it. If not entered, append at end.
			var inserted bool = false
			for comparisonIndex:=0; comparisonIndex < len(currentLineSlice) - 1; comparisonIndex++ {
				if (entity.xpos <= currentLineSlice[comparisonIndex].xpos) {
					currentLineSlice = slices.Insert(currentLineSlice, comparisonIndex, entity)
					inserted = true
					break;
				}
			}
			if(!inserted) {
				currentLineSlice = append(retSlice[entity.ypos], entity)
			}
		}
		for comparisonIndex := range retSlice[entity.ypos] {
			// reached end of Slice, new entitiy has highest x value
			if(comparisonIndex == len(retSlice) - 1) {
				retSlice[entity.ypos] = append(retSlice[entity.ypos], entity)
			}
		}
		retSlice[entity.ypos] = currentLineSlice
	}
	return retSlice
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
			entitySlice = append(entitySlice, createEnemy())
		case "enter", " ":
			//spaceAvailable := true
			//for laserIndex:=0; laserIndex < len(entitySlice); laserIndex++ {
			//	laserCmp := entitySlice[laserIndex]
			//	if(laserCmp.ypos == ymax - 2 && laserCmp.xpos == m.xpos) {
			//		spaceAvailable = false
			//		break
			//	}
			//}
			//if(spaceAvailable) {
			//	entitySlice = append(entitySlice, createLaser(m.xpos))			
			//}
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
	var yEntities = generateEntitiesByLine(headerOffset)
	// Debug info
	if len(entitySlice) >= 1 {
		laser := entitySlice[len(entitySlice) - 1]
		display += fmt.Sprintf("DEBUG: laser.xpos = %v\t laser.ypos = %v\t len(yEntities[0]) = %v\n", laser.xpos, laser.ypos, len(yEntities[0]))
	} else {
		display += fmt.Sprintf("DEBUG:\txpos = %v\tymax = %v\theaderOffset = %v\tlen(yEntities) = %v\n", m.xpos, ymax, headerOffset, len(yEntities))
	}
	headerOffset++
	// Print "space"
	// create line by line slice of entities
	for ySliceIndex := range yEntities{
		for entityIndex := range yEntities[ySliceIndex] {
			entity := yEntities[ySliceIndex][entityIndex]
			display += fmt.Sprintf("%*s", entity.xpos + 1, entity.symbol)
		}
		display += fmt.Sprintf("\n")
	}
	//for i:=headerOffset; i < ymax; i++ {
	//	// check for enemies on each line
	//	for enemyIndex:=0; enemyIndex < len(entitySlice); enemyIndex++ {
	//		enemy := entitySlice[enemyIndex]
	//		// Later we will grab all info on line and put it into a Slice before checking this
	//		if enemy.ypos == i - 1 {
	//			display += fmt.Sprintf("%*s", enemy.xpos, enemy.symbol)
	//		}
	//	}
	//	// check for lasers on each line
	//	for laserIndex:=0; laserIndex < len(entitySlice); laserIndex++ {
	//		laser := entitySlice[laserIndex]
	//		
	//		if laser.ypos == i - 1 {
	//			display += fmt.Sprintf("%*s", laser.xpos + 1, laser.symbol)
	//		}
	//	}
	//	display += fmt.Sprintf("\n")
	//}
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
