package main

import (
	_ "image/png"

	"github.com/x-hgg-x/goecsengine/loader"
	r "github.com/x-hgg-x/goecsengine/resources"
	s "github.com/x-hgg-x/goecsengine/states"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	windowWidth  = 600
	windowHeight = 600
)

type mainGame struct {
	world        w.World
	stateMachine s.StateMachine
}

func (game *mainGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return windowWidth, windowHeight
}

func (game *mainGame) Update() error {
	game.stateMachine.Update(game.world)
	return nil
}

func (game *mainGame) Draw(screen *ebiten.Image) {
	game.stateMachine.Draw(game.world, screen)
}

func main() {
	world := w.InitWorld(&Components{})

	// Init screen dimensions
	world.Resources.ScreenDimensions = &r.ScreenDimensions{Width: windowWidth, Height: windowHeight}

	// Load controls
	axes := []string{RotationAxis, DepthAxis}
	actions := []string{AddEntityAction, DeleteEntityAction}
	controls, inputHandler := loader.LoadControls("config/controls.toml", axes, actions)
	world.Resources.Controls = &controls
	world.Resources.InputHandler = &inputHandler

	// Load sprite sheets
	spriteSheets := loader.LoadSpriteSheets("metadata/spritesheets.toml")
	world.Resources.SpriteSheets = &spriteSheets

	// Load fonts
	fonts := loader.LoadFonts("metadata/fonts.toml")
	world.Resources.Fonts = &fonts

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Demo")

	utils.LogError(ebiten.RunGame(&mainGame{world, s.Init(&GameplayState{}, world)}))
}
