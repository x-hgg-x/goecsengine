package uisystem

import (
	c "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// UISystem sets mouse reactive components
func UISystem(world w.World) {
	world.Manager.Join(world.Components.Engine.SpriteRender, world.Components.Engine.Transform, world.Components.Engine.MouseReactive).Visit(ecs.Visit(func(entity ecs.Entity) {
		sprite := world.Components.Engine.SpriteRender.Get(entity).(*c.SpriteRender)
		transform := world.Components.Engine.Transform.Get(entity).(*c.Transform)
		mouseReactive := world.Components.Engine.MouseReactive.Get(entity).(*c.MouseReactive)

		screenHeight := float64(world.Resources.ScreenDimensions.Height)

		spriteWidth := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Width)
		spriteHeight := float64(sprite.SpriteSheet.Sprites[sprite.SpriteNumber].Height)

		minX := transform.Translation.X - spriteWidth/2
		maxX := transform.Translation.X + spriteWidth/2
		minY := screenHeight - transform.Translation.Y - spriteHeight/2
		maxY := screenHeight - transform.Translation.Y + spriteHeight/2

		x, y := ebiten.CursorPosition()

		mouseReactive.Hovered = minX <= float64(x) && float64(x) <= maxX && minY <= float64(y) && float64(y) <= maxY
		mouseReactive.JustClicked = mouseReactive.Hovered && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	}))
}
