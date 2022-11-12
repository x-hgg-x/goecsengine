package components

import (
	"fmt"
	"image/color"

	"github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"

	"golang.org/x/image/font"
)

// Text component
type Text struct {
	ID       string
	Text     string
	FontFace font.Face
	Color    color.RGBA
}

// Pivot variants
const (
	PivotDot          = "Dot"
	PivotTopLeft      = "TopLeft"
	PivotTopMiddle    = "TopMiddle"
	PivotTopRight     = "TopRight"
	PivotMiddleLeft   = "MiddleLeft"
	PivotMiddle       = "Middle"
	PivotMiddleRight  = "MiddleRight"
	PivotBottomLeft   = "BottomLeft"
	PivotBottomMiddle = "BottomMiddle"
	PivotBottomRight  = "BottomRight"
)

// UI transform origin variants
const (
	UITransformOriginTopLeft      = "TopLeft"
	UITransformOriginTopMiddle    = "TopMiddle"
	UITransformOriginTopRight     = "TopRight"
	UITransformOriginMiddleLeft   = "MiddleLeft"
	UITransformOriginMiddle       = "Middle"
	UITransformOriginMiddleRight  = "MiddleRight"
	UITransformOriginBottomLeft   = "BottomLeft"
	UITransformOriginBottomMiddle = "BottomMiddle"
	UITransformOriginBottomRight  = "BottomRight"
)

// UITransform component
type UITransform struct {
	// Translation defines the position of the pivot relative to the origin.
	Translation math.VectorInt2
	// Origin defines the origin (0, 0) relative to the screen. Default is "BottomLeft".
	Origin string
	// Pivot defines the position of the element relative to its translation (default is Middle).
	Pivot string
}

// ComputeOriginOffset returns the UI transform origin offset
func (t *UITransform) ComputeOriginOffset(screenWidth, screenHeight int) (offsetX, offsetY int) {
	switch t.Origin {
	case UITransformOriginTopLeft:
		offsetX, offsetY = 0, screenHeight
	case UITransformOriginTopMiddle:
		offsetX, offsetY = screenWidth/2, screenHeight
	case UITransformOriginTopRight:
		offsetX, offsetY = screenWidth, screenHeight
	case UITransformOriginMiddleLeft:
		offsetX, offsetY = 0, screenHeight/2
	case UITransformOriginMiddle:
		offsetX, offsetY = screenWidth/2, screenHeight/2
	case UITransformOriginMiddleRight:
		offsetX, offsetY = screenWidth, screenHeight/2
	case UITransformOriginBottomLeft:
		offsetX, offsetY = 0, 0
	case UITransformOriginBottomMiddle:
		offsetX, offsetY = screenWidth/2, 0
	case UITransformOriginBottomRight:
		offsetX, offsetY = screenWidth, 0
	case "": // UITransformOriginBottomLeft
		offsetX, offsetY = 0, 0
	default:
		utils.LogError(fmt.Errorf("unknown UI transform origin value: %s", t.Origin))
	}
	return
}

// MouseReactive component
type MouseReactive struct {
	ID          string
	Hovered     bool
	JustClicked bool
}

// ComputeDotOffset computes dot offset from text and pivot
func ComputeDotOffset(text string, fontFace font.Face, pivot string) (x, y int, err error) {
	bounds, _ := font.BoundString(fontFace, text)
	centerX := ((bounds.Min.X + bounds.Max.X) / 2).Round()
	centerY := ((bounds.Min.Y + bounds.Max.Y) / 2).Round()

	switch pivot {
	case PivotDot:
		x, y = 0, 0
	case PivotTopLeft:
		x, y = bounds.Min.X.Floor(), bounds.Min.Y.Floor()
	case PivotTopMiddle:
		x, y = centerX, bounds.Min.Y.Floor()
	case PivotTopRight:
		x, y = bounds.Max.X.Ceil(), bounds.Min.Y.Floor()
	case PivotMiddleLeft:
		x, y = bounds.Min.X.Floor(), centerY
	case PivotMiddle:
		x, y = centerX, centerY
	case PivotMiddleRight:
		x, y = bounds.Max.X.Ceil(), centerY
	case PivotBottomLeft:
		x, y = bounds.Min.X.Floor(), bounds.Max.Y.Ceil()
	case PivotBottomMiddle:
		x, y = centerX, bounds.Max.Y.Ceil()
	case PivotBottomRight:
		x, y = bounds.Max.X.Ceil(), bounds.Max.Y.Ceil()
	case "": // PivotMiddle
		x, y = centerX, centerY
	default:
		err = fmt.Errorf("unknown pivot value: %s", pivot)
	}
	return
}
