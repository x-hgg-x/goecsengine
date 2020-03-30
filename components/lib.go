package components

import (
	"reflect"

	ecs "github.com/x-hgg-x/goecs"
)

// EngineComponents contains references to all engine components
type EngineComponents struct {
	SpriteRender     *ecs.Component
	Transform        *ecs.Component
	AnimationControl *ecs.Component
	Text             *ecs.Component
	UITransform      *ecs.Component
	MouseReactive    *ecs.Component
}

// Components contains engine and game components
type Components struct {
	Engine *EngineComponents
	Game   interface{}
}

// InitComponents initializes components
func InitComponents(manager *ecs.Manager, gameComponents interface{}) *Components {
	components := &Components{Engine: &EngineComponents{}}

	ev := reflect.ValueOf(components.Engine).Elem()
	for iField := 0; iField < ev.NumField(); iField++ {
		ev.Field(iField).Set(reflect.ValueOf(manager.NewComponent()))
	}

	components.Game = gameComponents
	if gameComponents != nil {
		gv := reflect.ValueOf(components.Game).Elem()
		for iField := 0; iField < gv.NumField(); iField++ {
			gv.Field(iField).Set(reflect.ValueOf(manager.NewComponent()))
		}
	}

	return components
}
