package loader

import (
	"image/color"
	"reflect"

	c "github.com/x-hgg-x/goecsengine/components"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/BurntSushi/toml"
	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/x-hgg-x/goecs/v2"
	"golang.org/x/image/font"
)

// EngineComponentList is the list of engine components
type EngineComponentList struct {
	SpriteRender     *c.SpriteRender
	Transform        *c.Transform
	AnimationControl *c.AnimationControl
	Text             *c.Text
	UITransform      *c.UITransform
	MouseReactive    *c.MouseReactive
}

// EntityComponentList is a list of preloaded entities with components
type EntityComponentList struct {
	Engine []EngineComponentList
	Game   []interface{}
}

// LoadEntities creates entities with components from a TOML file
func LoadEntities(entityMetadataContent []byte, world w.World, gameComponentList []interface{}) []ecs.Entity {
	entityComponentList := EntityComponentList{
		Engine: LoadEngineComponents(entityMetadataContent, world),
		Game:   gameComponentList,
	}
	return AddEntities(world, entityComponentList)
}

// AddEntities adds entities with engine and game components
func AddEntities(world w.World, entityComponentList EntityComponentList) []ecs.Entity {
	// Create new entities and add engine components
	entities := make([]ecs.Entity, len(entityComponentList.Engine))
	for iEntity := range entityComponentList.Engine {
		entities[iEntity] = world.Manager.NewEntity()
		AddEntityComponents(entities[iEntity], world.Components.Engine, entityComponentList.Engine[iEntity])
	}

	// Add game components
	if entityComponentList.Game != nil {
		if len(entityComponentList.Game) != len(entityComponentList.Engine) {
			utils.LogFatalf("incorrect size for game component list")
		}
		for iEntity := range entities {
			AddEntityComponents(entities[iEntity], world.Components.Game, entityComponentList.Game[iEntity])
		}
	}
	return entities
}

// AddEntityComponents adds loaded components to an entity
func AddEntityComponents(entity ecs.Entity, ecsComponentList interface{}, components interface{}) ecs.Entity {
	ecv := reflect.ValueOf(ecsComponentList).Elem()
	cv := reflect.ValueOf(components)
	for iField := 0; iField < cv.NumField(); iField++ {
		if !cv.Field(iField).IsNil() {
			component := cv.Field(iField).Elem()
			value := reflect.New(reflect.TypeOf(component.Interface()))
			value.Elem().Set(component)

			ecsComponent := ecv.FieldByName(component.Type().Name()).Interface().(ecs.DataComponent)
			entity.AddComponent(ecsComponent, value.Interface())
		}
	}
	return entity
}

type engineComponentListData struct {
	SpriteRender     *spriteRenderData
	Transform        *c.Transform
	AnimationControl *animationControlData
	Text             *textData
	UITransform      *c.UITransform
	MouseReactive    *c.MouseReactive
}

type entity struct {
	Components engineComponentListData
}

type entityEngineMetadata struct {
	Entities []entity `toml:"entity"`
}

// LoadEngineComponents loads engine components from a TOML byte slice
func LoadEngineComponents(entityMetadataContent []byte, world w.World) []EngineComponentList {
	var entityEngineMetadata entityEngineMetadata
	utils.Try(toml.Decode(string(entityMetadataContent), &entityEngineMetadata))

	engineComponentList := make([]EngineComponentList, len(entityEngineMetadata.Entities))
	for iEntity, entity := range entityEngineMetadata.Entities {
		engineComponentList[iEntity] = processComponentsListData(world, entity.Components)
	}
	return engineComponentList
}

func processComponentsListData(world w.World, data engineComponentListData) EngineComponentList {
	return EngineComponentList{
		SpriteRender:     processSpriteRenderData(world, data.SpriteRender),
		Transform:        data.Transform,
		AnimationControl: processAnimationControlData(world, data),
		Text:             processTextData(world, data.Text),
		UITransform:      data.UITransform,
		MouseReactive:    data.MouseReactive,
	}
}

