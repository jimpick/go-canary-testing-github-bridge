package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q\n", r.URL.Path)
	})
	log.Fatal(http.ListenAndServe(":14001", nil))
}
