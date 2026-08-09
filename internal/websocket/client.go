package websocket

type Client struct {
	hub  *Hub
	conn *Websocket
}
