package main

import (
	"fmt"

	c "github.com/x-hgg-x/goecsengine/components"
	w "github.com/x-hgg-x/goecsengine/world"

	ecs "github.com/x-hgg-x/goecs/v2"
)

// DemoSystem is a demo system
func DemoSystem(world w.World) {
	lastAction := ""
	rateMultiplier := 0.0

	// Check if there are running animations
	runningAnimations := world.Manager.Join(world.Components.Engine.AnimationControl).Visit(
		func(index int) (skip bool) {
			return world.Components.Engine.AnimationControl.Get(ecs.Entity(index)).(*c.AnimationControl).GetState().Type == c.ControlStateRunning
		})

	// Restart
	if world.Resources.InputHandler.Actions[RestartAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl).Command.Type = c.AnimationCommandRestart
		}))
		lastAction = "Restart"
	}

	// Start and pause
	if world.Resources.InputHandler.Actions[StartPauseAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			if runningAnimations {
				world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl).Command.Type = c.AnimationCommandPause
				lastAction = "Pause"
			} else {
				world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl).Command.Type = c.AnimationCommandStart
				lastAction = "Start"
			}
		}))
	}

	// Step backward
	if world.Resources.InputHandler.Actions[StepBackwardAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl).Command.Type = c.AnimationCommandStepBackward
		}))
		lastAction = "StepBackward"
	}

	// Step forward
	if world.Resources.InputHandler.Actions[StepForwardAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl).Command.Type = c.AnimationCommandStepForward
		}))
		lastAction = "StepForward"
	}

	// Reverse animation
	if world.Resources.InputHandler.Actions[ReverseAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			animationControl := world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl)
			animationControl.RateMultiplier *= -1
			rateMultiplier = animationControl.RateMultiplier
		}))
		lastAction = "Reverse"
	}

	// Divide animation speed by 2
	if world.Resources.InputHandler.Actions[HalfSpeedAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			animationControl := world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl)
			animationControl.RateMultiplier *= 0.5
			rateMultiplier = animationControl.RateMultiplier
		}))
		lastAction = "HalfSpeed"
	}

	// Multiply animation speed by 2
	if world.Resources.InputHandler.Actions[DoubleSpeedAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			animationControl := world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl)
			animationControl.RateMultiplier *= 2
			rateMultiplier = animationControl.RateMultiplier
		}))
		lastAction = "DoubleSpeed"
	}

	// Set animation time to the middle of animation
	if world.Resources.InputHandler.Actions[SetTimeToMiddleAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			animationControl := world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl)
			times := animationControl.Animation.Time
			middleTime := (times[0] + times[len(times)-1]) / 2
			animationControl.Command.Type = c.AnimationCommandSetTime
			animationControl.Command.Time = middleTime
		}))
		lastAction = "SetTimeToMiddle"
	}

	// Abort and remove animation
	if world.Resources.InputHandler.Actions[AbortAction] {
		world.Manager.Join(world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
			world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl).Command.Type = c.AnimationCommandAbort
		}))
		lastAction = "Abort"
	}

	// Get animation state
	animationStates := []c.ControlState{}
	spriteNumbers := []int{}
	world.Manager.Join(world.Components.Engine.SpriteRender, world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
		animationStates = append(animationStates, world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl).GetState())
		spriteNumbers = append(spriteNumbers, world.Components.Engine.SpriteRender.Get(entity).(*c.SpriteRender).SpriteNumber)
	}))

	// Update text info
	currentBat := 0
	aborted := world.Manager.Join(world.Components.Engine.AnimationControl).Empty()
	world.Manager.Join(world.Components.Engine.Text, world.Components.Engine.UITransform).Visit(ecs.Visit(func(entity ecs.Entity) {
		text := world.Components.Engine.Text.Get(entity).(*c.Text)
		if text.ID == "last_command" && lastAction != "" {
			text.Text = fmt.Sprintf("Last command: %s", lastAction)
		}
		if text.ID == "aborted" && aborted {
			text.Text = "Aborted: true (Press Enter to reset)"
		}
		if text.ID == "rate_multiplier" && rateMultiplier != 0 {
			text.Text = fmt.Sprintf("Rate multiplier: %v", rateMultiplier)
		}
		if text.ID == fmt.Sprintf("bat%v", currentBat) && !aborted {
			state := ""
			switch animationStates[currentBat].Type {
			case c.ControlStateNotStarted:
				state = "NotStarted"
			case c.ControlStateRunning:
				state = "Running"
			case c.ControlStatePaused:
				state = "Paused"
			case c.ControlStateDone:
				state = "Done"
			}
			text.Text = fmt.Sprintf("Time: %.2f / Sprite: %v / %s", animationStates[currentBat].CurrentTime, spriteNumbers[currentBat], state)
			currentBat++
		}
	}))

	// Reset entities
	if world.Resources.InputHandler.Actions[ResetAction] {
		world.Manager.DeleteAllEntities()
		LoadEntities("metadata/game.toml", world)
		LoadEntities("metadata/text.toml", world)
	}
}
