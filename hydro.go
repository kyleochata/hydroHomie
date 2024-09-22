package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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

// Capitalize 'Plants' to make it accessible from other files
type HydroApp struct {
	Plants        []plant.Plant
	SelectedPlant int
}

func (ha HydroApp) Init() tea.Cmd {
	return load.LoadPlants
}

func (ha HydroApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case load.PlantsLoadedMessage:
		ha.Plants = msg.Plants
		return ha, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return ha, tea.Quit
		case "up":
			if ha.SelectedPlant > 0 {
				ha.SelectedPlant--
			} else {
				ha.SelectedPlant = len(ha.Plants) - 1
			}
		case "down":
			if ha.SelectedPlant < len(ha.Plants)-1 {
				ha.SelectedPlant++
			} else {
				ha.SelectedPlant = 0
			}
		case "enter":
			if len(ha.Plants) > 0 {
				ha.Plants[ha.SelectedPlant].WaterMe()
			}
		case "ctrl+s":
			return ha, savePlantsToMemory(ha.Plants)

		}
	}
	return ha, nil
}
func (ha HydroApp) View() string {
	var s strings.Builder
	s.WriteString("\u2191 or k to move up\t \u2193 or j to move down\n")
	for ind, plant := range ha.Plants {
		s.WriteString(ha.plantView(plant, ind) + "\n\n")
	}
	s.WriteString("Press Ctrl+S to save changes")
	s.WriteString("Press Ctrl+C or Esc to quit")
	return s.String()
}

func (ha HydroApp) plantView(plant plant.Plant, index int) string {
	s := "	%s\n	%s"
	if index == ha.SelectedPlant {
		s = "🚰\t%s\n	%s"
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
func savePlantsToMemory(plants []plant.Plant) tea.Cmd {
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
