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
		Percent   float64 `json:"percent"`
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

		timePos, _ := conn.Get("time-pos")
		r.Time.Current = timePos.(float64)

		timeRemaining, _ := conn.Get("time-remaining")
		r.Time.Remaining = timeRemaining.(float64)
		r.Time.Total = r.Time.Current + r.Time.Remaining

		percent, _ := conn.Get("percent-pos")
		r.Time.Percent = percent.(float64)
	}

	jsonString, _ := json.MarshalIndent(r, "", "  ")
	return string(jsonString)
}

func statusResponse(w http.ResponseWriter, conn *mpvipc.Connection) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, getPlayerStatusJSON(conn))
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
		statusResponse(w, conn)
	})

	http.HandleFunc("/pause", func(w http.ResponseWriter, r *http.Request) {
		if err = conn.Set("pause", true); err != nil {
			log.Print(err)
		}
		statusResponse(w, conn)
	})

	http.HandleFunc("/unpause", func(w http.ResponseWriter, r *http.Request) {
		if err = conn.Set("pause", false); err != nil {
			log.Print(err)
		}
		statusResponse(w, conn)
	})

	http.HandleFunc("/toggle", func(w http.ResponseWriter, r *http.Request) {
		paused, _ := conn.Get("pause")
		if err = conn.Set("pause", !paused.(bool)); err != nil {
			log.Print(err)
		}
		statusResponse(w, conn)
	})


	log.Print("Server running on http://localhost:8080 ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	go RunServer()
	RunPlayer()
}
