package components

import "fmt"

// Animation structure
type Animation struct {
	// List of times (must be in strictly increasing order, with first element equal to 0)
	Time []float64
	// List of sprite numbers (must have one less element than the Time field, and at least one element)
	SpriteNumber []int `toml:"sprite_number"`
}

// UnmarshalTOML fills structure fields from TOML data
func (a *Animation) UnmarshalTOML(i interface{}) error {
	data := i.(map[string]interface{})
	times := data["time"].([]interface{})
	spriteNumbers := data["sprite_number"].([]interface{})

	a.SpriteNumber = make([]int, len(spriteNumbers))
	for iSprite := range spriteNumbers {
		a.SpriteNumber[iSprite] = int(spriteNumbers[iSprite].(int64))
	}

	// Check animation length
	if len(spriteNumbers) < 1 || len(times) != len(spriteNumbers)+1 {
		return fmt.Errorf("incorrect animation length: len(Time) = %v and len(SpriteNumber) = %v", len(times), len(spriteNumbers))
	}

	// Check time values
	a.Time = make([]float64, len(times))
	if times[0].(float64) != 0 {
		return fmt.Errorf("first time value must be 0")
	}
	for iTime := 1; iTime < len(times); iTime++ {
		if times[iTime].(float64) <= times[iTime-1].(float64) {
			return fmt.Errorf("time values must be in strictly increasing order")
		}
		a.Time[iTime] = times[iTime].(float64)
	}

	return nil
}

// EndControlType is an end control type
type EndControlType int

const (
	// EndControlNormal goes back to the start of the animation
	EndControlNormal EndControlType = iota
	// EndControlStay stays at the end of the animation
	EndControlStay
	// EndControlLoop loops the animation
	EndControlLoop
)

// EndControl structure
type EndControl struct {
	// End control type
	Type EndControlType
}

// AnimationCommandType is an animation command type
type AnimationCommandType int

const (
	// AnimationCommandNone does nothing
	AnimationCommandNone AnimationCommandType = iota
	// AnimationCommandRestart restarts the animation
	AnimationCommandRestart
	// AnimationCommandStart starts the animation
	AnimationCommandStart
	// AnimationCommandStepBackward steps backward
	AnimationCommandStepBackward
	// AnimationCommandStepForward steps forward
	AnimationCommandStepForward
	// AnimationCommandSetTime sets animation time to the specified time value
	AnimationCommandSetTime
	// AnimationCommandPause pauses the animation
	AnimationCommandPause
	// AnimationCommandAbort aborts and removes the animation from entity
	AnimationCommandAbort
)

// AnimationCommand structure
type AnimationCommand struct {
	// Animation command type
	Type AnimationCommandType
	// Command time, used only with AnimationCommandSetTime
	Time float64
}

// ControlStateType is a control state type
type ControlStateType int

const (
	// ControlStateNotStarted is the default state
	ControlStateNotStarted ControlStateType = iota
	// ControlStateRunning is the running state
	ControlStateRunning
	// ControlStatePaused is the paused state
	ControlStatePaused
	// ControlStateDone is the done state
	ControlStateDone
)

// ControlState structure
type ControlState struct {
	// Control state type
	Type ControlStateType
	// Current animation time
	CurrentTime float64
}

// AnimationControl component
type AnimationControl struct {
	// Reference animation
	Animation *Animation
	// End control
	End EndControl
	// Animation command
	Command AnimationCommand
	// Animation speed multiplier
	RateMultiplier float64
	// Animation state
	state ControlState
}

// GetState returns animation control state
func (c *AnimationControl) GetState() ControlState {
	return c.state
}

// SetStateType sets animation control state type
func (c *AnimationControl) SetStateType(stateType ControlStateType) {
	c.state.Type = stateType
}

// SetCurrentTime sets current animation time
func (c *AnimationControl) SetCurrentTime(time float64) {
	c.state.CurrentTime = time
}
