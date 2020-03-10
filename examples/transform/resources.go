package main

const (
	// RotationAxis is the axis for rotating sprite
	RotationAxis = "Rotation"
	// DepthAxis is the axis for modifying sprite depth
	DepthAxis = "Depth"
	// AddEntityAction is the action for adding new entity
	AddEntityAction = "AddEntity"
	// DeleteEntityAction is the action for deleting an entity
	DeleteEntityAction = "DeleteEntity"
)

// Game contains game resources
type Game struct {
	Rotation float64
	Depth    float64
}

// NewGame creates a new game
func NewGame() *Game {
	return &Game{Depth: 0.25}
}
