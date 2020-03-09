package spritesystem

import (
	c "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	ecs "github.com/x-hgg-x/goecs"
)

// TransformSystem updates geometry matrix.
// Geometry matrix is first recentered, then scaled and rotated, and finally translated.
func TransformSystem(world w.World) {
	ecs.Join(world.Components.Engine.SpriteRender, world.Components.Engine.Transform).Visit(ecs.Visit(func(entity ecs.Entity) {
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
		screenHeight := float64(world.Resources.ScreenDimensions.Height)
		sprite.Options.GeoM.Translate(transform.Translation.X, screenHeight-transform.Translation.Y)
	}))
}
