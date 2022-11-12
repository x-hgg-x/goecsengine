package components

import (
	"fmt"

	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Sprite structure
type Sprite struct {
	// Horizontal position of the sprite in the sprite sheet
	X int
	// Vertical position of the sprite in the sprite sheet
	Y int
	// Width of the sprite
	Width int
	// Height of the sprite
	Height int
}

// Texture structure
type Texture struct {
	// Texture image
	Image *ebiten.Image
}

// UnmarshalText fills structure fields from text data
func (t *Texture) UnmarshalText(text []byte) error {
	textureImage, _ := utils.Try2(ebitenutil.NewImageFromFile(string(text)))
	t.Image = textureImage
	return nil
}

// SpriteSheet structure
type SpriteSheet struct {
	// Texture image
	Texture Texture `toml:"texture_image"`
	// List of sprites
	Sprites []Sprite
	// List of animations
	Animations map[string]*Animation
}

// SpriteRender component
type SpriteRender struct {
	// Reference sprite sheet
	SpriteSheet *SpriteSheet
	// Index of the sprite on the sprite sheet
	SpriteNumber int
	// Draw options
	Options ebiten.DrawImageOptions
}

// Transform origin variants
const (
	TransformOriginTopLeft      = "TopLeft"
	TransformOriginTopMiddle    = "TopMiddle"
	TransformOriginTopRight     = "TopRight"
	TransformOriginMiddleLeft   = "MiddleLeft"
	TransformOriginMiddle       = "Middle"
	TransformOriginMiddleRight  = "MiddleRight"
	TransformOriginBottomLeft   = "BottomLeft"
	TransformOriginBottomMiddle = "BottomMiddle"
	TransformOriginBottomRight  = "BottomRight"
)

// Transform component.
// The origin (0, 0) is the lower left part of screen.
// Image is first rotated, then scaled, and finally translated.
type Transform struct {
	// Scale1 vector defines image scaling. Contains scale value minus 1 so that zero value is identity.
	Scale1 math.Vector2 `toml:"scale_minus_1"`
	// Rotation angle is measured counterclockwise.
	Rotation float64
	// Translation defines the position of the image center relative to the origin.
	Translation math.Vector2
	// Origin defines the origin (0, 0) relative to the screen. Default is "BottomLeft".
	Origin string
	// Depth determines the drawing order on the screen. Images with higher depth are drawn above others.
	Depth float64
}

// NewTransform creates a new default transform, corresponding to identity.
func NewTransform() *Transform {
	return &Transform{}
}

// SetScale sets transform scale.
func (t *Transform) SetScale(sx, sy float64) *Transform {
	t.Scale1.X = sx - 1
	t.Scale1.Y = sy - 1
	return t
}

// SetRotation sets transform rotation.
func (t *Transform) SetRotation(angle float64) *Transform {
	t.Rotation = angle
	return t
}

// SetTranslation sets transform translation.
func (t *Transform) SetTranslation(tx, ty float64) *Transform {
	t.Translation.X = tx
	t.Translation.Y = ty
	return t
}

// SetDepth sets transform depth.
func (t *Transform) SetDepth(depth float64) *Transform {
	t.Depth = depth
	return t
}

// SetOrigin sets transform origin.
func (t *Transform) SetOrigin(origin string) *Transform {
	t.Origin = origin
	return t
}

// ComputeOriginOffset returns the transform origin offset.
func (t *Transform) ComputeOriginOffset(screenWidth, screenHeight float64) (offsetX, offsetY float64) {
	switch t.Origin {
	case TransformOriginTopLeft:
		offsetX, offsetY = 0, screenHeight
	case TransformOriginTopMiddle:
		offsetX, offsetY = screenWidth/2, screenHeight
	case TransformOriginTopRight:
		offsetX, offsetY = screenWidth, screenHeight
	case TransformOriginMiddleLeft:
		offsetX, offsetY = 0, screenHeight/2
	case TransformOriginMiddle:
		offsetX, offsetY = screenWidth/2, screenHeight/2
	case TransformOriginMiddleRight:
		offsetX, offsetY = screenWidth, screenHeight/2
	case TransformOriginBottomLeft:
		offsetX, offsetY = 0, 0
	case TransformOriginBottomMiddle:
		offsetX, offsetY = screenWidth/2, 0
	case TransformOriginBottomRight:
		offsetX, offsetY = screenWidth, 0
	case "": // TransformOriginBottomLeft
		offsetX, offsetY = 0, 0
	default:
		utils.LogError(fmt.Errorf("unknown transform origin value: %s", t.Origin))
	}
	return
}
