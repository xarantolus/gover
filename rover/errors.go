package rover

import (
	"fmt"
	"time"
)

// OperationError is an error that contains information about errors that happened to the rover
type OperationError struct {
	// Time describes the time when the error happened
	Time time.Time
	// The actual error that happened
	ActualError error
	// Operation describes what the robot was doing when the error happened
	Operation string
}

var (
	errMotorsInUse = fmt.Errorf("Motors are already in use")
)

func (o OperationError) Error() string {
	return fmt.Sprintf("%s: %s: %s", o.Time.Format("2006-01-02 15:04:05"), o.Operation, o.ActualError.Error())
}

// sendErr sends an error to the error channel
func (r *Rover) sendErr(err error, operation string, a ...interface{}) {
	r.errorChan <- OperationError{
		ActualError: err,
		Time:        time.Now(),
		Operation:   fmt.Sprintf(operation, a...),
	}
}

func (r *Rover) checkErr(err error, operation string, a ...interface{}) {
	if err != nil {
		r.sendErr(err, operation, a...)
	}
}

func (r *Rover) checkPanic(err error, operation string, a ...interface{}) {
	if err != nil {
		panic(fmt.Errorf("Error while %s: %s", fmt.Sprintf(operation, a...), err.Error()))
	}
}
