package loader

import (
	"fmt"

	"github.com/x-hgg-x/goecsengine/resources"
	"github.com/x-hgg-x/goecsengine/utils"

	"github.com/x-hgg-x/go-toml"
)

type controlsConfig struct {
	Controls resources.Controls `toml:"controls"`
}

// LoadControls loads controls from a TOML file
func LoadControls(controlsConfigPath string, axes []string, actions []string) (resources.Controls, resources.InputHandler) {
	var controlsConfig controlsConfig
	tree, err := toml.LoadFile(controlsConfigPath)
	utils.LogError(err)
	utils.LogError(tree.Unmarshal(&controlsConfig))

	var inputHandler resources.InputHandler
	inputHandler.Axes = make(map[string]float64)
	inputHandler.Actions = make(map[string]bool)

	// Check axes
	for _, axis := range axes {
		if _, ok := controlsConfig.Controls.Axes[axis]; !ok {
			utils.LogError(fmt.Errorf("unable to find controls for axis '%s'", axis))
		}
		inputHandler.Axes[axis] = 0
	}

	// Check actions
	for _, action := range actions {
		if _, ok := controlsConfig.Controls.Actions[action]; !ok {
			utils.LogError(fmt.Errorf("unable to find controls for action '%s'", action))
		}
		inputHandler.Actions[action] = false
	}

	return controlsConfig.Controls, inputHandler
}
