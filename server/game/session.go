package game

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Session struct {
	Id     uuid.UUID
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	done   chan error
	conn   *websocket.Conn
}

func NewSession(logger *zap.Logger, sugar *zap.SugaredLogger, conn *websocket.Conn) *Session {
	s := &Session{
		Id:     uuid.New(),
		logger: logger,
		sugar:  sugar,
		done:   make(chan error),
		conn:   conn,
	}
	go s.handleInput(conn)
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
		default:

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
