package uisystem

import (
	"fmt"

	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// RenderUISystem draws text entities
func RenderUISystem(world w.World, screen *ebiten.Image) {
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		textData := world.Components.Engine.Text.Get(entity).(*c.Text)
		uiTransform := world.Components.Engine.UITransform.Get(entity).(*c.UITransform)

		// Compute dot offset
		x, y := utils.Try2(c.ComputeDotOffset(textData.Text, textData.FontFace, uiTransform.Pivot))

		// Draw text
		var offsetX, offsetY int

		screenWidth := world.Resources.ScreenDimensions.Width
		screenHeight := world.Resources.ScreenDimensions.Height

		switch uiTransform.Origin {
		case c.UITransformOriginTopLeft:
			offsetX, offsetY = 0, screenHeight
		case c.UITransformOriginTopMiddle:
			offsetX, offsetY = screenWidth/2, screenHeight
		case c.UITransformOriginTopRight:
			offsetX, offsetY = screenWidth, screenHeight
		case c.UITransformOriginMiddleLeft:
			offsetX, offsetY = 0, screenHeight/2
		case c.UITransformOriginMiddle:
			offsetX, offsetY = screenWidth/2, screenHeight/2
		case c.UITransformOriginMiddleRight:
			offsetX, offsetY = screenWidth, screenHeight/2
		case c.UITransformOriginBottomLeft:
			offsetX, offsetY = 0, 0
		case c.UITransformOriginBottomMiddle:
			offsetX, offsetY = screenWidth/2, 0
		case c.UITransformOriginBottomRight:
			offsetX, offsetY = screenWidth, 0
		case "": // UITransformOriginBottomLeft
			offsetX, offsetY = 0, 0
		default:
			utils.LogError(fmt.Errorf("unknown UI transform origin value: %s", uiTransform.Origin))
		}

		text.Draw(screen, textData.Text, textData.FontFace, uiTransform.Translation.X+offsetX-x, screenHeight-uiTransform.Translation.Y-offsetY-y, textData.Color)
	}))
}
