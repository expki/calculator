package game

import (
	"fmt"
	"time"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const tick_rate = 5 * time.Second

type Session struct {
	Id     uuid.UUID
	state  *schema.Global
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	done   chan error
	conn   *websocket.Conn
}

func NewSession(logger *zap.Logger, sugar *zap.SugaredLogger, conn *websocket.Conn, state *schema.Global) *Session {
	s := &Session{
		Id:     uuid.New(),
		logger: logger,
		sugar:  sugar,
		done:   make(chan error),
		conn:   conn,
	}
	go s.handleInput(conn)
	go s.handleOutput(conn)
	return s
}

func (s *Session) Wait() error {
	return <-s.done
}

// handleInput handles client inputs
func (g *Session) handleInput(conn *websocket.Conn) {
	defer conn.Close()
	for {
		// Read message from browser
		mt, message, err := conn.ReadMessage()
		if err != nil {
			g.sugar.Error("Read error:", err)
			break // Exit the loop on error
		}
		switch mt {
		case websocket.BinaryMessage:
			// Decode the message
			data, _ := encoding.Decode(message)
			fmt.Println(data)
		default:
			g.sugar.Errorf("Invalid message type: %d", mt)
		}
	}
}

// handleOutput keeps the client up to date with latest state
func (g *Session) handleOutput(conn *websocket.Conn) {
	defer conn.Close()
	for {
		start := time.Now()
		//Send the state to the client
		if err := conn.WriteMessage(websocket.BinaryMessage, encoding.Encode(g.state)); err != nil {
			g.sugar.Error("Write error:", err)
			break
		}
		end := time.Since(start)
		if end < tick_rate {
			time.Sleep(tick_rate - end)
		}
	}
}
