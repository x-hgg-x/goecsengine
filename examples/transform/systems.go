package main

import (
	"fmt"
	"math/rand"

	ecs "github.com/x-hgg-x/goecs"
	c "github.com/x-hgg-x/goecsengine/components"
	m "github.com/x-hgg-x/goecsengine/math"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// DemoSystem is a demo system
func DemoSystem(world w.World) {
	gameComponents := world.Components.Game.(*Components)
	gameResources := world.Resources.Game.(*Game)

	gameResources.Rotation += world.Resources.InputHandler.Axes[RotationAxis] * 5 / ebiten.DefaultTPS
	gameResources.Depth += world.Resources.InputHandler.Axes[DepthAxis] * 2 / ebiten.DefaultTPS

	world.Manager.Join(gameComponents.Gopher, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		transform := world.Components.Engine.Transform.Get(entity).(*c.Transform)
		transform.Rotation = gameResources.Rotation
		transform.Depth = gameResources.Depth
	}))

	// Add a gopher entity
	if world.Resources.InputHandler.Actions[AddEntityAction] {
		gopherEntity := LoadEntities("assets/gopher.toml", world)
		for iEntity := range gopherEntity {
			transform := world.Components.Engine.Transform.Get(gopherEntity[iEntity]).(*c.Transform)
			transform.Rotation = gameResources.Rotation
			transform.Depth = gameResources.Depth
			transform.Translation = m.Vector2{
				X: float64(rand.Intn(world.Resources.ScreenDimensions.Width)),
				Y: float64(rand.Intn(world.Resources.ScreenDimensions.Height)),
			}
		}
	}

	// Delete a gopher entity
	if world.Resources.InputHandler.Actions[DeleteEntityAction] {
		gophers := world.Manager.Join(gameComponents.Gopher, gameComponents.Sticky.Not())
		firstGopher := ecs.Entity(*ecs.GetFirst(gophers))
		world.Manager.DeleteEntity(firstGopher)
	}

	// Update text info
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*c.Text)
		if text.ID == "rotation" {
			text.Text = fmt.Sprintf("Gopher rotation: %.2f", gameResources.Rotation)
		}
		if text.ID == "depth" {
			text.Text = fmt.Sprintf("Gopher depth: %.2f", gameResources.Depth)
		}
	}))
}
