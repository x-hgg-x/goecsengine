package states

import (
	"fmt"
	"os"

	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten"
)

// TransType is a transition type
type TransType int

const (
	// TransNone is the null transition
	TransNone TransType = iota
	// TransPop removes the active state and resume the next state
	TransPop
	// TransPush pauses the active state and add new states to the stack
	TransPush
	// TransSwitch removes the active state and replace it by a new one
	TransSwitch
	// TransReplace removes all states and insert a new stack
	TransReplace
	// TransQuit removes all states and quit
	TransQuit
)

// Transition is a state transition
type Transition struct {
	TransType TransType
	NewStates []State
}

// State is a game state
type State interface {
	// Executed when the state begins
	OnStart(world w.World)
	// Executed when the state exits
	OnStop(world w.World)
	// Executed when a new state is pushed over this one
	OnPause(world w.World)
	// Executed when the state become active again (states pushed over this one have been popped)
	OnResume(world w.World)
	// Executed on every frame when the state is active
	Update(world w.World, screen *ebiten.Image) Transition
}

// StateMachine contains a stack of states.
// Only the top state is active.
type StateMachine struct {
	states []State
}

// Init creates a new state machine with an initial state
func Init(s State, world w.World) StateMachine {
	s.OnStart(world)
	return StateMachine{[]State{s}}
}

// Update updates the state machine
func (sm *StateMachine) Update(world w.World, screen *ebiten.Image) {
	if len(sm.states) < 1 {
		os.Exit(0)
	}

	switch t := sm.states[len(sm.states)-1].Update(world, screen); t.TransType {
	case TransPop:
		sm._Pop(world)
	case TransPush:
		sm._Push(world, t.NewStates)
	case TransSwitch:
		sm._Switch(world, t.NewStates)
	case TransReplace:
		sm._Replace(world, t.NewStates)
	case TransQuit:
		sm._Quit(world)
	}
}

// Remove the active state and resume the next state
func (sm *StateMachine) _Pop(world w.World) {
	sm.states[len(sm.states)-1].OnStop(world)
	sm.states = sm.states[:len(sm.states)-1]

	if len(sm.states) > 0 {
		sm.states[len(sm.states)-1].OnResume(world)
	}
}

// Pause the active state and add new states to the stack
func (sm *StateMachine) _Push(world w.World, newStates []State) {
	if len(newStates) > 0 {
		sm.states[len(sm.states)-1].OnPause(world)

		for _, state := range newStates[:len(newStates)-1] {
			state.OnStart(world)
			state.OnPause(world)
		}
		newStates[len(newStates)-1].OnStart(world)

		sm.states = append(sm.states, newStates...)
	}
}

// Remove the active state and replace it by a new one
func (sm *StateMachine) _Switch(world w.World, newStates []State) {
	if len(newStates) != 1 {
		utils.LogError(fmt.Errorf("switch transition accept only one new state"))
	}

	sm.states[len(sm.states)-1].OnStop(world)
	newStates[0].OnStart(world)
	sm.states[len(sm.states)-1] = newStates[0]
}

// Remove all states and insert a new stack
func (sm *StateMachine) _Replace(world w.World, newStates []State) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}

	if len(newStates) > 0 {
		for _, state := range newStates[:len(newStates)-1] {
			state.OnStart(world)
			state.OnPause(world)
		}
		newStates[len(newStates)-1].OnStart(world)
	}
	sm.states = newStates
}

// Remove all states and quit
func (sm *StateMachine) _Quit(world w.World) {
	for len(sm.states) > 0 {
		sm.states[len(sm.states)-1].OnStop(world)
		sm.states = sm.states[:len(sm.states)-1]
	}
	os.Exit(0)
}
