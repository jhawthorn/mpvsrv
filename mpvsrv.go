package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/DexterLB/mpvipc"
)

const socketPath = "/tmp/mpv_socket"

func RunPlayer() {
	cmd := exec.Command("mpv", "--input-ipc-server", socketPath, "--idle", "--force-window")
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	log.Printf("Command finished with error: %v", err)
}

type StatusResponse struct {
	Paused bool   `json:"paused"`
	Path   string `json:"path"`
	Title  string `json:"title"`
}

func waitForSocket() {
	// FIXME: this should sleep less and check for the socket to be created
	time.Sleep(time.Second)
}

func getPlayerStatusJSON(conn *mpvipc.Connection) string {
	response := &StatusResponse{
		Paused: false,
		Path:   "/foo/bar",
		Title:  "Foobar",
	}
	json, _ := json.MarshalIndent(response, "", "  ")
	return string(json)
}

func RunServer() {
	waitForSocket()

	conn := mpvipc.NewConnection(socketPath)
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, getPlayerStatusJSON(conn))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	go RunServer()
	RunPlayer()
}
