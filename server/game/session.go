package game

import (
	"calculator/logger"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Session struct {
	Id    uuid.UUID
	state *schema.Global
	done  chan error
	conn  *websocket.Conn
}

func NewSession(id uuid.UUID, conn *websocket.Conn, state *schema.Global) *Session {
	s := &Session{
		Id:    id,
		done:  make(chan error, 1),
		conn:  conn,
		state: state,
	}
	go s.handleInput()
	return s
}

func (s *Session) Wait() error {
	return <-s.done
}

// handleInput handles client inputs
func (s *Session) handleInput() {
	for {
		// Read message from browser
		mt, message, err := s.conn.ReadMessage()
		if err != nil {
			logger.Sugar().Error("Read error:", err)
			close(s.done)
			s.conn.Close()
			return // End the session
		}
		// Handle message
		func() {
			if mt != websocket.BinaryMessage {
				logger.Sugar().Errorf("Invalid message type: %d", mt)
				return
			}
			// Decode the message
			data, err := encoding.DecodeWithCompression(message)
			if err != nil {
				logger.Sugar().Errorf("Decode error: %v", err)
				return
			}
			var userIn schema.Input
			err = encoding.Engrain(data.(map[string]any), &userIn)
			if err != nil {
				logger.Sugar().Errorf("Decode client message")
				return
			}
			me := s.state.GetMember(s.Id.String())
			if me == nil {
				logger.Sugar().Errorf("Member not found: %s", s.Id.String())
				return
			}
			me.SetX(userIn.X)
			me.SetY(userIn.Y)
		}()
	}
}
