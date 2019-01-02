package socket

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"../rover"
	"github.com/googollee/go-socket.io"
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

	go func() {
		for dir := range rov.Directions() {
			server.BroadcastTo(defaultRoom, "directionchanged", dir)
		}
	}()

	// Incoming connections
	server.On("connection", func(so socketio.Socket) {
		so.Join(defaultRoom)
		fmt.Printf("Connected client %s\n", so.Id())

		so.On("direction", func(d directionEvent) {
			dir, ok := d.Direction()
			if !ok {
				fmt.Printf("Couldn't parse direction %s\n", d.RawDirection)
				return
			}

			durationSinceSend := time.Now().Sub(d.Date)
			log.Printf("Direction: %v, duration=%s\n", dir, durationSinceSend)

			if durationSinceSend > 3*time.Second {
				// Drop package
				log.Println("Dropped package")
				return
			}

			// Set the direction for the rover
			rov.SetDirection(dir)
		})

		so.On("disconnection", func() {
			log.Println("Disconnected")

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
