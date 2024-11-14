package game

import (
	"bytes"
	"calculator/logger"
	"sync"
	"time"

	"github.com/expki/calculator/lib/encoding"
	"github.com/gorilla/websocket"
)

const game_tick_rate = time.Second / 60

func (g *Game) gameLoop() {
	var lastEncodedState []byte
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
		g.stateLock.Unlock()

		// Skip if state hasn't changed
		if bytes.Equal(msg, lastEncodedState) {
			continue
		}

		// Send game state to all clients
		g.clientLock.RLock()
		var wg sync.WaitGroup
		wg.Add(len(g.clientMap))
		for _, client := range g.clientMap {
			go func() {
				err := client.conn.WriteMessage(websocket.BinaryMessage, msg)
				if err != nil {
					logger.Sugar().Error("Write error %s:", client.Id.String(), err)
				}
				wg.Done()
			}()
		}
		g.clientLock.RUnlock()
		wg.Wait()
	}
}
