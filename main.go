package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"./rover"
	"./socket"
)

const port = "80"

func main() {
	var r = rover.Current()

	defer func() {
		if rec := recover(); rec != nil {
			r.Stop()

			fmt.Println("Recovered panic, exiting:", rec)
		}
	}()

	go func() {
		for err := range r.Errors() {
			fmt.Println("Error: " + err.Error())
		}
	}()

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.Handle("/socket.io/", socket.CreateServer())

	fmt.Printf("Gover server listening on port %s\n", port)
	go panic(http.ListenAndServe(":"+port, nil))

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-signals

	fmt.Println("Stopping")
	r.Stop()
}
