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
	Idle   bool   `json:"idle"`
	Paused bool   `json:"paused"`
	Path   string `json:"path"`
	Title  string `json:"title"`
	Time   struct {
		Current   float64 `json:"current"`
		Remaining float64 `json:"remaining"`
		Total     float64 `json:"total"`
	} `json:"time"`
}

func getPlayerStatusJSON(conn *mpvipc.Connection) string {
	var r StatusResponse
	idle, _ := conn.Get("idle")
	r.Idle = idle.(bool)

	if !r.Idle {
		paused, _ := conn.Get("pause")
		r.Paused = paused.(bool)

		path, _ := conn.Get("path")
		r.Path = path.(string)

		title, _ := conn.Get("media-title")
		r.Title = title.(string)
	}

	jsonString, _ := json.MarshalIndent(r, "", "  ")
	return string(jsonString)
}

func RunServer() {
	conn := mpvipc.NewConnection(socketPath)

	var err error
	for i := 0; i < 1000; i++ {
		err = conn.Open()
		if err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, getPlayerStatusJSON(conn))
	})

	log.Print("Server running on http://localhost:8080 ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	go RunServer()
	RunPlayer()
}
