package main

import (
	"fmt"
	"net/http"
)

const port = "903"

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	fmt.Printf("Gover server listening on port %s\n", port)
	panic(http.ListenAndServe(":903", nil))
}
