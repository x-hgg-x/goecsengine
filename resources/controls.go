package resources

import (
	"fmt"
	"reflect"

	"github.com/x-hgg-x/goecsengine/utils"

	"github.com/hajimehoshi/ebiten"
	"github.com/x-hgg-x/go-toml"
)

// Key is a US keyboard key
type Key struct {
	Key ebiten.Key
}

// UnmarshalText fills structure fields from text data
func (k *Key) UnmarshalText(text []byte) error {
	if key, ok := utils.KeyMap[string(text)]; ok {
		k.Key = key
		return nil
	}
	return fmt.Errorf("unknown key: '%s'", string(text))
}

// MouseButton is a mouse button
type MouseButton struct {
	MouseButton ebiten.MouseButton
}

// UnmarshalText fills structure fields from text data
func (b *MouseButton) UnmarshalText(text []byte) error {
	if mouseButton, ok := utils.MouseButtonMap[string(text)]; ok {
		b.MouseButton = mouseButton
		return nil
	}
	return fmt.Errorf("unknown mouse button: '%s'", string(text))
}

// ControllerButton is a gamepad button
type ControllerButton struct {
	ID            int
	GamepadButton ebiten.GamepadButton
}

// UnmarshalTOML fills structure fields from TOML data
func (b *ControllerButton) UnmarshalTOML(i interface{}) error {
	data := i.(map[string]interface{})
	if gamepadButton, ok := utils.GamepadButtonMap[data["button"].(string)]; ok {
		b.ID = int(data["id"].(int64))
		b.GamepadButton = gamepadButton
		return nil
	}
	return fmt.Errorf("unknown gamepad button: '%s'", data["button"].(string))
}

type button struct {
	Key              *Key
	MouseButton      *MouseButton      `toml:"mouse_button"`
	ControllerButton *ControllerButton `toml:"controller"`
}

// Button can be a US keyboard key, a mouse button or a gamepad button
type Button struct {
	Value interface{}
}

// UnmarshalTOML fills structure fields from TOML data
func (b *Button) UnmarshalTOML(i interface{}) error {
	var err error
	b.Value, err = getInterfaceValue(i, &button{})
	return err
}

// Emulated is an emulated axis
type Emulated struct {
	Pos Button
	Neg Button
}

// ControllerAxis is a gamepad axis
type ControllerAxis struct {
	ID       int
	Axis     int
	Invert   bool
	DeadZone float64 `toml:"dead_zone"`
}

// MouseAxis is a mouse axis
type MouseAxis struct {
	Axis int
}

type axis struct {
	Emulated       *Emulated
	ControllerAxis *ControllerAxis `toml:"controller_axis"`
	MouseAxis      *MouseAxis      `toml:"mouse_axis"`
}

// Axis can be an emulated axis, a gamepad axis or a mouse axis
type Axis struct {
	Value interface{}
}

// UnmarshalTOML fills structure fields from TOML data
func (a *Axis) UnmarshalTOML(i interface{}) error {
	var err error
	a.Value, err = getInterfaceValue(i, &axis{})
	return err
}

// Action contains buttons combinations with settings
type Action struct {
	// Combinations contains buttons combinations
	Combinations [][]Button
	// Once determines if the action should be triggered every frame when the button is pressed (default) or only once
	Once bool
}

// Controls contains input controls
type Controls struct {
	// Axes contains axis controls, used for inputs represented by a float value from -1 to 1
	Axes map[string]Axis
	// Actions contains buttons combinations, used for general inputs
	Actions map[string]Action
}

// InputHandler contains input axis values and actions corresponding to specified controls
type InputHandler struct {
	// Axes contains input axis values
	Axes map[string]float64
	// Actions contains input actions
	Actions map[string]bool
}

func getInterfaceValue(treeMap interface{}, data interface{}) (interface{}, error) {
	// Unmarshal from tree
	if tree, err := toml.TreeFromMap(treeMap.(map[string]interface{})); err != nil {
		return nil, err
	} else if err := tree.Unmarshal(data); err != nil {
		return nil, err
	}

	v := reflect.ValueOf(data).Elem()

	// Get non-nil field
	var typeName string
	var value interface{}
	for iField := 0; iField < v.NumField(); iField++ {
		field := v.Field(iField)
		if field.Kind() == reflect.Ptr && !field.IsNil() {
			if typeName != "" {
				return nil, fmt.Errorf("duplicate fields found: %s, %s", field.Elem().Type().Name(), typeName)
			}
			typeName = field.Elem().Type().Name()
			value = field.Interface()
		}
	}
	return value, nil
}
