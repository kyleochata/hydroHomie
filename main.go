package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func InitialModel() HydroApp {
	return HydroApp{Plants: AllPlants}
}

func main() {
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
