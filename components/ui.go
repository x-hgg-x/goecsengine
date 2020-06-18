package components

import (
	"fmt"
	"image/color"

	"github.com/x-hgg-x/goecsengine/math"

	"golang.org/x/image/font"
)

// Text component
type Text struct {
	ID       string
	Text     string
	FontFace font.Face
	Color    color.RGBA
	OffsetX  int
	OffsetY  int
}

// Pivot variants
const (
	Dot          = "Dot"
	TopLeft      = "TopLeft"
	TopMiddle    = "TopMiddle"
	TopRight     = "TopRight"
	MiddleLeft   = "MiddleLeft"
	Middle       = "Middle"
	MiddleRight  = "MiddleRight"
	BottomLeft   = "BottomLeft"
	BottomMiddle = "BottomMiddle"
	BottomRight  = "BottomRight"
)

// UITransform component
type UITransform struct {
	// Translation defines the position of the pivot relative to the origin.
	Translation math.VectorInt2
	// Pivot defines the position of the element relative to its translation (default is Middle).
	Pivot string
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
	case Dot:
		x, y = 0, 0
	case TopLeft:
		x, y = bounds.Min.X.Floor(), bounds.Min.Y.Floor()
	case TopMiddle:
		x, y = centerX, bounds.Min.Y.Floor()
	case TopRight:
		x, y = bounds.Max.X.Ceil(), bounds.Min.Y.Floor()
	case MiddleLeft:
		x, y = bounds.Min.X.Floor(), centerY
	case Middle:
		x, y = centerX, centerY
	case MiddleRight:
		x, y = bounds.Max.X.Ceil(), centerY
	case BottomLeft:
		x, y = bounds.Min.X.Floor(), bounds.Max.Y.Ceil()
	case BottomMiddle:
		x, y = centerX, bounds.Max.Y.Ceil()
	case BottomRight:
		x, y = bounds.Max.X.Ceil(), bounds.Max.Y.Ceil()
	case "": // Middle
		x, y = centerX, centerY
	default:
		err = fmt.Errorf("unknown pivot value: %s", pivot)
	}
	return
}
