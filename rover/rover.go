package rover

import (
	"sync"

	rpi "github.com/nathan-osman/go-rpigpio"
)

// singleRover is the Rover instance. Don't create another one
var singleRover *Rover

type Rover struct {
	errorChan     chan OperationError
	directionChan chan Direction

	openMotorPins map[int]*rpi.Pin
	motorsMutex   sync.Mutex
	motorsLocked  bool
}

// Current returns the current rover instance and creates a new one if it doesn't exist yet
func Current() *Rover {
	if singleRover == nil {
		// Rover creation
		singleRover = &Rover{
			errorChan:     make(chan OperationError),
			directionChan: make(chan Direction),

			// Motor Controls
			openMotorPins: make(map[int]*rpi.Pin),
			motorsMutex:   sync.Mutex{},
			motorsLocked:  false,
		}
	}
	return singleRover
}

// Errors returns the channel that will receive all errors
func (r *Rover) Errors() chan OperationError {
	return r.errorChan
}

// Directions returns a channel that sends all directions the rover is taking
func (r *Rover) Directions() chan Direction {
	return r.directionChan
}

// SetDirection sets the current direction
func (r *Rover) SetDirection(d Direction) {
	// Stop before setting new direction
	if r.motorsLocked {
		r.Stop()
	}

	switch d {
	case Left:
		r.turnLeft()
	case Right:
		r.turnRight()
	case Forward:
		r.forward()
	case Reverse:
		r.reverse()
	case PivotLeft:
		r.pivotLeft()
	case PivotRight:
		r.pivotRight()
	case ReverseLeft:
		r.reverseLeft()
	case ReverseRight:
		r.reverseRight()
	default: // case Stop:
		r.Stop()
	}

	r.directionChan <- d
}

// Stop stops the motors
func (r *Rover) Stop() {
	if !r.motorsLocked {
		return
	}

	r.closeMotorPins()
}
