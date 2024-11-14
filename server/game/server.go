package game

import (
	"calculator/logger"
	"net/http"

	"github.com/google/uuid"
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
		logger.Sugar().Error("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Create a new session
	id := uuid.New()
	session := NewSession(id, conn, &g.state)

	// Add member to state
	g.stateLock.Lock()
	g.state.AddMember(id.String(), 0, 0)
	g.stateLock.Unlock()
	// Add session to client map
	g.clientLock.Lock()
	g.clientMap[id] = session
	g.clientLock.Unlock()

	// Wait for session to close
	err = session.Wait()
	if err != nil {
		logger.Sugar().Errorf("Session closed with error: %v", err)
	}

	// Remove member from state
	g.stateLock.Lock()
	g.state.RemoveMember(id.String())
	g.stateLock.Unlock()
	// Remove session from client map
	g.clientLock.Lock()
	delete(g.clientMap, id)
	g.clientLock.Unlock()
}
