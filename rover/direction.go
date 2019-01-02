package rover

type Direction string

// Direction for the robot
const (
	Left  Direction = "LEFT"
	Right Direction = "RIGHT"

	Forward Direction = "FORWARD"
	Reverse Direction = "REVERSE"

	PivotLeft  Direction = "PIVOT_LEFT"
	PivotRight Direction = "PIVOT_RIGHT"

	ReverseLeft  Direction = "REVERSE_LEFT"
	ReverseRight Direction = "REVERSE_RIGHT"

	Stop Direction = "STOP"
)
