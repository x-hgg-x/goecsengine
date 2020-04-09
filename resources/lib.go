package resources

import (
	"github.com/x-hgg-x/goecsengine/components"
)

// Resources contains references to data not related to any entity
type Resources struct {
	ScreenDimensions *ScreenDimensions
	Controls         *Controls
	InputHandler     *InputHandler
	SpriteSheets     *map[string]components.SpriteSheet
	Fonts            *map[string]Font
	Prefabs          interface{}
	Game             interface{}
}

// InitResources initializes resources
func InitResources() *Resources {
	return &Resources{Controls: &Controls{}, InputHandler: &InputHandler{}}
}
