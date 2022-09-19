package main

import (
	"github.com/gorilla/websocket"
)

// client refers to a single user
type client struct {
	// socket is the web socket for this client
	socket *websocket.Conn
	// send is a channel on which msgs are sent
	send chan []byte
	// room is the room this client is chatting in
	room *room
}

// reading from the web-socket

// msgs are read from socket using ReadMessage and are then sent to forward channel in Room
func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			// print err and return
			return
		}

		c.room.forward <- msg
	}

}

// writing to the web-socket
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			//print err
			return
		}
	}
}
