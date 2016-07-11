package main

import (
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/DexterLB/mpvipc"
)

func RunPlayer() {
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

func getPlayerStatusJSON(conn *mpvipc.Connection) string {
	return "{}"
}

func RunServer() {
	waitForSocket()

	conn := mpvipc.NewConnection("/tmp/mpv_socket")
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, getPlayerStatusJSON(conn))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	go RunServer()
	RunPlayer()
}
