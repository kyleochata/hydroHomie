package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/kyleochata/hydrohomie/load"
	"github.com/kyleochata/hydrohomie/plant"
)

var AllPlants = []plant.Plant{{
	Name:            "Pothos",
	TimeLastWatered: time.Now(),
}, {
	Name:            "Cactus",
	TimeLastWatered: time.Now().Add(-3 * 24 * time.Hour), // watered 3 days ago
}, {
	Name:            "Green Onions",
	TimeLastWatered: time.Time{},
},
}

const (
	FocusStatePlantList FocusState = iota
	FocusStateNewPlantInput
)

type FocusState int

type HydroApp struct {
	Plants        []plant.Plant
	SelectedPlant int
	Width, Height int
	NewPlantInput textinput.Model
	FocusState    FocusState
}

func (ha HydroApp) Init() tea.Cmd {
	return load.LoadPlants
}

func (ha HydroApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// getting a msg with the type of tea.Msg but that interface type can be literally anything. with the type switch, I'm i am getting that msg that was passed from the parent and I am examining it's concrete type at runtime or the instance that this update is called. Once it is in it's concrete type I then run to the correct case. Within the case, i know that msg is already filling the interface of the case so i can access the fields without having to do another type assertion
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		ha.Width, ha.Height = msg.Width, msg.Height
		return ha, nil
	case load.PlantsLoadedMessage:
		ha.Plants = msg.Plants
		if len(ha.Plants) == 0 {
			ha.FocusState = FocusStateNewPlantInput
		}
		return ha, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return ha, tea.Quit
		case "up":
			if ha.SelectedPlant > 0 && ha.FocusState == FocusStatePlantList {
				ha.SelectedPlant--
			} else if ha.FocusState == FocusStatePlantList {
				ha.SelectedPlant = len(ha.Plants) - 1
			}
		case "down":
			if ha.SelectedPlant < len(ha.Plants)-1 && ha.FocusState == FocusStatePlantList {
				ha.SelectedPlant++
			} else if ha.FocusState == FocusStatePlantList {
				ha.SelectedPlant = 0
			}
		case "enter":
			if ha.FocusState == FocusStatePlantList {
				if len(ha.Plants) > 0 {
					ha.Plants[ha.SelectedPlant].WaterMe()
				}
			} else if ha.FocusState == FocusStateNewPlantInput && ha.NewPlantInput.Value() != "" {
				ha.Plants = append(ha.Plants, plant.Plant{Name: ha.NewPlantInput.Value()})
				ha.NewPlantInput.Reset()
			}
		case "ctrl+s":
			return ha, SavePlantsToMemory(ha.Plants)
		case "tab", "shift+tab":
			ha.toggleFocusState()
			//nav to add plant page
		case "ctrl+a":
			//return the new model, first cmd of the new model is to return a tea.Msg of type tea.WindowSize with the height and width k-v pair
			return newAddPlantModel(ha.Plants), func() tea.Msg {
				return tea.WindowSizeMsg{
					Height: ha.Height,
					Width:  ha.Width,
				}
			}
		}
	}
	//Render int the changes for the text input when we type
	var cmd tea.Cmd
	if ha.NewPlantInput, cmd = ha.NewPlantInput.Update(msg); cmd != nil {
		return ha, cmd
	}
	return ha, nil
}
func (ha HydroApp) View() string {
	var uiEl = []string{"\u2191 or k: move up\t \u2193 or j to move down\n"}
	for i, p := range ha.Plants {
		uiEl = append(uiEl, ha.plantView(p, i))
	}
	uiEl = append(uiEl, ha.NewPlantInputView())
	uiEl = append(uiEl, "Ctrl+A: New Plant Page")
	uiEl = append(uiEl, "Ctrl+S: Save | Ctrl+C or Esc: Quit")
	return gloss.JoinVertical(gloss.Left, uiEl...)
}
func newAddPlantModel(ha []plant.Plant) tea.Model {
	NewPlantInput := textinput.New()
	NewPlantInput.Focus()
	return NewPlant{Plants: ha, NewPlantInput: NewPlantInput}
}
func (ha *HydroApp) toggleFocusState() {
	switch ha.FocusState {
	case FocusStateNewPlantInput:
		if len(ha.Plants) > 0 {
			ha.FocusState = FocusStatePlantList
			ha.NewPlantInput.Blur()
		}
	case FocusStatePlantList:
		ha.FocusState = FocusStateNewPlantInput
		ha.NewPlantInput.Focus()
	}
}

func (ha HydroApp) plantView(plant plant.Plant, index int) string {
	plantText := ha.plantText(plant, index)
	if index == ha.SelectedPlant && ha.FocusState == FocusStatePlantList {
		return ha.boxStyleSelected().Render(plantText)
	}
	return ha.boxStyle().Render(plantText)
}

func (ha HydroApp) plantText(plant plant.Plant, index int) string {
	s := "%s\n%s"
	if index == ha.SelectedPlant && ha.FocusState == FocusStatePlantList {
		s = "🚰 %s\n%s"
	}
	return fmt.Sprintf(s, plant.Name, ha.plantLastWatered(plant))
}
func (ha HydroApp) plantLastWatered(plant plant.Plant) string {
	switch day, ok := plant.LastWatered(); {
	case !ok:
		return "Not watered yet! It's a desert out here..."
	case day <= 0:
		return "Just drank water today!"
	default:
		return fmt.Sprintf("Last watered %d days ago", day)
	}
}
func SavePlantsToMemory(plants []plant.Plant) tea.Cmd {
	return func() tea.Msg {
		memoryDir := "./memory"
		// Ensure the 'memory' directory exists
		if err := os.MkdirAll(memoryDir, os.ModePerm); err != nil {
			// CHORE: look at bubbletea log to file
			return nil
		}

		filePath := filepath.Join(memoryDir, "plants.json")
		// Create or overwrite the plants.json file in the memory directory
		if f, err := os.Create(filePath); err == nil {
			defer f.Close()
			if err := json.NewEncoder(f).Encode(&plants); err != nil {
				// log
				return nil
			}
		} else {
			//log
			return nil
		}
		return nil
	}
}

func (ha HydroApp) NewPlantInputView() string {
	var s strings.Builder
	s.WriteString("Add new plant\n")
	s.WriteString(ha.NewPlantInput.View())
	if ha.FocusState == FocusStateNewPlantInput {
		return ha.boxStyleSelected().Render(s.String())
	}
	return ha.boxStyle().Render(s.String())
}

// ========================LipGloss=====================================

func (ha HydroApp) boxStyle() gloss.Style {
	return gloss.NewStyle().
		BorderForeground(gloss.Color("#FFFFFF")).
		Foreground(gloss.Color("#FFFFFF")).
		Border(gloss.RoundedBorder()).
		Width(ha.Width-4).
		Padding(0, 1, 0, 1)
}
func (ha HydroApp) boxStyleSelected() gloss.Style {
	return gloss.NewStyle().
		BorderForeground(gloss.Color("#00FF00")).
		Border(gloss.RoundedBorder()).
		Bold(true).
		Foreground(gloss.Color(""))
}
