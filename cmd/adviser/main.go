package main

import (
	"log"

	"bonefabric/adviser/app"
)

func main() {
	application := app.New()
	var err error

	if err = application.Init(); err != nil {
		log.Fatalf("failed to init application: %s", err)
	}
	application.Run()
}
