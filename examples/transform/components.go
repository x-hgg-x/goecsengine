package main

import ecs "github.com/x-hgg-x/goecs/v2"

// Components contains references to all game components
type Components struct {
	Gopher *ecs.NullComponent
	Sticky *ecs.NullComponent
}

// Gopher component
type Gopher struct{}

// Sticky component
type Sticky struct{}
