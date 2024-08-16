package game

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
	EnableCompression: false,
}

// UpgradeHandler upgrades the connection to a WebSocket
func (g *Game) UpgradeHandler() http.Handler {
	return http.HandlerFunc(g.UpgradeHandlerFunc)
}

// UpgradeHandlerFunc upgrades the connection to a WebSocket
func (g *Game) UpgradeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		g.sugar.Error("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Handle the WebSocket connection
	g.session(conn)
}

// session handles a single websocket connection
func (g *Game) session(conn *websocket.Conn) {
	defer conn.Close()

	for {
		// Read message from browser
		mt, message, err := conn.ReadMessage()
		if err != nil {
			g.sugar.Error("Read error:", err)
			break // Exit the loop on error
		}

		// Print the message to the console
		fmt.Printf("Received: %s\n", message)

		// Write message back to browser
		if err := conn.WriteMessage(mt, message); err != nil {
			g.sugar.Error("Write error:", err)
			break // Exit the loop on error
		}
	}
}
