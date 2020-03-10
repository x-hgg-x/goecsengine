package main

import (
	"github.com/x-hgg-x/goecsengine/states"
	i "github.com/x-hgg-x/goecsengine/systems/input"
	s "github.com/x-hgg-x/goecsengine/systems/sprite"
	u "github.com/x-hgg-x/goecsengine/systems/ui"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	// Load game and text entities
	LoadEntities("assets/start.toml", world)
	LoadEntities("assets/text.toml", world)

	world.Resources.Game = NewGame()
}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {
	world.Resources.Game = nil
	world.Manager.DeleteAllEntities()
}

// Update method
func (st *GameplayState) Update(world w.World, screen *ebiten.Image) states.Transition {
	//
	// Pre-game systems
	//
	// Register input commands
	i.InputSystem(world)
	// Set MouseReactive components (optional here)
	u.UISystem(world)

	//
	// Game systems
	//
	DemoSystem(world)

	//
	// Post-game systems
	//
	// Update ebiten geometry matrix
	s.TransformSystem(world)
	// Draw sprites
	s.RenderSpriteSystem(world, screen)
	// Draw text
	u.RenderUISystem(world, screen)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return states.Transition{TransType: states.TransQuit}
	}
	return states.Transition{}
}
