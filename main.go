package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

func InitialModel() HydroApp {
	return HydroApp{Plants: AllPlants, SelectedPlant: 0, NewPlantInput: textinput.New()}
}

func main() {
	p := tea.NewProgram(InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
