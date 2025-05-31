package game

import (
	"bytes"
	"calculator/logger"
	"encoding/json"
	"time"

	"github.com/expki/calculator/lib/encoding"
)

const game_tick_rate = time.Second / 60

func (g *Game) gameLoop() {
	var lastEncodedState []byte
	var lastClients = make(map[int]struct{})
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
			logger.Logger().Debug("no clients...")
			continue
		}

		// Encode game state
		g.stateLock.RLock()
		msg := encoding.EncodeWithCompression(g.state.State())
		g.stateLock.RUnlock()

		ssss, _ := json.Marshal(g.state.State())
		logger.Logger().Debug(string(ssss))

		// Send game state to all clients
		g.clientLock.RLock()
		for _, client := range g.clientMap {
			// Skip if client and state hasn't changed
			if _, ok := lastClients[client.Id]; ok && bytes.Equal(msg, lastEncodedState) {
				continue
			}
			if ok := client.TrySend(msg); !ok {
				logger.Sugar().Debugf("Client is not ready to receive message: %d", client.Id)
				continue
			}
		}
		lastEncodedState = msg
		lastClients = make(map[int]struct{})
		for _, client := range g.clientMap {
			lastClients[client.Id] = struct{}{}
		}
		g.clientLock.RUnlock()
	}
}
