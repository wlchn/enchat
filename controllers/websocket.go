package controllers

import (
	"encoding/json"
	"net/http"

	"enchat/models"
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
)

// WebSocketController handles WebSocket requests.
type WebSocketController struct {
	MainController
}

// Get method handles GET requests for WebSocketController.
func (c *WebSocketController) Get() {
	// Safe check.
	email := c.GetString("email")
	if len(email) == 0 {
		c.Redirect("/", 302)
		return
	}

	c.TplName = "chatroom.html"
	c.Data["IsWebSocket"] = true
	c.Data["UserName"] = email
}

// Join method handles WebSocket requests for WebSocketController.
func (c *WebSocketController) Join() {
	email := c.GetString("email")
	if len(email) == 0 {
		c.Redirect("/", 302)
		return
	}

	// Upgrade from http request to WebSocket.
	ws, err := websocket.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	// Join chat room.
	Join(email, ws)
	defer Leave(email)

	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		if err != nil {
			return
		}
		publish <- newEvent(models.EVENT_MESSAGE, email, string(p))
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event models.Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- sub.Value.(Subscriber).Name
			}
		}
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func sendMsg(event models.Event, toUserEmail string) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if sub.Value.(Subscriber).Name == toUserEmail {
			if ws != nil {
				if ws.WriteMessage(websocket.TextMessage, data) != nil {
					// User disconnected.
					unsubscribe <- sub.Value.(Subscriber).Name
				}
			}
		}
	}
}
