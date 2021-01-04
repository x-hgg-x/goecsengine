package resources

import (
	"github.com/x-hgg-x/goecsengine/components"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// Resources contains references to data not related to any entity
type Resources struct {
	ScreenDimensions *ScreenDimensions
	Controls         *Controls
	InputHandler     *InputHandler
	SpriteSheets     *map[string]components.SpriteSheet
	Fonts            *map[string]Font
	AudioContext     *audio.Context
	AudioPlayers     *map[string]*audio.Player
	Prefabs          interface{}
	Game             interface{}
}

// InitResources initializes resources
func InitResources() *Resources {
	return &Resources{Controls: &Controls{}, InputHandler: &InputHandler{}}
}
