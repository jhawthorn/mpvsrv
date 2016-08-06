package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/static"
	"github.com/DexterLB/mpvipc"
)

const socketPath = "/tmp/mpv_socket"

func RunPlayer() {
	cmd := exec.Command("mpv", "--input-ipc-server", socketPath, "--idle")
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

func getPlayerStatus(conn *mpvipc.Connection) StatusResponse {
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

	return r
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

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("static", true)))
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, getPlayerStatus(conn))
	})
	r.POST("/pause", func(c *gin.Context) {
		if err = conn.Set("pause", true); err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, getPlayerStatus(conn))
	})
	r.POST("/unpause", func(c *gin.Context) {
		if err = conn.Set("pause", false); err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, getPlayerStatus(conn))
	})
	r.POST("/toggle", func(c *gin.Context) {
		paused, _ := conn.Get("pause")
		if err = conn.Set("pause", !paused.(bool)); err != nil {
			log.Print(err)
		}
		c.JSON(http.StatusOK, getPlayerStatus(conn))
	})


	log.Print("Server running on http://localhost:8080 ")
	r.Run(":8080")
}

func main() {
	go RunServer()
	RunPlayer()
}
