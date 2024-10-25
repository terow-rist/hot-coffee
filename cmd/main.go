package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/{SomeShit}", jopaHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func jopaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Dog shit")
}
