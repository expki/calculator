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
	// Create a new session
	session := NewSession(g.logger, g.sugar, conn, &g.state)
	// Add session to state
	g.state.AddMember(session.Id.String(), 0, 0)
	// Wait for session to close
	err = session.Wait()
	if err != nil {
		g.sugar.Errorf("Session closed with error: %v", err)
	}
	g.state.RemoveMember(session.Id.String())
}
