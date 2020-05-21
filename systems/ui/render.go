package uisystem

import (
	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// RenderUISystem draws text entities
func RenderUISystem(world w.World, screen *ebiten.Image) {
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		textData := world.Components.Engine.Text.Get(entity).(*c.Text)
		uiTransform := world.Components.Engine.UITransform.Get(entity).(*c.UITransform)

		// Compute dot offset
		x, y, err := c.ComputeDotOffset(textData.Text, textData.FontFace, uiTransform.Pivot)
		utils.LogError(err)

		// Draw text
		screenHeight := world.Resources.ScreenDimensions.Height
		text.Draw(screen, textData.Text, textData.FontFace, uiTransform.Translation.X+textData.OffsetX-x, screenHeight-uiTransform.Translation.Y-textData.OffsetY-y, textData.Color)
	}))
}
