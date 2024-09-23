package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/kyleochata/hydrohomie/plant"
)

type NewPlant struct {
	Plants        []plant.Plant
	NewPlantInput textinput.Model
	Width, Height int
}

func (np NewPlant) Init() tea.Cmd { return nil }
func (np NewPlant) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//type switch -->> directly matches the concrete type of msg
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		np.Width, np.Height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return np.goToListPage()
		case "enter":
			if np.NewPlantInput.Value() != "" {
				np.Plants = append(np.Plants, plant.Plant{Name: np.NewPlantInput.Value()})
				np.NewPlantInput.Reset()
				return np, SavePlantsToMemory(np.Plants)
			}
		}

	}
	var cmd tea.Cmd
	if np.NewPlantInput, cmd = np.NewPlantInput.Update(msg); cmd != nil {
		return np, cmd
	}
	return np, nil
}
func (np NewPlant) View() string {
	return gloss.JoinVertical(gloss.Left, []string{
		titleStyle.Render("Add new plant"),
		np.NewPlantInputView(),
		"Esc: Main Page | Ctrl+C: Quit",
	}...)
}

var titleStyle = gloss.NewStyle().Bold(true).Foreground(gloss.Color("#00FF00"))

func (np NewPlant) NewPlantInputView() string {
	var s strings.Builder
	s.WriteString("Add new Plant\n")
	s.WriteString(np.NewPlantInput.View())
	return np.inputStyle().Render(s.String())
}
func (np NewPlant) inputStyle() gloss.Style {
	return gloss.NewStyle().
		Padding(1, 1, 1, 1).
		Width(np.Width - 4).
		Border(gloss.RoundedBorder()).
		BorderForeground(gloss.Color("#00FF00"))
}
func (np NewPlant) goToListPage() (tea.Model, tea.Cmd) {
	m := HydroApp{Plants: np.Plants}
	return m, func() tea.Msg {
		return tea.WindowSizeMsg{
			Height: np.Height,
			Width:  np.Width,
		}
	}

}
