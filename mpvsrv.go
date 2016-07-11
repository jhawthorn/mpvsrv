package main

import (
	"fmt"
	"log"
	"time"
	"os/exec"

	"github.com/DexterLB/mpvipc"
)

func NewPlayer() {
	cmd := exec.Command("mpv", "--input-ipc-server=/tmp/mpv_socket", "--idle", "--force-window")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

func waitForSocket() {
	// FIXME: this should sleep less and check for the socket to be created
	time.Sleep(time.Second)
}

func main() {
	go NewPlayer()
	waitForSocket()

	conn := mpvipc.NewConnection("/tmp/mpv_socket")
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	events, stopListening := conn.NewEventListener()

	path, err := conn.Get("path")
	if err != nil {
		// log.Fatal(err)
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
