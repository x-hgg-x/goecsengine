package main

import (
	"os"

	"github.com/x-hgg-x/goecsengine/loader"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	ecs "github.com/x-hgg-x/goecs/v2"
)

type gameComponentList struct {
	Gopher *Gopher
	Sticky *Sticky
}

type entity struct {
	Components gameComponentList
}

type entityGameMetadata struct {
	Entities []entity `toml:"entity"`
}

func loadGameComponents(entityMetadataContent []byte, world w.World) []interface{} {
	var entityGameMetadata entityGameMetadata
	utils.Try(toml.Decode(string(entityMetadataContent), &entityGameMetadata))

	gameComponentList := make([]interface{}, len(entityGameMetadata.Entities))
	for iEntity, entity := range entityGameMetadata.Entities {
		gameComponentList[iEntity] = entity.Components
	}
	return gameComponentList
}

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataPath string, world w.World) []ecs.Entity {
	entityMetadataContent := utils.Try(os.ReadFile(entityMetadataPath))
	gameComponentList := loadGameComponents(entityMetadataContent, world)
	return loader.LoadEntities(entityMetadataContent, world, gameComponentList)
}
