package websocket

import "sync"

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
