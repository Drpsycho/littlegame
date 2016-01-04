package main

import (
	"flag"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

var addr = flag.String("addr", "10.0.0.4:8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("/home/drpsycho/js/jsgame/sandbox.html"))
var upgrader = websocket.Upgrader{} // use default options

func handler(rw http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(rw, "ws://"+req.Host+"/handler")
}

func main() {

	go HubHandler.run()
	http.HandleFunc("/", handler)
	http.HandleFunc("/handler", serverWs)
	http.Handle("/stuff/", http.StripPrefix("/stuff/", http.FileServer(http.Dir("/home/drpsycho/js/jsgame/stuff/"))))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
