package inputsystem

import (
	"math"

	"github.com/x-hgg-x/goecsengine/resources"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// InputSystem updates input axis values and actions
func InputSystem(world w.World) {
	for k, v := range world.Resources.Controls.Axes {
		world.Resources.InputHandler.Axes[k] = getAxisValue(world, v)
	}

	for k, v := range world.Resources.Controls.Actions {
		world.Resources.InputHandler.Actions[k] = isActionDone(v)
	}
}

func getAxisValue(world w.World, axis resources.Axis) float64 {
	axisValue := 0.0

	switch value := axis.Value.(type) {
	case *resources.Emulated:
		if isPressed(value.Pos) {
			axisValue++
		}
		if isPressed(value.Neg) {
			axisValue--
		}
	case *resources.ControllerAxis:
		deadZone := math.Abs(value.DeadZone)
		axisValue = ebiten.GamepadAxis(value.ID, value.Axis)

		if axisValue < -deadZone {
			axisValue = (axisValue + deadZone) / (1.0 - deadZone)
		} else if axisValue > deadZone {
			axisValue = (axisValue - deadZone) / (1.0 - deadZone)
		} else {
			axisValue = 0
		}

		if value.Invert {
			axisValue *= -1
		}
	case *resources.MouseAxis:
		screenWidth := float64(world.Resources.ScreenDimensions.Width)
		screenHeight := float64(world.Resources.ScreenDimensions.Height)

		x, y := ebiten.CursorPosition()
		switch value.Axis {
		case 0:
			axisValue = float64(x) / screenWidth
		case 1:
			axisValue = (screenHeight - float64(y)) / screenHeight
		}
		axisValue = 2*axisValue - 1
	}

	// Axis value must be between -1 and 1
	axisValue = math.Min(math.Max(axisValue, -1), 1)
	return axisValue
}

func isActionDone(action resources.Action) bool {
	var funcPressed func(resources.Button) bool
	if action.Once {
		funcPressed = isJustPressed
	} else {
		funcPressed = isPressed
	}

	for _, combination := range action.Combinations {
		actionDone := true
		for _, button := range combination {
			actionDone = actionDone && funcPressed(button)
		}
		if actionDone {
			return true
		}
	}
	return false
}

func isPressed(b resources.Button) bool {
	switch value := b.Value.(type) {
	case *resources.Key:
		return ebiten.IsKeyPressed(value.Key)
	case *resources.MouseButton:
		return ebiten.IsMouseButtonPressed(value.MouseButton)
	case *resources.ControllerButton:
		return ebiten.IsGamepadButtonPressed(value.ID, value.GamepadButton)
	}
	return false
}

func isJustPressed(b resources.Button) bool {
	switch value := b.Value.(type) {
	case *resources.Key:
		return inpututil.IsKeyJustPressed(value.Key)
	case *resources.MouseButton:
		return inpututil.IsMouseButtonJustPressed(value.MouseButton)
	case *resources.ControllerButton:
		return inpututil.IsGamepadButtonJustPressed(value.ID, value.GamepadButton)
	}
	return false
}
