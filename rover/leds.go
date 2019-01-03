package rover

import (
	"github.com/nathan-osman/go-rpigpio"
)

const (
	// One cable needs to be connected to a Ground, another one to the Pin defined below

	// frontLed is the status led that is on when the server is running. It goes off when the robot shuts down
	frontLed = 17 // Pin 11, Ground is pin 06
)

// ToggleFrontLED toggles the front led of the robot
func (r *Rover) ToggleFrontLED() {
	r.frontLEDMutex.Lock()
	if r.frontLEDOn {
		r.checkErr(r.setLEDPin(frontLed, rpi.LOW), "while setting front led pin %d to LOW", frontLed)
	} else {
		r.checkErr(r.setLEDPin(frontLed, rpi.HIGH), "while setting front led pin %d to HIGH", frontLed)
	}

	// Reflect new state in variable
	r.frontLEDOn = !r.frontLEDOn
	r.frontLEDMutex.Unlock()
}

// IsFrontLEDOn returns whether the front led is turned on or off
func (r *Rover) IsFrontLEDOn() bool {
	return r.frontLEDOn
}

func (r *Rover) closeAllLEDs() {
	r.openLEDPinsMutex.Lock()
	defer r.openLEDPinsMutex.Unlock()

	for number, pin := range r.openLEDPins {
		if pin != nil {
			r.checkErr(pin.Close(), "closing led pin %d", number)
			delete(r.openLEDPins, number)
		}
	}
}

// setLEDPin sets the pin with the number `pinNumber` to the value `newValue`
func (r *Rover) setLEDPin(pinNumber int, newValue rpi.Value) (err error) {
	// See if the given pin is already in our map
	r.openLEDPinsMutex.Lock()
	pin, ok := r.openLEDPins[pinNumber]
	r.openLEDPinsMutex.Unlock()

	// If not, we load it
	if !ok {
		pin, err = rpi.OpenPin(pinNumber, rpi.OUT)
		if err != nil {
			return
		}
	}

	// Write the desired value
	err = pin.Write(newValue)
	if err != nil {
		return
	}

	// Make sure map access is safe when we delete
	r.openLEDPinsMutex.Lock()
	defer r.openLEDPinsMutex.Unlock()

	// Add to or remove from open pins
	if newValue == rpi.HIGH {
		r.openLEDPins[pinNumber] = pin
	} else {
		// Turned the led off, so we close it
		err = pin.Close()
		if err != nil {
			return
		}

		delete(r.openLEDPins, pinNumber)
	}

	return
}
