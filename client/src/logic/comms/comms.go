package comms

import (
	"calculator/src/logic"
	"calculator/src/userinput"
	"context"
	"log"
	"time"

	"github.com/coder/websocket"
	"github.com/expki/calculator/lib/compression"
	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

const target_tick_rate = time.Second / 60

type Comms struct {
	conn *websocket.Conn
}

func New(port string, lgc *logic.Logic, usi *userinput.UserInput) *Comms {
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

	go func() {
		defer comms.Close()
		var userIn schema.Input
		for {
			start := time.Now()
			if usi.GetMouseX() != userIn.X || usi.GetMouseY() != userIn.Y {
				userIn.X = usi.GetMouseX()
				userIn.Y = usi.GetMouseY()
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				err = conn.Write(ctx, websocket.MessageBinary, compression.Compress(encoding.Encode(userIn)))
				if err != nil {
					log.Printf("websocket.Message.Send: %v", err)
					cancel()
					return
				}
				cancel()
			}
			end := time.Since(start)
			if end < target_tick_rate {
				time.Sleep(target_tick_rate - end)
			}
		}
	}()

	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			_, msg, err := conn.Read(ctx)
			if err != nil {
				log.Printf("websocket.Message.Read: %v", err)
				cancel()
				return
			}
			data, err := encoding.DecodeWithCompression(msg)
			if err != nil {
				log.Printf("websocket.Message.Compress: %v", err)
				cancel()
				return
			}
			state, done := lgc.LockState()
			err = encoding.Engrain(data.(map[string]any), state)
			if err != nil {
				log.Printf("encoding.Engrain: %v", err)
				cancel()
				return
			}
			cancel()
			done()
		}
	}()

	return comms
}

func (c *Comms) Close() {
	defer c.conn.CloseNow()
	c.conn.Close(websocket.StatusNormalClosure, "")
}
