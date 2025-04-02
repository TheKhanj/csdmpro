package main

import (
	"log"
	"sync"
)

func main() {
	app := WireBuild()

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()

		app.TgServer.Listen()
	}()

	go func() {
		defer wg.Done()

		log.Println("observing stats...")
		app.CoreObserver.Start()
	}()

	wg.Wait()
}
