package loader

import (
	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"

	"github.com/x-hgg-x/go-toml"
)

type spriteSheetMetadata struct {
	SpriteSheets map[string]c.SpriteSheet `toml:"sprite_sheet"`
}

// LoadSpriteSheets loads sprite sheets from a TOML file
func LoadSpriteSheets(spriteSheetMetadataPath string) map[string]c.SpriteSheet {
	var spriteSheetMetadata spriteSheetMetadata
	tree, err := toml.LoadFile(spriteSheetMetadataPath)
	utils.LogError(err)
	utils.LogError(tree.Unmarshal(&spriteSheetMetadata))
	return spriteSheetMetadata.SpriteSheets
}
