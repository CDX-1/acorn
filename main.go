package main

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"
)

//go:embed web/acorn/dist/*
var embeddedFiles embed.FS

func main() {
	distFS, err := fs.Sub(embeddedFiles, "web/acorn/dist")
	if err != nil {
		log.Fatal(err)
	}

	fileServer := http.FileServer(http.FS(distFS))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			path = "index.html"
		} else {
			path = path[1:]
		}

		_, err := fs.Stat(distFS, path)
		if err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		index, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", 500)
			return
		}
		defer index.Close()

		content, err := io.ReadAll(index)
		if err != nil {
			http.Error(w, "failed to read index.html", 500)
			return
		}

		http.ServeContent(w, r, "index.html", fsModTime(), bytes.NewReader(content))
	})

	log.Println("Running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fsModTime() (t time.Time) { return }