//
// SpriteRender
//

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
		utils.LogFatalf("fill and sprite_sheet_name fields are exclusive")
	}

	// Sprite is included in sprite sheet
	if spriteRenderData.SpriteSheetName != "" {
		// Add reference to sprite sheet
		spriteSheet, ok := (*world.Resources.SpriteSheets)[spriteRenderData.SpriteSheetName]
		if !ok {
			utils.LogFatalf("unable to find sprite sheet with name '%s'", spriteRenderData.SpriteSheetName)
		}
		return &c.SpriteRender{
			SpriteSheet:  &spriteSheet,
			SpriteNumber: spriteRenderData.SpriteNumber,
		}
	}

	// Sprite is a colored rectangle
	textureImage := ebiten.NewImage(spriteRenderData.Fill.Width, spriteRenderData.Fill.Height)
	textureImage.Fill(color.RGBA{
		R: spriteRenderData.Fill.Color[0],
		G: spriteRenderData.Fill.Color[1],
		B: spriteRenderData.Fill.Color[2],
		A: spriteRenderData.Fill.Color[3],
	})

	return &c.SpriteRender{
		SpriteSheet: &c.SpriteSheet{
			Texture: c.Texture{Image: textureImage},
			Sprites: []c.Sprite{{X: 0, Y: 0, Width: spriteRenderData.Fill.Width, Height: spriteRenderData.Fill.Height}},
		},
		SpriteNumber: 0,
	}
}

//
// AnimationControl
//

var endControlMap = map[string]c.EndControlType{
	"":       c.EndControlNormal,
	"Normal": c.EndControlNormal,
	"Stay":   c.EndControlStay,
	"Loop":   c.EndControlLoop,
}

type endControlData struct {
	Type string
}

var animationCommandMap = map[string]c.AnimationCommandType{
	"":             c.AnimationCommandStart,
	"None":         c.AnimationCommandNone,
	"Restart":      c.AnimationCommandRestart,
	"Start":        c.AnimationCommandStart,
	"StepBackward": c.AnimationCommandStepBackward,
	"StepForward":  c.AnimationCommandStepForward,
	"SetTime":      c.AnimationCommandSetTime,
	"Pause":        c.AnimationCommandPause,
	"Abort":        c.AnimationCommandAbort,
}

type animationCommandData struct {
	Type string
	Time float64
}

type animationControlData struct {
	SpriteSheetName string `toml:"sprite_sheet_name"`
	AnimationName   string `toml:"animation_name"`
	End             endControlData
	Command         animationCommandData
	RateMultiplier1 float64 `toml:"rate_multiplier_minus_1"`
}

func processAnimationControlData(world w.World, data engineComponentListData) *c.AnimationControl {
	animationControlData := data.AnimationControl
	spriteRenderData := data.SpriteRender
	if animationControlData == nil || spriteRenderData == nil {
		return nil
	}

	// Find spritesheet
	if animationControlData.SpriteSheetName != spriteRenderData.SpriteSheetName {
		utils.LogFatalf("AnimationControl and SpriteRender components don't have the same sprite sheet ('%s' vs '%s')", animationControlData.SpriteSheetName, spriteRenderData.SpriteSheetName)
	}
	spriteSheet, ok := (*world.Resources.SpriteSheets)[animationControlData.SpriteSheetName]
	if !ok {
		utils.LogFatalf("unable to find sprite sheet with name '%s'", animationControlData.SpriteSheetName)
	}

	// Find animation
	animation, ok := spriteSheet.Animations[animationControlData.AnimationName]
	if !ok {
		utils.LogFatalf("unable to find animation with name '%s'", animationControlData.AnimationName)
	}

	// Check end control
	endControl, ok := endControlMap[animationControlData.End.Type]
	if !ok {
		utils.LogFatalf("unknown end control option: '%s'", animationControlData.End.Type)
	}

	// Check animation command
	animationCommand, ok := animationCommandMap[animationControlData.Command.Type]
	if !ok {
		utils.LogFatalf("unknown animation command option: '%s'", animationControlData.Command.Type)
	}

	return &c.AnimationControl{
		Animation:      animation,
		End:            c.EndControl{Type: endControl},
		Command:        c.AnimationCommand{Type: animationCommand, Time: animationControlData.Command.Time},
		RateMultiplier: animationControlData.RateMultiplier1 + 1,
	}
}

//
// Text
//

type fontFaceOptions struct {
	Size              float64
	DPI               float64
	Hinting           string
	GlyphCacheEntries int `toml:"glyph_cache_entries"`
	SubPixelsX        int `toml:"sub_pixels_x"`
	SubPixelsY        int `toml:"sub_pixels_y"`
}

var hintingMap = map[string]font.Hinting{
	"":         font.HintingFull,
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
		utils.LogFatalf("unable to find font with name '%s'", textData.FontFace.Font)
	}

	// Check hinting
	hinting, ok := hintingMap[textData.FontFace.Options.Hinting]
	if !ok {
		utils.LogFatalf("unknown hinting option: '%s'", textData.FontFace.Options.Hinting)
	}

	options := &truetype.Options{
		Size:              textData.FontFace.Options.Size,
		DPI:               textData.FontFace.Options.DPI,
		Hinting:           hinting,
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
