package main

import ecs "github.com/x-hgg-x/goecs"

// Components contains references to all game components
type Components struct {
	Gopher *ecs.Component
	Sticky *ecs.Component
}

// Gopher component
type Gopher struct{}

// Sticky component
type Sticky struct{}
