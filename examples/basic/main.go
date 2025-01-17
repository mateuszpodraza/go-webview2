package main

import (
	"log"

	"github.com/mattpodraza/webview2/v2/pkg/webview2"
)

func main() {
	wv, err := webview2.New(
		webview2.WithTitle("The Go Programming Language"),
		webview2.WithSize(800, 600),
		webview2.WithURL("https://golang.org"),
	)

	if err != nil {
		log.Fatalf("Failed to create webview2: %v", err)
	}

	if err := wv.Run(); err != nil {
		log.Fatalf("Failed while running webview: %v", err)
	}
}
