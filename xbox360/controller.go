package xbox360
import (
	"io"
	"log"
	"math"
	"time"

	"github.com/simulatedsimian/joystick"
	"github.com/xarantolus/gover/rover"
)

const (
	serialName = "Microsoft Corp. Xbox 360 Wireless Adapter"
)

var (
	controllerPort io.ReadWriteCloser
)

func init() {
	js, err := joystick.Open(0)
	if err != nil {
		log.Println("Error while setting up controller:", err.Error())
	}

	tick := time.NewTicker(time.Millisecond * 100)

	var rov = rover.Current()

	go func() {
		defer js.Close()
		defer tick.Stop()

		defer rov.Stop() // In case of disconnect

		for {
			state, err := js.Read()
			if err != nil {
				log.Printf("Controller disconnected (error: %s); you need to restart in order to reconnect", err.Error())
				break
			}
			lr, fr := leftStickInterpret(state.AxisData[0]), buttonsInterpret(state.Buttons)

			rov.SetDirection(interpretDirections(lr, fr))
			<-tick.C
		}
	}()
}

func interpretDirections(lr, fr rover.Direction) rover.Direction {
	// Stick is more important than buttons
	// Should reverse in direction if B is pressed (fr == rover.Reverse)

	var finalDir = lr

	if fr == rover.Stop {
		// No input from buttons
		return lr
	}
	if fr == rover.Forward && finalDir == rover.Stop {
		// Overwrite Stop with forward
		finalDir = fr
	}

	if fr == rover.Reverse {
		switch lr {
		case rover.Left:
			finalDir = rover.ReverseLeft // TODO: Might have to switch these
		case rover.Right:
			finalDir = rover.ReverseRight
		case rover.Stop:
			finalDir = rover.Reverse // Accept the normal input
		}
	}

	return finalDir
}

const conTreshold = 22500 // Min = 0; max = 32767

func leftStickInterpret(axisStateZero int) rover.Direction {
	if math.Abs(float64(axisStateZero)) < conTreshold {
		return rover.Stop
	}

	if axisStateZero > 0 {
		return rover.Right
	}

	return rover.Left
}

// A = 2; B = 1; X = 4; Y = 8
func buttonsInterpret(buttonCombinedNumber uint32) rover.Direction {
	if buttonCombinedNumber == 2 {
		return rover.Forward
	}

	if buttonCombinedNumber == 1 {
		return rover.Reverse
	}

	return rover.Stop
}

