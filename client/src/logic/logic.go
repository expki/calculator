package logic

import (
	"bytes"
	"calculator/src/types"
	"context"
	"log"
	"syscall/js"
	"time"

	"github.com/coder/websocket"
	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

const game_tick_rate = time.Second / 60

type Logic struct {
	url  string
	conn *websocket.Conn
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
	js.CopyBytesToJS(sharedArray, encoding.Encode(types.State{}))

	// Connect to the server
	logic.connect()

	// Client (self) → Server
	var lastMsg []byte
	notify := func(input schema.Input) {
		msg := encoding.EncodeWithCompression(input)
		if bytes.Equal(lastMsg, msg) {
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = logic.conn.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			log.Fatalf("websocket.Message.Send exception: %v", err)
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
		var pref float32
		for {
			func() {
				t, msg, err := logic.conn.Read(context.Background())
				if err != nil {
					log.Fatalf("websocket.Message.Read exception: %v", err)
					return
				}
				switch t {
				case websocket.MessageBinary:
					// continue
				default:
					log.Fatalf("unexpected message type: %v", t)
				}
				prefStart := time.Now()
				data, err := encoding.DecodeWithCompression(msg)
				if err != nil {
					log.Fatalf("websocket.Message.Compress exception: %v", err)
					return
				}
				var global schema.Global
				err = encoding.Engrain(data.(map[string]any), &global)
				if err != nil {
					log.Fatalf("encoding.Engrain exception: %v", err)
					return
				}
				state := types.State{
					Global:   global,
					CpuLogic: pref,
				}
				stateData := encoding.Encode(state)
				if !bytes.Equal(lastStateData, stateData) {
					js.CopyBytesToJS(sharedArray, stateData)
					lastStateData = stateData
					pref = float32(time.Duration(time.Since(prefStart).Nanoseconds())) / float32(game_tick_rate)
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
