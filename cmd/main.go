package main

import (
	"log"
)

func main() {
	app, err := NewApplication()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	defer app.Shutdown()

	if err := app.Run(); err != nil {
		log.Fatalf("Application run error: %v", err)
	}
}
