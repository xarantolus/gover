package main

import (
	"fmt"
	"net/http"

	"./socket"
)

const port = "903"

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.Handle("/socket.io/", socket.SocketIOServer())

	fmt.Printf("Gover server listening on port %s\n", port)
	panic(http.ListenAndServe(":903", nil))
}
