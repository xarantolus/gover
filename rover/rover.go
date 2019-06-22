package rover

import (
	"fmt"
	"sync"
	"time"

	rpi "github.com/nathan-osman/go-rpigpio"
)

// singleRover is the Rover instance. Don't create another one
var (
	singleRover      *Rover
	ErrRoverShutdown = fmt.Errorf("The rover is shutting down")

	// don't allow going forward when the front distance is lower than this
	// 25 is perfect for the speed. If the rover catches the wall about 20cm before hitting it, it has enough time to "brake"
	roverMinFrontDistanceCM float32 = 25
)

type Rover struct {
	// Channels
	errorChan     chan OperationError
	directionChan chan Direction

	currentDirection Direction

	// LEDs
	frontLEDOn    bool
	frontLEDMutex sync.Mutex

	// Keep track of all leds (there's only one right now, so this is not very necessary)
	openLEDPinsMutex sync.Mutex
	openLEDPins      map[int]*rpi.Pin

	// Sensors
	sensorFrontMutex        sync.Mutex
	sensorFrontLastDistance float32
	sensorFrontLastDate     time.Time

	// Motor Controls
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

			// Initial direction
			currentDirection: Stop,

			// LEDs
			frontLEDOn:    false,
			frontLEDMutex: sync.Mutex{},

			openLEDPinsMutex: sync.Mutex{},
			openLEDPins:      make(map[int]*rpi.Pin),

			// Sensors
			sensorFrontMutex:        sync.Mutex{},
			sensorFrontLastDistance: 100, // Allow going forward from the beginning
			sensorFrontLastDate:     time.Now(),

			// Motor Controls
			openMotorPins: make(map[int]*rpi.Pin),
			motorsMutex:   sync.Mutex{},
			motorsLocked:  false,
		}
	}
	return singleRover
}

// CurrentDirection gives the direction the robot is currently going
func (r *Rover) CurrentDirection() Direction {
	return r.currentDirection
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
	// Don't stop if we already have that direction
	if r.currentDirection == d {
		return
	}

	// Stop before setting new direction
	if r.motorsLocked {
		r.Stop()
	}

	var success = true
	switch d {
	case Left:
		r.turnLeft()
	case Right:
		r.turnRight()
	case Forward:
		// Forward might cancel itself
		success = r.forward()
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

	if success {
		// change & notify of direction state
		r.currentDirection = d
		r.directionChan <- d
	}
}

// Stop stops the motors
func (r *Rover) Stop() {
	if !r.motorsLocked {
		return
	}

	r.closeMotorPins()

	// Notify
	r.currentDirection = Stop
	r.directionChan <- Stop
}

// Shutdown stops and shuts down the rover and closes all open GPIO pins. The rover shouldn't be used after this
func (r *Rover) Shutdown() {
	// Shutdown motors,
	r.Stop()
	// .. and leds
	r.closeAllLEDs()

	close(r.errorChan)
	close(r.directionChan)

	// Make sure that the last distance scan has returned
	r.sensorFrontMutex.Lock()
	r.sensorFrontMutex.Unlock()
}
