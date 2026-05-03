package main

import (
	"fmt"
	"os"
	"math/rand"
	"sort"	
	"strings"
	"log"
	"time"
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

//===================================
// Globals
//=================================== 
var xmax int
var ymax int
var enemyCount int = 0
var enemyCap int
var headerOffset int = 3
var tickCounter int = 0
var debug bool = false
var entitySlice = []entityModel{}

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

type TickMsg time.Time

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

// Updates entities for each tick
func processTick(tickCount int) {
	// Update entities
	var entityCeiling int = ymax - headerOffset + 1
	totalEntities := len(entitySlice)
	for entityIndex:=0; entityIndex < totalEntities; entityIndex++{
		// bandaid solution to fix race issue
		if (entityIndex >= len(entitySlice)) { continue }
		entity := entitySlice[entityIndex]
		// Check collisions
		if (entity.entityType == laserType) {
			entity.ypos = entity.ypos - 1
			killed := false
			// Check to see if laser has colided with an enemy
			for enemyCheckIndex := range entitySlice{
				checkedEntity := entitySlice[enemyCheckIndex]
				if(checkedEntity.entityType == enemyType && checkedEntity.ypos == entity.ypos && checkedEntity.xpos == entity.xpos)  {
					log.Printf("Deleting enemy at x:%v y:%v index:%v\n", checkedEntity.xpos, checkedEntity.ypos, enemyCheckIndex)
					entitySlice = slices.Delete(entitySlice, enemyCheckIndex, enemyCheckIndex + 1)
					totalEntities--
					var entityIndexAdjust int = 0
					if (enemyCheckIndex < entityIndex) {
						entityIndex--
					}
					log.Printf("entitySlice length: %v\n", len(entitySlice))
					log.Printf("Deleting laser at x:%v y:%v index:%v\n", entity.xpos, entity.ypos, entityIndex + entityIndexAdjust)
					log.Printf("slices.Delete(entitySlice (len of %v), %v, %v)\n", len(entitySlice), entityIndex + entityIndexAdjust, entityIndex + entityIndexAdjust + 1)
					entitySlice = slices.Delete(entitySlice, entityIndex, entityIndex + 1)
					totalEntities--
					killed = true
					break;
				}
			}	
			if killed { continue }
		} else if (entity.entityType == enemyType && tickCount % 5 == 0) {
			entity.ypos = entity.ypos + 1
		}
		// if entity has invalid ypos then remove it
		if (entity.ypos < 0 || entity.ypos >= entityCeiling) {
			entitySlice = slices.Delete(entitySlice, entityIndex, entityIndex + 1)
			totalEntities--
		} else {
			entitySlice[entityIndex] = entity
		}
	}
}

//===================================
// Bubble tea functions
//===================================

func doTick() tea.Cmd { 
	tickDuration := time.Millisecond * 100
	return tea.Tick(tickDuration, func(t time.Time) tea.Msg{
		return TickMsg(t)
	})
}

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
	return doTick()
}

func (m entityModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		//what was the actual key pressed?
		switch msg.String() {

		//exit program
		case "ctrl+c", "q":
			log.Print("Quitting!")
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
			// Check to see if enemy can be created
			enemyCap = xmax
			if (enemyCount + 1 >= enemyCap) {
				break
			}
			var generatedEnemy entityModel = createEnemy()
			// check for space conflicts
			for !validAreaCheck(generatedEnemy) {
				generatedEnemy = createEnemy()
			}
			enemyCount++
			log.Printf("new enemy x: %v\tnew enemy y: %v\n", generatedEnemy.xpos, generatedEnemy.ypos)
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
	// Used to register Ticks
	case TickMsg:
		tickCounter++
		processTick(tickCounter)	
		return m, doTick()
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
	// Error logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
	// setup
	p := tea.NewProgram(initalModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
