package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var workdir = flag.String("workdir", "./", "work dir")

var homeTempl *template.Template
var upgrader = websocket.Upgrader{}

func handler(rw http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(rw, "ws://"+req.Host+"/handler")
}

func main() {
	flag.Parse()

	homeTempl = template.Must(template.ParseFiles(*workdir+"/sandbox.html"))
	go HubHandler.run()
	http.HandleFunc("/", handler)
	http.HandleFunc("/handler", serverWs)
	http.Handle("/stuff/", http.StripPrefix("/stuff/", http.FileServer(http.Dir(*workdir+"/stuff/"))))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
