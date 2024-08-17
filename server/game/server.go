package game

import (
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
	err = NewSession(g.logger, g.sugar, conn).Wait()
	if err != nil {
		g.sugar.Errorf("Session closed with error: %v", err)
	}
}
