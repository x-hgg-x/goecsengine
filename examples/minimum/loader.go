package main

import (
	"github.com/x-hgg-x/goecsengine/loader"
	w "github.com/x-hgg-x/goecsengine/world"

	ecs "github.com/x-hgg-x/goecs"
)

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataPath string, world w.World) []ecs.Entity {
	return loader.LoadEntities(entityMetadataPath, world, nil)
}
