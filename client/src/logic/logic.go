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

func New(port string, sharedArray js.Value) *Logic {
	// Set initial state
	js.CopyBytesToJS(sharedArray, encoding.Encode(types.LocalState{}))

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
	userInput := LocalInput{}
	window.Set("handleInput", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return nil
		}
		arg := args[0]
		if value := arg.Get("mouseleft"); !value.IsUndefined() {
			userInput.mouseleft = value.Bool()
		}
		if value := arg.Get("mousex"); !value.IsUndefined() {
			userInput.mousex = value.Int()
		}
		if value := arg.Get("mousey"); !value.IsUndefined() {
			userInput.mousey = value.Int()
		}
		if value := arg.Get("width"); !value.IsUndefined() {
			userInput.width = value.Int()
		}
		if value := arg.Get("height"); !value.IsUndefined() {
			userInput.height = value.Int()
		}
		if value := arg.Get("x"); !value.IsUndefined() {
			userInput.x = value.Int()
		}
		if value := arg.Get("y"); !value.IsUndefined() {
			userInput.y = value.Int()
		}
		// row 1
		if value := arg.Get("clear"); !value.IsUndefined() {
			userInput.clear = value.Bool()
		}
		if value := arg.Get("bracket"); !value.IsUndefined() {
			userInput.bracket = value.Bool()
		}
		if value := arg.Get("percentage"); !value.IsUndefined() {
			userInput.percentage = value.Bool()
		}
		if value := arg.Get("divide"); !value.IsUndefined() {
			userInput.divide = value.Bool()
		}
		// row 2
		if value := arg.Get("seven"); !value.IsUndefined() {
			userInput.seven = value.Bool()
		}
		if value := arg.Get("eight"); !value.IsUndefined() {
			userInput.eight = value.Bool()
		}
		if value := arg.Get("nine"); !value.IsUndefined() {
			userInput.nine = value.Bool()
		}
		if value := arg.Get("times"); !value.IsUndefined() {
			userInput.times = value.Bool()
		}
		// row 3
		if value := arg.Get("four"); !value.IsUndefined() {
			userInput.four = value.Bool()
		}
		if value := arg.Get("five"); !value.IsUndefined() {
			userInput.five = value.Bool()
		}
		if value := arg.Get("six"); !value.IsUndefined() {
			userInput.six = value.Bool()
		}
		if value := arg.Get("minus"); !value.IsUndefined() {
			userInput.minus = value.Bool()
		}
		// row 4
		if value := arg.Get("one"); !value.IsUndefined() {
			userInput.one = value.Bool()
		}
		if value := arg.Get("two"); !value.IsUndefined() {
			userInput.two = value.Bool()
		}
		if value := arg.Get("three"); !value.IsUndefined() {
			userInput.three = value.Bool()
		}
		if value := arg.Get("plus"); !value.IsUndefined() {
			userInput.plus = value.Bool()
		}
		// row 5
		if value := arg.Get("negate"); !value.IsUndefined() {
			userInput.negate = value.Bool()
		}
		if value := arg.Get("zero"); !value.IsUndefined() {
			userInput.zero = value.Bool()
		}
		if value := arg.Get("decimal"); !value.IsUndefined() {
			userInput.decimal = value.Bool()
		}
		if value := arg.Get("equals"); !value.IsUndefined() {
			userInput.equals = value.Bool()
		}
		notify(userInput.Translate())
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
				var global schema.State
				err = encoding.Engrain(data.(map[string]any), &global)
				if err != nil {
					log.Fatalf("encoding.Engrain exception: %v", err)
					return
				}
				state := types.LocalState{
					State:   global,
					Id:      0,
					CpuLoad: pref,
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
