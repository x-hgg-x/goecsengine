package loader

import (
	"github.com/x-hgg-x/goecsengine/resources"
	"github.com/x-hgg-x/goecsengine/utils"

	"github.com/BurntSushi/toml"
)

type fontMetadata struct {
	Fonts map[string]resources.Font `toml:"font"`
}

// LoadFonts loads fonts from a TOML file
func LoadFonts(fontPath string) map[string]resources.Font {
	var fontMetadata fontMetadata
	_, err := toml.DecodeFile(fontPath, &fontMetadata)
	utils.LogError(err)
	return fontMetadata.Fonts
}
