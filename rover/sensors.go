package rover

import (
	"time"

	"github.com/nathan-osman/go-rpigpio"
)

const (
	frontTrigger = 23 // 16
	frontEcho    = 24 // 18
)

// DistanceFront gets the distance from the front sensor in cm
// If the signal takes more than one second we abort and return not ok
func (r *Rover) DistanceFront() (cm float32, ok bool) {
	r.sensorFrontMutex.Lock()
	defer r.sensorFrontMutex.Unlock()
	// Open trigger pin
	pinTrigger, err := rpi.OpenPin(frontTrigger, rpi.OUT)
	r.checkPanic(err, "while opening front sensor trigger pin %d", frontTrigger)

	// Open echo receiver pin
	pinEcho, err := rpi.OpenPin(frontEcho, rpi.IN)
	r.checkPanic(err, "while opening front sensor echo pin %d", frontEcho)

	// Send impulse to trigger
	r.checkPanic(pinTrigger.Write(rpi.HIGH), "writing HIGH to trigger pin %d", frontTrigger)
	time.Sleep(30 * time.Microsecond)
	r.checkPanic(pinTrigger.Write(rpi.LOW), "writing LOW to trigger pin %d", frontTrigger)
	time.Sleep(30 * time.Microsecond)

	// Abort after 1 second
	abortNow := false
	go func() {
		time.Sleep(1 * time.Second)
		abortNow = true
	}()

	for {
		val, err := pinEcho.Read()
		r.checkErr(err, "reading from front sensor echo pin %d (while waiting for HIGH)", frontEcho)

		if val == rpi.HIGH {
			break
		} else if abortNow {
			return 0, false
		}
	}
	startTime := time.Now()
	for {
		val, err := pinEcho.Read()
		r.checkErr(err, "reading from front sensor echo pin %d (while waiting for LOW)", frontEcho)

		if val == rpi.LOW {
			break
		} else if abortNow {
			return 0, false
		}
	}
	endTime := time.Now()

	r.checkErr(pinTrigger.Close(), "closing front sensor trigger pin %d", frontTrigger)
	r.checkErr(pinEcho.Close(), "closing front sensor echo pin %d", frontEcho)

	dist := distanceFromSoundSpeed(endTime.Sub(startTime))

	// According to https://cdn.sparkfun.com/datasheets/Sensors/Proximity/HCSR04.pdf the range is from 2cm-400cm -> we say it's not ok
	if dist < 2 || dist > 400 {
		return 0, false
	}

	return dist, true
}

func distanceFromSoundSpeed(dur time.Duration) float32 {
	return float32(dur.Seconds() * 34300 / 2)
}
