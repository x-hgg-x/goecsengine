package main

import (
	"os"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataPath string, world w.World) []ecs.Entity {
	entityMetadataContent := utils.Try(os.ReadFile(entityMetadataPath))
	return loader.LoadEntities(entityMetadataContent, world, nil)
}
