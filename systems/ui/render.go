package uisystem

import (
	"fmt"

	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	ecs "github.com/x-hgg-x/goecs/v2"
	"golang.org/x/image/font"
)

// RenderUISystem draws text entities
func RenderUISystem(world w.World, screen *ebiten.Image) {
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		textData := world.Components.Engine.Text.Get(entity).(*c.Text)
		uiTransform := world.Components.Engine.UITransform.Get(entity).(*c.UITransform)

		bounds, _ := font.BoundString(textData.FontFace, textData.Text)
		centerX := ((bounds.Min.X + bounds.Max.X) / 2).Round()
		centerY := ((bounds.Min.Y + bounds.Max.Y) / 2).Round()

		// Compute dot offset from pivot
		var x, y int
		switch uiTransform.Pivot {
		case c.TopLeft:
			x, y = bounds.Min.X.Floor(), bounds.Min.Y.Floor()
		case c.TopMiddle:
			x, y = centerX, bounds.Min.Y.Floor()
		case c.TopRight:
			x, y = bounds.Max.X.Ceil(), bounds.Min.Y.Floor()
		case c.MiddleLeft:
			x, y = bounds.Min.X.Floor(), centerY
		case c.Middle:
			x, y = centerX, centerY
		case c.MiddleRight:
			x, y = bounds.Max.X.Ceil(), centerY
		case c.BottomLeft:
			x, y = bounds.Min.X.Floor(), bounds.Max.Y.Ceil()
		case c.BottomMiddle:
			x, y = centerX, bounds.Max.Y.Ceil()
		case c.BottomRight:
			x, y = bounds.Max.X.Ceil(), bounds.Max.Y.Ceil()
		case "": // Middle
			x, y = centerX, centerY
		default:
			utils.LogError(fmt.Errorf("unknown pivot value: %s", uiTransform.Pivot))
		}

		// Draw text
		screenHeight := world.Resources.ScreenDimensions.Height
		text.Draw(screen, textData.Text, textData.FontFace, uiTransform.Translation.X-x, screenHeight-uiTransform.Translation.Y-y, textData.Color)
	}))
}
