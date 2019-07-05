package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/xarantolus/gover/rover"
	"github.com/xarantolus/gover/socket"
	_ "github.com/xarantolus/gover/xbox360"
)

const port = "80"

func main() {
	var r = rover.Current()

	// Print all erros the rover encounters
	go func() {
		for err := range r.Errors() {
			log.Println("Error: " + err.Error())
		}
	}()

	// Create web server
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/", fs)
	mux.Handle("/socket.io/", socket.CreateServer())

	server := http.Server{
		Addr: ":" + port, Handler: mux,
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Wait for shutdown signals
	go func() {
		<-sc

		log.Println("Received shutdown signal, stopping rover and http server")
		// Shutdown rover & http server
		r.Shutdown()

		if err := server.Shutdown(nil); err != nil {
			panic(fmt.Errorf("Error while shutting down web server: %s", err.Error()))
		}
	}()

	log.Printf("Gover server listening on port %s\n", port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
