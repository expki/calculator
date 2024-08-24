package game

import (
	"bytes"
	"time"

	"github.com/expki/calculator/lib/compression"
	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const tick_rate = time.Second / 60

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
		state:  state,
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
	defer func() {
		select {
		case <-g.done:
			return
		default:
			close(g.done)
		}
		conn.Close()
	}()
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
			out, err := compression.Decompress(message)
			if err != nil {
				g.sugar.Errorf("Decompress client message")
				continue
			}
			data, _ := encoding.Decode(out)
			var userIn schema.Input
			err = encoding.Engrain(data.(map[string]any), &userIn)
			if err != nil {
				g.sugar.Errorf("Decode client message")
				continue
			}
			me := g.state.GetMember(g.Id.String())
			me.SetX(userIn.X)
			me.SetY(userIn.Y)
		default:
			g.sugar.Errorf("Invalid message type: %d", mt)
		}
	}
}

// handleOutput keeps the client up to date with latest state
func (g *Session) handleOutput(conn *websocket.Conn) {
	defer func() {
		select {
		case <-g.done:
			return
		default:
			close(g.done)
		}
		conn.Close()
	}()
	var lastEncodedState []byte
	for {
		start := time.Now()
		// Encode state
		encodedState := encoding.Encode(g.state)
		if !bytes.Equal(encodedState, lastEncodedState) {
			//Compress message
			msg := compression.Compress(encodedState)
			g.sugar.Debugf("message compression: %.2f", float64(len(msg))/float64(len(encodedState)))
			//Send the state to the client
			if err := conn.WriteMessage(websocket.BinaryMessage, msg); err != nil {
				g.sugar.Error("Write error:", err)
				break
			}
			lastEncodedState = encodedState
		}

		end := time.Since(start)
		if end < tick_rate {
			time.Sleep(tick_rate - end)
		}
	}
}
