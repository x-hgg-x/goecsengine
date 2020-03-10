package loader

import (
	"fmt"
	"image/color"
	"reflect"

	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	ecs "github.com/x-hgg-x/goecs"
	"golang.org/x/image/font"
)

// EngineComponentList is the list of engine components
type EngineComponentList struct {
	SpriteRender  *c.SpriteRender
	Transform     *c.Transform
	Text          *c.Text
	UITransform   *c.UITransform
	MouseReactive *c.MouseReactive
}

type engineComponentListData struct {
	SpriteRender  *spriteRenderData
	Transform     *c.Transform
	Text          *textData
	UITransform   *c.UITransform
	MouseReactive *c.MouseReactive
}

type entity struct {
	Components engineComponentListData
}

type entityEngineMetadata struct {
	Entities []entity `toml:"entity"`
}

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataPath string, world w.World, gameComponentList []interface{}) []ecs.Entity {
	engineComponentList := LoadEngineComponents(entityMetadataPath, world)

	entities := make([]ecs.Entity, len(engineComponentList))
	for iEntity := range engineComponentList {
		// Add components to a new entity
		entities[iEntity] = world.Manager.NewEntity()
		AddEntityComponents(entities[iEntity], world.Components.Engine, engineComponentList[iEntity])
		AddEntityComponents(entities[iEntity], world.Components.Game, gameComponentList[iEntity])
	}
	return entities
}

// AddEntityComponents adds loaded components to an entity
func AddEntityComponents(entity ecs.Entity, ecsComponentList interface{}, components interface{}) ecs.Entity {
	ecv := reflect.ValueOf(ecsComponentList).Elem()
	cv := reflect.ValueOf(components)
	for iField := 0; iField < cv.NumField(); iField++ {
		if !cv.Field(iField).IsNil() {
			component := cv.Field(iField)
			componentName := component.Elem().Type().Name()
			ecsComponent := ecv.FieldByName(componentName).Interface().(*ecs.Component)
			entity.AddComponent(ecsComponent, component.Interface())
		}
	}
	return entity
}

// LoadEngineComponents loads engine components from a TOML file
func LoadEngineComponents(entityMetadataPath string, world w.World) []EngineComponentList {
	var entityEngineMetadata entityEngineMetadata
	_, err := toml.DecodeFile(entityMetadataPath, &entityEngineMetadata)
	utils.LogError(err)

	engineComponentList := make([]EngineComponentList, len(entityEngineMetadata.Entities))
	for iEntity, entity := range entityEngineMetadata.Entities {
		engineComponentList[iEntity] = processComponentsListData(world, entity.Components)
	}
	return engineComponentList
}

func processComponentsListData(world w.World, data engineComponentListData) EngineComponentList {
	return EngineComponentList{
		SpriteRender:  processSpriteRenderData(world, data.SpriteRender),
		Transform:     data.Transform,
		Text:          processTextData(world, data.Text),
		UITransform:   data.UITransform,
		MouseReactive: data.MouseReactive,
	}
}

type fillData struct {
	Width  int
	Height int
	Color  [4]uint8
}

type spriteRenderData struct {
	Fill            *fillData
	SpriteSheetName string `toml:"sprite_sheet_name"`
	SpriteNumber    int    `toml:"sprite_number"`
}

func processSpriteRenderData(world w.World, spriteRenderData *spriteRenderData) *c.SpriteRender {
	if spriteRenderData == nil {
		return nil
	}
	if spriteRenderData.Fill != nil && spriteRenderData.SpriteSheetName != "" {
		utils.LogError(fmt.Errorf("fill and sprite_sheet_name fields are exclusive"))
	}

	// Sprite is included in sprite sheet
	if spriteRenderData.SpriteSheetName != "" {
		// Add reference to sprite sheet from its name
		spriteSheet, ok := (*world.Resources.SpriteSheets)[spriteRenderData.SpriteSheetName]
		if !ok {
			utils.LogError(fmt.Errorf("unable to find sprite sheet with name '%s'", spriteRenderData.SpriteSheetName))
		}
		return &c.SpriteRender{
			SpriteSheet:  &spriteSheet,
			SpriteNumber: spriteRenderData.SpriteNumber,
		}
	}

	// Sprite is a colored rectangle
	textureImage, err := ebiten.NewImage(spriteRenderData.Fill.Width, spriteRenderData.Fill.Height, ebiten.FilterNearest)
	utils.LogError(err)

	textureImage.Fill(color.RGBA{
		R: spriteRenderData.Fill.Color[0],
		G: spriteRenderData.Fill.Color[1],
		B: spriteRenderData.Fill.Color[2],
		A: spriteRenderData.Fill.Color[3],
	})

	return &c.SpriteRender{
		SpriteSheet: &c.SpriteSheet{
			Texture: c.Texture{Image: textureImage},
			Sprites: []c.Sprite{c.Sprite{X: 0, Y: 0, Width: spriteRenderData.Fill.Width, Height: spriteRenderData.Fill.Height}},
		},
		SpriteNumber: 0,
	}
}

type fontFaceOptions struct {
	Size              float64
	DPI               float64
	Hinting           string
	GlyphCacheEntries int `toml:"glyph_cache_entries"`
	SubPixelsX        int `toml:"sub_pixels_x"`
	SubPixelsY        int `toml:"sub_pixels_y"`
}

var hintingMap = map[string]font.Hinting{
	"":         font.HintingNone,
	"None":     font.HintingNone,
	"Vertical": font.HintingVertical,
	"Full":     font.HintingFull,
}

type fontFaceData struct {
	Font    string
	Options fontFaceOptions
}

type textData struct {
	ID       string
	Text     string
	FontFace fontFaceData `toml:"font_face"`
	Color    [4]uint8
}

func processTextData(world w.World, textData *textData) *c.Text {
	if textData == nil {
		return nil
	}

	// Search font from its name
	textFont, ok := (*world.Resources.Fonts)[textData.FontFace.Font]
	if !ok {
		utils.LogError(fmt.Errorf("unable to find font with name '%s'", textData.FontFace.Font))
	}

	options := &truetype.Options{
		Size:              textData.FontFace.Options.Size,
		DPI:               textData.FontFace.Options.DPI,
		Hinting:           hintingMap[textData.FontFace.Options.Hinting],
		GlyphCacheEntries: textData.FontFace.Options.GlyphCacheEntries,
		SubPixelsX:        textData.FontFace.Options.SubPixelsX,
		SubPixelsY:        textData.FontFace.Options.SubPixelsY,
	}

	return &c.Text{
		ID:       textData.ID,
		Text:     textData.Text,
		FontFace: truetype.NewFace(textFont.Font, options),
		Color:    color.RGBA{R: textData.Color[0], G: textData.Color[1], B: textData.Color[2], A: textData.Color[3]},
	}
}
