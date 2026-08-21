package websocket

import (
	"log"
	"sync"
)

type Hub struct {
	mu         sync.RWMutex
	rooms      map[string]map[*Client]struct{}
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]struct{}),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
	}
}
func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.AddClient(c)
		case c := <-h.unregister:
			h.RemoveClient(c)
		}
	}
}

func (h *Hub) AddClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[c.rideID]
	if !ok {
		room = make(map[*Client]struct{})
		h.rooms[c.rideID] = room
	}
	room[c] = struct{}{}
}

func (h *Hub) RemoveClient(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.rooms[c.rideID]
	if ok {
		delete(room, c)

		if len(room) == 0 {
			delete(h.rooms, c.rideID)
		}
	}

	c.closeSend()
}

func (h *Hub) BroadcastToRide(rideID string, payload []byte) {
	h.mu.RLock()

	room := h.rooms[rideID]

	clients := make([]*Client, 0, len(room))
	for c := range room {
		clients = append(clients, c)
	}

	h.mu.RUnlock()

	for _, c := range clients {
		select {
		case c.send <- payload:
		default:
			log.Printf(
				"ws: dropping slow client user=%s ws=%s",
				c.userID,
				c.rideID,
			)

			select {
			case h.unregister <- c:
			default:
			}
		}
	}
}
