package logic

import (
	"bytes"
	"calculator/src/types"
	"context"
	"log"
	"sync"
	"syscall/js"
	"time"

	"github.com/coder/websocket"
	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

type Logic struct {
	connectLock sync.Mutex
	url         string
	conn        *websocket.Conn
}

type UserInput struct {
	width     int
	height    int
	key       map[string]bool
	mouseleft bool
	mousex    int
	mousey    int
}

func New(port string, sharedArray js.Value) *Logic {
	// Get server websocket url
	uri, err := getWebsocketURL(port)
	if err != nil {
		log.Fatalf("websocketURL: %v", err)
	}

	// Instantiate the logic module
	logic := &Logic{
		url:  uri,
		conn: nil,
	}

	// Set initial state
	js.CopyBytesToJS(sharedArray, encoding.EncodeWithCompression(types.State{}))

	// Connect to the server
	logic.connect()

	// Client (self) → Server
	var lastMsg []byte
	notify := func(input schema.Input) {
		msg := encoding.EncodeWithCompression(input)
		if bytes.Equal(lastMsg, msg) {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		logic.connectLock.Lock()
		err = logic.conn.Write(ctx, websocket.MessageBinary, msg)
		logic.connectLock.Unlock()
		if err != nil {
			log.Printf("websocket.Message.Send exception: %v", err)
			return
		}
		lastMsg = msg
	}
	window := js.Global()
	userInput := UserInput{key: make(map[string]bool)}
	window.Set("handleInput", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return nil
		}
		arg := args[0]
		if value := arg.Get("width"); !value.IsUndefined() {
			userInput.width = value.Int()
		}
		if value := arg.Get("height"); !value.IsUndefined() {
			userInput.height = value.Int()
		}
		if value := arg.Get("keyup"); !value.IsUndefined() {
			userInput.key[value.String()] = false
		}
		if value := arg.Get("keydown"); !value.IsUndefined() {
			userInput.key[value.String()] = true
		}
		if value := arg.Get("mouseleft"); !value.IsUndefined() {
			userInput.mouseleft = value.Bool()
		}
		if value := arg.Get("mousex"); !value.IsUndefined() {
			userInput.mousex = value.Int()
		}
		if value := arg.Get("mousey"); !value.IsUndefined() {
			userInput.mousey = value.Int()
		}
		notify(schema.Input{X: userInput.mousex, Y: userInput.mousey})
		return nil
	}))

	// Server → Client (self)
	go func() {
		var lastStateData []byte
		for {
			func() {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
				defer cancel()
				_, msg, err := logic.conn.Read(ctx)
				if err != nil {
					log.Printf("websocket.Message.Read exception: %v", err)
					logic.connect()
					return
				}
				data, err := encoding.DecodeWithCompression(msg)
				if err != nil {
					log.Printf("websocket.Message.Compress exception: %v", err)
					return
				}
				var global schema.Global
				err = encoding.Engrain(data.(map[string]any), &global)
				if err != nil {
					log.Printf("encoding.Engrain exception: %v", err)
					return
				}
				state := types.State{
					Global: global,
				}
				stateData := encoding.EncodeWithCompression(state)
				if !bytes.Equal(lastStateData, stateData) {
					js.CopyBytesToJS(sharedArray, stateData)
					lastStateData = stateData
				}
			}()
		}
	}()

	// Return logic module
	return logic
}

func (l *Logic) Close() {
	l.conn.Close(websocket.StatusGoingAway, "normal closure")
}
