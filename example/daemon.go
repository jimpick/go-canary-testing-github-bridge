package main

import (
	"fmt"
	"log"
	"net/http"

	ghbridge "github.com/jimpick/go-canary-testing-github-bridge"
)

var webhookSecretKey = []byte("ipfs_secret")

func main() {
	c := make(chan interface{})
	handler := ghbridge.GetHandler(c, webhookSecretKey)
	http.HandleFunc("/webhook", handler)
	go func() {
		for msg := range c {
			fmt.Printf("%T %v\n", msg, msg)
		}
	}()
	log.Fatal(http.ListenAndServe(":14001", nil))
}
