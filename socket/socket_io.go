package socket

import (
	"fmt"
	"log"
	"strings"
	"time"

	"../controls"
	"github.com/googollee/go-socket.io"
)

var (
	server     *socketio.Server
	haveServer = false
)

const defaultRoom = "default"

// SocketIOServer returns the current Socket.IO server
func SocketIOServer() *socketio.Server {
	if haveServer {
		return server
	}
	var err error
	server, err = socketio.NewServer(nil)
	if err != nil {
		panic(err)
	}

	// Incoming connections
	server.On("connection", func(so socketio.Socket) {
		so.Join(defaultRoom)
		fmt.Printf("Connected client %s\n", so.Id())

		so.On("direction", func(d directionEvent) {
			dir, ok := d.Direction()
			if ok {
				log.Printf("Direction: %v, duration=%s\n", dir, time.Now().Sub(d.Date))

				// TODO: Actually change direction now

				server.BroadcastTo(defaultRoom, "directionchanged", dir)
				// so.BroadcastTo(defaultRoom, "directionchanged", dir)
			}

		})

		so.On("disconnection", func() {
			log.Println("Disconnected")
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

var dirMap = map[string]controls.Direction{
	"LEFT":  controls.LEFT,
	"RIGHT": controls.RIGHT,
	"FRONT": controls.FRONT,
	"BACK":  controls.BACK,
	"STOP":  controls.STOP,
}

func (d *directionEvent) Direction() (dir controls.Direction, ok bool) {
	dir, ok = dirMap[strings.ToUpper(d.RawDirection)]
	return
}
