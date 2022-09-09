package animationsystem

import (
	"math"

	c "github.com/x-hgg-x/goecsengine/components"
	m "github.com/x-hgg-x/goecsengine/math"
	"github.com/x-hgg-x/goecsengine/utils"
	w "github.com/x-hgg-x/goecsengine/world"

	"github.com/hajimehoshi/ebiten/v2"
	ecs "github.com/x-hgg-x/goecs/v2"
)

// AnimationSystem updates animations
func AnimationSystem(world w.World) {
	world.Manager.Join(world.Components.Engine.SpriteRender, world.Components.Engine.AnimationControl).Visit(ecs.Visit(func(entity ecs.Entity) {
		sprite := world.Components.Engine.SpriteRender.Get(entity).(*c.SpriteRender)
		animationControl := world.Components.Engine.AnimationControl.Get(entity).(*c.AnimationControl)

		times := animationControl.Animation.Time
		spriteNumbers := animationControl.Animation.SpriteNumber

		currentTime := animationControl.GetState().CurrentTime

		// Process command
		switch animationControl.Command.Type {
		case c.AnimationCommandRestart:
			animationControl.SetStateType(c.ControlStateRunning)
			currentTime = 0

		case c.AnimationCommandStart:
			animationControl.SetStateType(c.ControlStateRunning)

		case c.AnimationCommandStepBackward:
			animationPos := computeAnimationPos(currentTime, times)
			if animationControl.End.Type == c.EndControlLoop && animationPos == 0 {
				animationPos = len(spriteNumbers) - 1
			} else {
				animationPos = m.Max(0, animationPos-1)
			}
			currentTime = times[animationPos]

		case c.AnimationCommandStepForward:
			animationPos := computeAnimationPos(currentTime, times)
			if animationControl.End.Type == c.EndControlLoop && animationPos == len(spriteNumbers)-1 {
				animationPos = 0
			} else {
				animationPos = m.Min(len(spriteNumbers)-1, animationPos+1)
			}
			currentTime = times[animationPos]

		case c.AnimationCommandSetTime:
			currentTime = math.Min(math.Max(animationControl.Command.Time, times[0]), times[len(times)-1])

		case c.AnimationCommandPause:
			animationControl.SetStateType(c.ControlStatePaused)

		case c.AnimationCommandAbort:
			entity.RemoveComponent(world.Components.Engine.AnimationControl)
			return

		case c.AnimationCommandNone:
			break

		default:
			utils.LogFatalf("unknown animation command: %v", animationControl.Command.Type)
		}

		// Reset command
		animationControl.Command = c.AnimationCommand{}

		// Run animation
		if animationControl.GetState().Type == c.ControlStateRunning {
			currentTime += animationControl.RateMultiplier / float64(ebiten.DefaultTPS)
		}

		// Check animation end
		if animationControl.End.Type == c.EndControlLoop {
			currentTime = mod(currentTime, times[len(times)-1])
		} else if currentTime >= times[len(times)-1] {
			switch animationControl.End.Type {
			case c.EndControlNormal:
				animationControl.SetStateType(c.ControlStateDone)
				currentTime = 0
			case c.EndControlStay:
				animationControl.SetStateType(c.ControlStateDone)
				currentTime = times[len(times)-1]
			default:
				utils.LogFatalf("unknown end control: %v", animationControl.End.Type)
			}
		}

		// Set animation state
		animationControl.SetCurrentTime(currentTime)
		sprite.SpriteNumber = spriteNumbers[computeAnimationPos(currentTime, times)]
	}))
}

func computeAnimationPos(currentTime float64, times []float64) int {
	animationPos := 0
	for iTime := range times[1:] {
		animationPos = iTime
		if times[iTime+1] > currentTime {
			break
		}
	}
	return animationPos
}

func mod(a, b float64) float64 {
	m := math.Mod(a, b)
	if m < 0 {
		m += math.Abs(b)
	}
	return m
}
