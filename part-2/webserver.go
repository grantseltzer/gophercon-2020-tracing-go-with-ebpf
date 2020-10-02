package main

import (
	"fmt"
	"log"
	"net/http"
)

//go:noinline
func handlerFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Received connection from %s", r.Host)
}

func main() {
	http.HandleFunc("/", handlerFunction)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
