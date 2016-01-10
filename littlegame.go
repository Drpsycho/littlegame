package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

// var workdir = flag.String("workdir", "./", "work dir")

// var homeTempl *template.Template
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found.", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	io.WriteString(w, "<html><body>game server</body></html>")
}

func main() {
	flag.Parse()

	ParseMap()
	FillTileMap()

	// homeTempl = template.Must(template.ParseFiles(*workdir + "/sandbox.html"))
	go HubHandler.run()
	go StageUpdater()
	go ActionHandler()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/handler", serverWs)
	// http.Handle("/stuff/", http.StripPrefix("/stuff/", http.FileServer(http.Dir(*workdir+"/stuff/"))))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
