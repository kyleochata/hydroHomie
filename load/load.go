package load

import (
	"encoding/json"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyleochata/hydrohomie/plant"
)

type PlantsLoadedMessage struct{ Plants []plant.Plant }

func LoadPlants() tea.Msg {
	memoryDir := "./memory"
	if err := os.MkdirAll(memoryDir, os.ModePerm); err != nil {
		// log err
		return nil
	}
	// Create or overwrite the plants.json file in the memory directory
	filePath := filepath.Join(memoryDir, "plants.json")
	file, err := os.Open(filePath)
	if err != nil {
		return nil
	}
	defer file.Close()
	var plants []plant.Plant
	if err := json.NewDecoder(file).Decode(&plants); err != nil {
		return nil
	}
	return PlantsLoadedMessage{Plants: plants}
}
