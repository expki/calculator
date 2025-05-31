package game

import (
	"calculator/logger"
	"encoding/json"
	"time"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

const game_tick_rate = time.Second / 60

func (g *Game) gameLoop() {
	var (
		lastClientCount int
		lastState       string
	)
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
		if clientCount != lastClientCount {
			logger.Sugar().Debugf("Number of clients: %d", lastClientCount)
			lastClientCount = clientCount
		}
		if clientCount == 0 {
			continue
		}

		// Encode game state
		g.stateLock.RLock()
		state := g.state.State()
		g.stateLock.RUnlock()
		currentState, _ := json.Marshal(state)
		if lastState != string(currentState) {
			logger.Logger().Debug(string(currentState))
			lastState = string(currentState)
		} else {
			// skip if no change
			continue
		}

		// Send game state to all clients
		g.clientLock.RLock()
		for _, client := range g.clientMap {
			msg := encoding.EncodeWithCompression(schema.PersonalizedState{
				Id:    client.Id,
				State: state.WithoutMember(client.Id),
			})
			if ok := client.TrySend(msg); !ok {
				logger.Sugar().Debugf("Client is not ready to receive message: %d", client.Id)
				continue
			}
		}
		g.clientLock.RUnlock()
	}
}
