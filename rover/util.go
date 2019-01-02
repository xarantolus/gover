package rover

import (
	"github.com/nathan-osman/go-rpigpio"
)

func getValueFromBool(b bool) rpi.Value {
	if b {
		return rpi.HIGH
	}
	return rpi.LOW
}
