package world

import (
	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/resources"

	ecs "github.com/x-hgg-x/goecs"
)

// World is the main ECS structure
type World struct {
	Manager    *ecs.Manager
	Components *c.Components
	Resources  *resources.Resources
}

// InitWorld initializes the world
func InitWorld(gameComponents interface{}, gameResources interface{}) World {
	manager := &ecs.Manager{}
	components := c.InitComponents(manager, gameComponents)
	resources := resources.InitResources(gameResources)

	return World{
		Manager:    manager,
		Components: components,
		Resources:  resources,
	}
}
