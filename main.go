package main

import (
	"log"
	"sync"
)

func main() {
	app := WireBuild()

	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()

		app.TgServer.Listen()
	}()

	go func() {
		defer wg.Done()

		log.Println("observing stats...")
		app.CoreObserver.Start()
	}()

	go func() {
		defer wg.Done()

		log.Println("starting notifier...")
		app.Notifier.Start()
	}()

	wg.Wait()
}
