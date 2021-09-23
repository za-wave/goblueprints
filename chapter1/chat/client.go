package main

import (
	"github.com/gorilla/websocket"
)

// client is a chatting users
type client struct {
	// socket is Websocket for the users
	// send is chanel to send  messages
	// room is a chat room the user is in
	socket *websocket.Conn
	send   chan []byte
	room   *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
