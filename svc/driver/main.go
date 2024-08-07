package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

const (
	serveURL  = "localhost:8081"
	driverURI = "v1/driver/"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	locationURI = fmt.Sprintf("%sws", driverURI)
)

func main() {
	http.HandleFunc(locationURI, handleDriverConnections)
	log.Printf("HTTP server started on %s\n", serveURL)
	err := http.ListenAndServe(serveURL, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleDriverConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	in := make(chan []byte, 256)
	out := make(chan []byte, 256)

	defer ws.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go driverSvcLoop(ctx, in, out)

	for {
		select {
		case <-ctx.Done():
			log.Println("Driver service loop done")
			return
		case msg := <-out:
			err := ws.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				// If there is an error writing to the websocket, log it and continue
				// We may want to handle this differently in a production system
				log.Println("write:", err)
				continue
			}
		default:
			// TODO: check message types
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)
			in <- msg
		}

	}
}
