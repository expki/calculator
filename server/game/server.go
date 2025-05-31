package game

import (
	"calculator/logger"
	"net/http"

	"github.com/expki/calculator/lib/schema"
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

	// Add member to state
	g.stateLock.Lock()
	var id int = 0
	for _, member := range g.state.Members {
		if id <= member.Member.Id {
			id = 1 + member.Member.Id
		}
	}
	g.state.Members = append(g.state.Members, schema.MemberState{Member: schema.Member{Id: id}})
	g.stateLock.Unlock()

	// Create session
	session := NewSession(id, conn, &g.state, &g.stateLock)

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
	filteredMembers := make([]schema.MemberState, 0, len(g.state.Members)-1)
	for _, member := range g.state.Members {
		if member.Member.Id != id {
			filteredMembers = append(filteredMembers, member)
		}
	}
	g.stateLock.Unlock()

	// Remove session from client map
	g.clientLock.Lock()
	delete(g.clientMap, id)
	g.clientLock.Unlock()
}
