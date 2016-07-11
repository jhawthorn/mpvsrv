package main

import (
	"fmt"
	"log"

	"github.com/DexterLB/mpvipc"
)

func main() {
	conn := mpvipc.NewConnection("/tmp/mpv_rpc")
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	events, stopListening := conn.NewEventListener()

	path, err := conn.Get("path")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("current file playing: %s", path)

	err = conn.Set("pause", true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("paused!")

	_, err = conn.Call("observe_property", 42, "volume")
	if err != nil {
		fmt.Print(err)
	}

	go func() {
		conn.WaitUntilClosed()
		stopListening <- struct{}{}
	}()

	for event := range events {
		if event.ID == 42 {
			log.Printf("volume now is %f", event.Data.(float64))
		} else {
			log.Printf("received event: %s", event.Name)
		}
	}

	log.Printf("mpv closed socket")
}
