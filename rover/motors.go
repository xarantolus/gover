package rover

import (
	"time"

	"github.com/nathan-osman/go-rpigpio"
)

// Input pins for the L298N Dual H Bridge Motor Controller board
const (
	// GPIO Numbers, in comments pin numbers
	in1 = 5  //29
	in2 = 6  // 31
	in3 = 16 // 36
	in4 = 26 // 37
)

var motorPins = []int{in1, in2, in3, in4}

// getMotorPins returns all motor pins as a map[pinNumber]Pin
func (r *Rover) getMotorPins() (map[int]*rpi.Pin, error) {
	if r.motorsLocked {
		return nil, errMotorsInUse
	}

	r.motorsMutex.Lock()
	r.motorsLocked = true

	var err error

	for _, pinNumber := range motorPins {
		r.openMotorPins[pinNumber], err = rpi.OpenPin(pinNumber, rpi.OUT)
		r.checkPanic(err, "opening pin %d", pinNumber)
	}

	// Turn status led on as we're going to start the motors
	if !r.IsFrontLEDOn() {
		r.ToggleFrontLED()
	}

	return r.openMotorPins, nil
}

// closeMotorPins closes all motors
func (r *Rover) closeMotorPins() {
	for number, pin := range r.openMotorPins {
		if pin != nil {
			r.checkErr(pin.Close(), "closing motor pin %d", number)
			delete(r.openMotorPins, number)
		}
	}
	r.motorsLocked = false
	r.motorsMutex.Unlock()

	// Turn status led off
	if r.IsFrontLEDOn() {
		r.ToggleFrontLED()
	}
}

func (r *Rover) outputMotors(val1, val2, val3, val4 bool) error {
	pins, err := r.getMotorPins()
	if err != nil {
		return err
	}

	err = pins[in1].Write(getValueFromBool(val1))
	if err != nil {
		return err
	}

	err = pins[in2].Write(getValueFromBool(val2))
	if err != nil {
		return err
	}

	err = pins[in3].Write(getValueFromBool(val3))
	if err != nil {
		return err
	}

	err = pins[in4].Write(getValueFromBool(val4))
	if err != nil {
		return err
	}

	return nil
}

// canGoForward checks if the robot can go forward and stops it if it can't
func (r *Rover) canGoForward() (canGo bool) {
	// Make sure the rover doesn't run straight into a wall if we have a recent distance
	if r.sensorFrontLastDistance < roverMinFrontDistanceCM && time.Now().Sub(r.sensorFrontLastDate) < 2*time.Second {
		r.SetDirection(Stop)
		return false
	}
	return true
}

func (r *Rover) forward() (ok bool) {
	if r.canGoForward() {
		r.checkErr(r.outputMotors(false, true, true, false), "going forward")
		return true
	}
	return false
}

func (r *Rover) reverse() {
	r.checkErr(r.outputMotors(true, false, false, true), "reversing")
}

func (r *Rover) turnLeft() {
	r.checkErr(r.outputMotors(true, true, true, false), "turning left")
}

func (r *Rover) turnRight() {
	r.checkErr(r.outputMotors(false, true, false, false), "turning right")
}

func (r *Rover) reverseLeft() {
	r.checkErr(r.outputMotors(false, false, false, true), "reversing left")
}

func (r *Rover) reverseRight() {
	r.checkErr(r.outputMotors(true, false, false, false), "reversing right")
}

func (r *Rover) pivotLeft() {
	r.checkErr(r.outputMotors(true, false, true, false), "pivoting left")
}

func (r *Rover) pivotRight() {
	r.checkErr(r.outputMotors(false, true, false, true), "pivoting right")
}
