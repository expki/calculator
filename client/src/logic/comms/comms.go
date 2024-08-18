package comms

import (
	"calculator/src/logic"
	"context"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/expki/calculator/lib/encoding"
)

type Comms struct {
	conn *websocket.Conn
}

func New(port string, lgc *logic.Logic) *Comms {
	uri, err := getWebsocketURL(port)
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
		conn: conn,
	}

	/*go func() {
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
	}()*/

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_, msg, err := conn.Read(ctx)
			if err != nil {
				log.Printf("websocket.Message.Read: %v", err)
				return
			}
			data, _ := encoding.Decode(msg)
			state, done := lgc.LockState()
			err = encoding.Engrain(data.(map[string]any), state)
			if err != nil {
				log.Printf("encoding.Engrain: %v", err)
				return
			}
			done()
		}
	}()

	return comms
}

func (c *Comms) Close() {
	defer c.conn.CloseNow()
	c.conn.Close(websocket.StatusNormalClosure, "")
}
