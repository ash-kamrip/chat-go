package main

import (
	"log"
	"net/http"

	trace "github.com/ash-kamrip/tracer"
	"github.com/gorilla/websocket"
)

type room struct {
	// forward is a channel that holds the incoming message
	// msgs that needs to be forwarded to other client
	forward chan []byte

	// join channel adds the users to the room
	join chan *client
	// leave channel removes the users from the room
	leave chan *client

	// clients keep a record of users present in the room
	clients map[*client]bool

	// tracer will receive trace information of activity in the room.
	tracer trace.Tracer
}

const (
	socketBuffersize  = 1024
	messageBuffersize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBuffersize, WriteBufferSize: messageBuffersize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// the request must be upgraded to get a web-socket
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	// create a client
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBuffersize),
		room:   r,
	}

	// add the client to the room
	r.join <- client
	// tidy up when the client leave
	defer func() { r.leave <- client }()

	// the write method here is supposed to run in background thread
	go client.write()
	// blocks the main thread and read messages.
	client.read()

}

// this code will run indefinitely in the background
func (r *room) run() {

	for {
		// keeps watch on the 3 of our channels
		// we wants synchronized access to the clients map in room, hence this implementatioResponseWriter
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.tracer.Trace("New client joined! 😁")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("client left 😔 ")
		// msgs received over forward channel is forwarded to all the client(send channel of the client) in the room
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace("-- sent to client")
			}
		}
	}

}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}
