package websocket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Mikhail-Tal63/Orbit/internal/driver"
	"github.com/Mikhail-Tal63/Orbit/internal/location"
	"github.com/Mikhail-Tal63/Orbit/middleware"
	"github.com/Mikhail-Tal63/Orbit/utils"
	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongwait       = 10 * time.Second
	pingoeriod     = (pongwait * 9) / 10
	maxMessageSize = 1 << 16
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		return middleware.IsAllowedOrigin(origin)
	},
}


func ServerWS(
	h *Hub,
	jwtSecret []byte,
	driverRepo driver.DriverRepository,
	w http.ResponseWriter,
	r *http.Request,
) {
	token := r.URL.Query().Get("token")

	if token == "" {
		http.Error(
			w,
			"missing token",
			http.StatusUnauthorized,
		)
		return
	}

	userID, err := utils.VerifyJWT(jwtSecret, token)
	if err != nil {
		http.Error(
			w,
			"invalid or expired token",
			http.StatusUnauthorized,
		)
		return
	}

	
	userUUID, err := uuid.Parse(userID.String())
	if err != nil {
		http.Error(
			w,
			"invalid user id",
			http.StatusUnauthorized,
		)
		return
	}

	
	driverRecord, err := driverRepo.GetDriverByUserId(
		r.Context(),
		userUUID,
	)
	if err != nil {
		http.Error(
			w,
			"user is not a driver",
			http.StatusForbidden,
		)
		return
	}

	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws: upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID.String(),
		driverID: driverRecord.ID.String(),
	}

	h.register <- client

	go client.WritePump()
	go client.ReadPump()
}


func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongwait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongwait))
	})

	for {
		_, row, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		_ = c.conn.SetReadDeadline(time.Now().Add(pongwait))

		var env struct {
			Type   string `json:"type"`
			RideID string `json:"ride_id"`
		}

		if err := json.Unmarshal(row, &env); err != nil {
			continue
		}
		switch env.Type {

case "join_ride":
	c.rideID = env.RideID
	c.hub.AddClient(c)
	

case "leave_ride":
	c.hub.RemoveClient(c)

case EventLocationUpdate:
	var locationUpdate location.LocationUpdate

	if err := json.Unmarshal(row, &locationUpdate); err != nil {
		log.Printf(
			"ws: invalid location update user=%s: %v",
			c.userID,
			err,
		)
		continue
	}

	err := c.hub.locationService.UpdateDriverLocation(
		context.Background(),
		c.driverID,
		locationUpdate.Latitude,
		locationUpdate.Longitude,
	)

	if err != nil {
		log.Printf(
			"ws: failed to update location driver=%s: %v",
			c.driverID,
			err,
		)
		continue
	}
case "ping":
	pong, _ := json.Marshal(
		map[string]string{
			"type": "pong",
		},
	)

	select {
	case c.send <- pong:
	default:
	}
}

	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingoeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("ws: write error user=%s: %v", c.userID, err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PongMessage, nil); err != nil {
				return
			}
		}
	}
}
