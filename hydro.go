package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	Plants []plant.Plant
	cursor int
}

func (ha HydroApp) Init() tea.Cmd {
	return nil
}

func (ha HydroApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return ha, tea.Quit
		case "up":
			if ha.cursor <= 0 {
				ha.cursor = len(ha.Plants) - 1
			} else {
				ha.cursor--
			}
		case "down":
			if ha.cursor >= len(ha.Plants)-1 {
				ha.cursor = 0
			} else {
				ha.cursor++
			}

		}

	}
	return ha, nil
}

func (ha HydroApp) View() string {
	var s strings.Builder
	s.WriteString("\u2191 or k to move up\t \u2193 or j to move down\n")
	for _, plant := range ha.Plants {
		s.WriteString(ha.plantView(plant) + "\n\n")
	}
	s.WriteString("Press Ctrl+C or Esc to quit")
	return s.String()
}

func (ha HydroApp) plantView(plant plant.Plant) string {
	return fmt.Sprintf(
		"	%s\n	%s", plant.Name, ha.plantLastWatered(plant),
	)
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
