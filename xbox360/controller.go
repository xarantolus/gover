package xbox360

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/simulatedsimian/joystick"
	"github.com/xarantolus/gover/rover"
)

func init() {
	var tries int
	go func() {
	retry:
		js, err := joystick.Open(0)
		if err != nil {
			log.Printf("Error while setting up controller (try #%d): %s\n", tries, err.Error())

			if tries == 10 {
				log.Println("Exceeded reconnection limit for controller, the controller will not be available")
				return
			}
			tries++

			time.Sleep(time.Second * 2)
			goto retry
		}

		log.Printf("Successfully connected controller %s\n", js.Name())

		tick := time.NewTicker(time.Millisecond * 100)

		var rov = rover.Current()
		go func() {
			defer js.Close()
			defer tick.Stop()

			defer rov.Stop() // In case of disconnect

			for {
				state, err := js.Read()
				if err != nil {
					panic(fmt.Errorf("Controller disconnected (error: %s); you need to restart in order to reconnect", err.Error()))
				}

				_, _, rest := buttonsCheckXY(state.Buttons)

				lr, fr := leftStickInterpret(state.AxisData[0]), buttonsFRInterpret(rest)

				rov.SetDirection(interpretDirections(lr, fr))
				<-tick.C
			}
		}()
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
		finalDir = fr // TODO: Maybe only go left or right if A is also pressed (fr == Forward)
	}

	if fr == rover.Reverse {
		switch lr {
		case rover.Left:
			finalDir = rover.ReverseLeft
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
// buttonsFRInterpret checks whether to go forward, reverse or stop
func buttonsFRInterpret(buttonCombinedNumber uint32) rover.Direction {
	if buttonCombinedNumber == 2 {
		return rover.Forward
	}

	if buttonCombinedNumber == 1 {
		return rover.Reverse
	}

	return rover.Stop
}

// buttonsCheckXY checks if any of the X or Y buttons are pressed
func buttonsCheckXY(buttonCombinedNumber uint32) (haveX, haveY bool, rest uint32) {
	rest = buttonCombinedNumber

	if rest > 1000 {
		// Is using the left control cross or the RB/LB button
		return
	}

	if rest >= 8 {
		haveY = true
		rest -= 8
	}

	if rest >= 4 {
		haveX = true
		rest -= 4
	}

	return
}
