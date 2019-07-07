package socket

import (
	"log"
	"strings"
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/xarantolus/gover/rover"
)

var (
	server     *socketio.Server
	haveServer = false
)

const defaultRoom = "default"

// CreateServer returns the current Socket.IO server
func CreateServer() *socketio.Server {
	if haveServer {
		return server
	}
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		panic(err)
	}

	var rov = rover.Current()

	// Send all directions to clients
	go func() {
		for dir := range rov.Directions() {
			server.BroadcastTo(defaultRoom, "directionchanged", dir)
		}
	}()

	// Read sensor and send distance to clients
	go func() {
		for /* ever */ {
			if server.Count() > 0 {
				now := time.Now()
				dist, ok := rov.DistanceFront()
				server.BroadcastTo(defaultRoom, "sensor:front", sensorEvent{
					Distance: dist,
					Ok:       ok,
				})
				// Try to send this signal once per 1/4 second
				time.Sleep(250*time.Millisecond - time.Since(now))
			} else {
				time.Sleep(time.Second)
			}
		}
	}()

	// Incoming connections
	server.On("connection", func(so socketio.Socket) {
		so.Join(defaultRoom)
		log.Printf("Connected client %s\n", so.Id())

		// Send the initial direction
		so.Emit("directionchanged", rov.CurrentDirection())

		so.On("direction", func(d directionEvent) {
			dir, ok := d.Direction()
			if !ok {
				log.Printf("Couldn't parse direction %s\n", d.RawDirection)
				return
			}

			durationSinceSend := time.Since(d.Date)
			log.Printf("Direction: %v, duration=%s\n", dir, durationSinceSend)

			// Only drop packages that aren't stop
			if dir != rover.Stop && durationSinceSend > 3*time.Second {
				// Drop package
				log.Println("Dropped package")
				return
			}

			// Set the direction for the rover
			rov.SetDirection(dir)
		})

		so.On("disconnection", func() {
			log.Println("Client disconnected, stopping rover")
			rov.Stop()
		})
	})
	// Log errors
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("Socket.io error:", err)
	})

	haveServer = true

	return server
}

type socketResponse struct {
	Success bool               `json:"success"`
	Error   *socketResponseErr `json:"error,omitempty"`
}
type socketResponseErr struct {
	Message string `json:"message"`
}

type directionEvent struct {
	Date         time.Time `json:"date"`
	RawDirection string    `json:"direction"`
}

type sensorEvent struct {
	Ok       bool    `json:"ok"`
	Distance float32 `json:"distance"`
}

var (
	dirMap = map[string]rover.Direction{
		"LEFT":    rover.Left,
		"RIGHT":   rover.Right,
		"FORWARD": rover.Forward,
		"REVERSE": rover.Reverse,

		"PIVOT_LEFT":  rover.PivotLeft,
		"PIVOT_RIGHT": rover.PivotRight,

		"REVERSE_LEFT":  rover.ReverseLeft,
		"REVERSE_RIGHT": rover.ReverseRight,

		"STOP": rover.Stop,
	}
	dirMapMut = sync.Mutex{}
)

func (d *directionEvent) Direction() (dir rover.Direction, ok bool) {
	dirMapMut.Lock()
	dir, ok = dirMap[strings.ToUpper(d.RawDirection)]
	dirMapMut.Unlock()
	return
}
