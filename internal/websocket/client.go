package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
)
type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte

	userID string

	rideID string

	closeOnce sync.Once
}

func (c *Client) closeSend(){
	c.closeOnce.Do(func() {
		close(c.send)
	})
}
