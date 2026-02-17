package main

import (
	"fmt"
	"os"
	"math/rand"
	"sort"	
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

//===================================
// Globals
//=================================== 
var xmax int
var ymax int
var debug bool = false
var headerOffset int = 3
var entitySlice = []entityModel{}
var laserSlice = []entityModel{}

const (
	playerType 	int = 0
	enemyType 	int	= 1
	laserType 	int	= 2
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
// player entityModel - The current entityModel of the player
// returns: a entityModel with a laserType above the player's position
func createLaser(player entityModel) entityModel{
	return entityModel{
		entityType: laserType,
		xpos: player.xpos,
		ypos: player.ypos - 1 , //above player
		symbol: "|",
	}
}

// Generates a sorted 2d slice of entities by line [y pos (asending)][entity (asending xpos order)]
// playerEntity entityModel - the player's entity that has been excluded from the entitySlice
// returns: a sorted 2d slice of entityModels
func generateEntitiesByLine(playerEntity entityModel) [][]entityModel{
	// Error handling
	if (ymax - (headerOffset - 1)) < 0 {
		return make([][]entityModel, 0)
	}

	// generate empty slice
	retSlice := make([][]entityModel, ymax - (headerOffset - 1))
	for i := range retSlice {
		retSlice[i] = []entityModel{}
	}

	// enter in entity to correct slice
	var modEntitySlice []entityModel = append(entitySlice, playerEntity) 
	for entityIndex := range modEntitySlice {
		var currentEntity = modEntitySlice[entityIndex]
		retSlice[currentEntity.ypos] = append(retSlice[currentEntity.ypos], currentEntity) 
	}
	// sort each slice
	for yIndex := range retSlice {
		sort.Slice(retSlice[yIndex], func(p, q int) bool {
			return retSlice[yIndex][p].xpos < retSlice[yIndex][q].xpos
		})
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

func initalModel() entityModel {
	return entityModel{
		entityType:playerType,
		xpos:0,	
		ypos:0,
		symbol:"^",
	}
}

func (m entityModel) Init() tea.Cmd {
	if(debug) {headerOffset = 4}
	return nil
}

func (m entityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			var generatedLaser entityModel = createLaser(m)
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
		if (ymax - headerOffset >= 0 && m.ypos != ymax - (headerOffset)) {
			m.ypos = ymax - (headerOffset)
		}
	}
	return m, nil
}

func (m entityModel) View() string {
	// The header
	display := "What matters most?\n\n"
	var yEntities = generateEntitiesByLine(m)
	// Debug info
	if(debug) {
		if (len(yEntities) > 0 && len(entitySlice) - 1 > 0) {
			laser := entitySlice[len(entitySlice) - 1]
			display += fmt.Sprintf("DEBUG: entity.xpos = %v entity.ypos = %v player.xpos = %v player.ypos = %v\n", laser.xpos, laser.ypos, m.xpos, m.ypos)
		} else {
			display += fmt.Sprintf("DEBUG:\txpos = %v\tymax = %v\theaderOffset = %v\t len(yEntities) = %v\n", m.xpos, ymax, headerOffset, len(yEntities))
		}
	}
	// Print "space"
	// create line by line slice of entities
	for ySliceIndex := range yEntities{
		var previousXPos int = 0
		for entityIndex := range yEntities[ySliceIndex] {
			entity := yEntities[ySliceIndex][entityIndex]
			var xPosDelta int = (entity.xpos - previousXPos)
			if (xPosDelta < 0) {
				xPosDelta = 0
			}
			display += strings.Repeat(" ", xPosDelta)
			display += fmt.Sprintf("%v", entity.symbol)
			previousXPos = entity.xpos + 1
		}
		if (ySliceIndex != len(yEntities) - 1) {
			display += fmt.Sprintf("\n")
		}
	}
	return display
}

func main() {
	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
