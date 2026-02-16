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
var headerOffset int = 3
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

// Creates an entityModel of type enemy at ypos 0 and a random valid x index
// Returns an entityModel of type enemy
func createEnemy() entityModel{
	return entityModel{
		entityType: enemyType,
		xpos: rand.Intn(xmax - 1),
		ypos: 0,
		symbol: "@",
	}
}

// Creates an entityModel of type laser above the player's position
// playerXpos int - The current xpos of the player
// returns: a entityModel with a laserType above the player's position
func createLaser(playerXpos int) entityModel{
	return entityModel{
		entityType: laserType,
		xpos: playerXpos,
		ypos: ymax - (headerOffset + 2) , //above player
		symbol: "|",
	}
}

// Generates a sorted 2d slice of entities by line [y pos (asending)][entity (asending xpos order)]
// returns: a sorted 2d slice of entityModels
func generateEntitiesByLine() [][]entityModel{
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

// Checks the entire entity slice to see if entering in an entity at a give position, would create a coordinate conflict
// 
func validAreaCheck(givenEntity entityModel) bool {
	var validArea bool = true
	for entityIndex := range entitySlice {
		if (givenEntity.xpos == entitySlice[entityIndex].xpos && givenEntity.ypos == entitySlice[entityIndex].ypos) {
			validArea = false
			break
		}
	}
	return validArea
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
			var generatedEnemy entityModel = createEnemy()
			// check for space conflicts
			for !validAreaCheck(generatedEnemy) {
				generatedEnemy = createEnemy()
			}
			entitySlice = append(entitySlice, generatedEnemy)
		case "enter", " ":
			var generatedLaser entityModel = createLaser(m.xpos)
			// check for space conflicts
			if validAreaCheck(generatedLaser) {
				entitySlice = append(entitySlice, generatedLaser)
			}
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
	// The header
	display := "What matters most?\n\n"
	var yEntities = generateEntitiesByLine()
	// Debug info
	if (len(yEntities) > 0 && len(yEntities[0]) > 0) {
		laser := yEntities[0][len(yEntities[0]) - 1]
		display += fmt.Sprintf("DEBUG: entity.xpos = %v entity.ypos = %v\n", laser.xpos, laser.ypos)
	} else {
		display += fmt.Sprintf("DEBUG:\txpos = %v\tymax = %v\theaderOffset = %v\n", m.xpos, ymax, headerOffset)
	}
	// Print "space"
	// create line by line slice of entities
	for ySliceIndex := range yEntities{
		var previousXPos int = 0
		for entityIndex := range yEntities[ySliceIndex] {
			entity := yEntities[ySliceIndex][entityIndex]
			//display += fmt.Sprintf("%v ", entity.symbol)
			display += fmt.Sprintf("%*s", (entity.xpos - previousXPos), entity.symbol)
			previousXPos = entity.xpos
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
