package main

import (
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"time"
	"fmt"
	"flag"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/DexterLB/mpvipc"
	"github.com/elazarl/go-bindata-assetfs"
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

func dirList(f http.File) []gin.H {
	dirs, err := f.Readdir(-1)
	if err != nil {
		log.Print(err)
	}
	results := make([]gin.H, len(dirs))
	for i, d := range dirs {
		results[i] = gin.H{
			"name": d.Name(),
			"size": d.Size(),
			"mode": d.Mode(),
			"modtime": d.ModTime(),
			"is_dir": d.IsDir(),
		}
	}
	return results
}

func ServeStatic() gin.HandlerFunc {
	filesystem := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "static"}
	fileserver := http.FileServer(filesystem)
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		_, err := filesystem.Open(path)
		if err == nil {
			fileserver.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}

func RunServer(basepath string) {
	conn := mpvipc.NewConnection(socketPath)
	dir := http.Dir(basepath)

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
	r.Use(ServeStatic())
	r.GET("/browse/*path", func(c *gin.Context) {
		path := c.Param("path")
		file, err := dir.Open(path)
		if err != nil {
			log.Print(err)
		}

		c.JSON(200, dirList(file))
	})
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, getPlayerStatus(conn))
	})
	r.POST("/play", func(c *gin.Context) {
		var json struct {
			Path     string `form:"path" json:"path" binding:"required"`
		}

		if c.Bind(&json) == nil {
			var fullpath string
			log.Print(json.Path)
			if _, err = url.ParseRequestURI(json.Path); err == nil {
				fullpath = json.Path
			} else {
				log.Print(err)
				fullpath = path.Join(basepath, path.Clean(json.Path))
			}
			log.Print(fullpath)
			if _, err = conn.Call("loadfile", fullpath, "replace"); err != nil {
				log.Print(err)
			}
			c.JSON(http.StatusOK, getPlayerStatus(conn))
		}
	})
	r.POST("/seek", func(c *gin.Context) {
	    var data struct {
		Seconds float64 `form:"seconds" json:"seconds" binding:"required"`
	    }
	    if c.Bind(&data) == nil {
		if _, err = conn.Call("seek", data.Seconds, "absolute"); err != nil {
		    log.Print(err)
		}
		c.JSON(http.StatusOK, getPlayerStatus(conn))
	    }
	})
	r.POST("/stop", func(c *gin.Context) {
		if _, err = conn.Call("stop"); err != nil {
			log.Print(err)
		}
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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: mpvsrv DIR\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) !=  1 {
		usage()
	}
	go RunServer(args[0])
	RunPlayer()
}
