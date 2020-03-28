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
	Game             interface{}
}

// InitResources initializes resources
func InitResources(gameResources interface{}) *Resources {
	return &Resources{Controls: &Controls{}, InputHandler: &InputHandler{}, Game: gameResources}
}
