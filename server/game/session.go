package game

import (
	"calculator/logger"
	"sync"
	"time"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"

	"github.com/gorilla/websocket"
)

type Session struct {
	Id        int
	stateLock *sync.RWMutex
	state     *schema.StateState
	done      chan error
	conn      *websocket.Conn
	send      chan []byte
}

func NewSession(id int, conn *websocket.Conn, state *schema.StateState, stateLock *sync.RWMutex) *Session {
	s := &Session{
		Id:        id,
		done:      make(chan error, 1),
		conn:      conn,
		stateLock: stateLock,
		state:     state,
		send:      make(chan []byte, 1),
	}
	go s.handlePing()
	go s.handleInput()
	go s.handleOutput()
	return s
}

func (s *Session) Wait() error {
	return <-s.done
}

func (s *Session) handlePing() {
	// Set up a ticker to send ping messages periodically
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-s.done:
			ticker.Stop()
			return
		case <-ticker.C:
			// Send a ping message
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Sugar().Errorf("Ping error %d: %v", s.Id, err)
				close(s.done)
				s.conn.Close()
			}
		}
	}
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
			if mt == websocket.PongMessage {
				logger.Sugar().Debugf("Pong message: %d", s.Id)
				return // Ignore pong messages
			}
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
			s.stateLock.Lock()
			s.handleUserInput(userIn)
			s.stateLock.Unlock()
		}()
	}
}

// Send sends a message to the client
func (s *Session) handleOutput() {
	for {
		select {
		case <-s.done:
			return
		case msg := <-s.send:
			err := s.conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				logger.Sugar().Error("Write error:", err)
			}
		}
	}
}

func (s *Session) TrySend(msg []byte) bool {
	select {
	case s.send <- msg:
		return true
	default:
		return false
	}
}
