package main

import (
	"os"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"
)

// GameplayState is the main game state
type GameplayState struct{}

// OnPause method
func (st *GameplayState) OnPause(world w.World) {}

// OnResume method
func (st *GameplayState) OnResume(world w.World) {}

// OnStart method
func (st *GameplayState) OnStart(world w.World) {
	loader.LoadEntities(utils.Try(os.ReadFile("game.toml")), world, nil)
}

// OnStop method
func (st *GameplayState) OnStop(world w.World) {
	world.Manager.DeleteAllEntities()
}

// Update method
func (st *GameplayState) Update(world w.World) states.Transition {
	return states.Transition{}
}
