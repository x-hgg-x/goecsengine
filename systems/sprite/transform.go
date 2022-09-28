package spritesystem

import (
	"fmt"

	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// TransformSystem updates geometry matrix.
// Geometry matrix is first recentered, then scaled and rotated, and finally translated.
func TransformSystem(world w.World) {
	world.Manager.Join(world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
		sprite := world.Components.Engine.SpriteRender.Get(entity).(*c.SpriteRender)
		transform := world.Components.Engine.Transform.Get(entity).(*c.Transform)

		spriteWidth := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
		spriteHeight := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

		// Reset geometry matrix
		sprite.Options.GeoM.Reset()

		// Center sprite on top left pixel
		sprite.Options.GeoM.Translate(-spriteWidth/2, -spriteHeight/2)

		// Perform scale
		sprite.Options.GeoM.Scale(transform.Scale1.X+1, transform.Scale1.Y+1)

		// Perform rotation
		sprite.Options.GeoM.Rotate(-transform.Rotation)

		// Perform translation
		var offsetX, offsetY float64

		screenWidth := float64(world.Resources.ScreenDimensions.Width)
		screenHeight := float64(world.Resources.ScreenDimensions.Height)

		switch transform.Origin {
		case c.TransformOriginTopLeft:
			offsetX, offsetY = 0, screenHeight
		case c.TransformOriginTopMiddle:
			offsetX, offsetY = screenWidth/2, screenHeight
		case c.TransformOriginTopRight:
			offsetX, offsetY = screenWidth, screenHeight
		case c.TransformOriginMiddleLeft:
			offsetX, offsetY = 0, screenHeight/2
		case c.TransformOriginMiddle:
			offsetX, offsetY = screenWidth/2, screenHeight/2
		case c.TransformOriginMiddleRight:
			offsetX, offsetY = screenWidth, screenHeight/2
		case c.TransformOriginBottomLeft:
			offsetX, offsetY = 0, 0
		case c.TransformOriginBottomMiddle:
			offsetX, offsetY = screenWidth/2, 0
		case c.TransformOriginBottomRight:
			offsetX, offsetY = screenWidth, 0
		case "": // TransformOriginBottomLeft
			offsetX, offsetY = 0, 0
		default:
			utils.LogError(fmt.Errorf("unknown transform origin value: %s", transform.Origin))
		}

		sprite.Options.GeoM.Translate(transform.Translation.X+offsetX, screenHeight-transform.Translation.Y-offsetY)
	}))
}
