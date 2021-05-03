package main

import (
	"log"

	"github.com/mattpodraza/webview2/pkg/webview2"
)

func main() {
	wv, err := webview2.New(
		webview2.WithTitle("Minimal webview example"),
		webview2.WithSize(800, 600, webview2.HintNone),
		webview2.WithURL("https://en.m.wikipedia.org/wiki/Main_Page"),
	)

	if err != nil {
		log.Fatalf("Failed to create webview2: %v", err)
	}

	// Error handling omitted for brevity
	_ = wv.Run()
}
