package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client is a chatting users
type client struct {
	// socket is Websocket for the users
	// send is chanel to send  messages
	// room is a chat room the user is in
	socket   *websocket.Conn
	send     chan *message
	room     *room
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		err := c.socket.ReadJSON(&msg)
		if err != nil {
			return
		}
		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		c.room.forward <- msg
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			return
		}
	}
	c.socket.Close()
}
