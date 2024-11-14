package game

import (
	"bytes"
	"calculator/logger"
	"time"

	"github.com/expki/calculator/lib/encoding"
	"github.com/google/uuid"
)

const game_tick_rate = time.Second / 60

func (g *Game) gameLoop() {
	var lastEncodedState []byte
	var lastClients = make(map[uuid.UUID]struct{})
	ticker := time.NewTicker(game_tick_rate)
	for {
		select {
		case <-g.appCtx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			// continue
		}
		// skip if no clients
		g.stateLock.RLock()
		clientCount := len(g.state.Members)
		g.stateLock.RUnlock()
		if clientCount == 0 {
			continue
		}

		// Update game state
		//g.stateLock.Lock()
		//g.stateLock.Unlock()

		// Encode game state
		g.stateLock.RLock()
		msg := encoding.EncodeWithCompression(g.state)
		g.stateLock.RUnlock()

		// Send game state to all clients
		g.clientLock.RLock()
		for _, client := range g.clientMap {
			// Skip if client and state hasn't changed
			if _, ok := lastClients[client.Id]; ok && bytes.Equal(msg, lastEncodedState) {
				continue
			}
			if ok := client.TrySend(msg); !ok {
				logger.Sugar().Debugf("Client is not ready to receive message: %s", client.Id.String())
				continue
			}
		}
		lastEncodedState = msg
		lastClients = make(map[uuid.UUID]struct{})
		for _, client := range g.clientMap {
			lastClients[client.Id] = struct{}{}
		}
		g.clientLock.RUnlock()
	}
}
