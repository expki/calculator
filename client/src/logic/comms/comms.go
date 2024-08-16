package comms

import (
	"context"
	"log"
	"time"

	"github.com/coder/websocket"
)

type Comms struct {
	conn      *websocket.Conn
	readChan  chan []byte
	writeChan chan []byte
}

func New() *Comms {
	uri, err := getWebsocketURL()
	if err != nil {
		log.Fatalf("websocketURL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, uri, nil)
	if err != nil {
		log.Fatalf("websocket.Dial: %v", err)
	}

	comms := &Comms{
		conn:      conn,
		readChan:  make(chan []byte),
		writeChan: make(chan []byte),
	}

	go func() {
		defer comms.Close()
		for {
			msg := <-comms.writeChan
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := conn.Write(ctx, websocket.MessageBinary, msg)
			if err != nil {
				log.Printf("websocket.Message.Send: %v", err)
				return
			}
		}
	}()

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_, msg, err := conn.Read(ctx)
			if err != nil {
				log.Printf("websocket.Message.Read: %v", err)
				return
			}
			comms.readChan <- msg
		}
	}()

	return comms
}

func (c *Comms) Close() {
	defer c.conn.CloseNow()
	c.conn.Close(websocket.StatusNormalClosure, "")
}

func (c *Comms) TryRead() (msg []byte, ok bool) {
	select {
	case msg = <-c.readChan:
		return msg, true
	default:
		return nil, false
	}
}

func (c *Comms) TryWrite(msg []byte) bool {
	select {
	case c.writeChan <- msg:
		return true
	default:
		return false
	}
}
