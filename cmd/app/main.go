package main

import (
	"log"
	"tmballNews/internal/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed start app main.go: %s", err.Error())
	}

	if err := a.Run(); err != nil {
		log.Fatalf("failed to run app main.go: %s", err.Error())
	}
}
