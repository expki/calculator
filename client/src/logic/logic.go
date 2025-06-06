package logic

import (
	"bytes"
	"calculator/src/types"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
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

func New(ctx context.Context, cancel context.CancelFunc, port string, sharedArray js.Value) (logic *Logic) {
	window := js.Global()

	// Instantiate the logic module
	logic = &Logic{}
	var err error

	// Set initial state
	js.CopyBytesToJS(sharedArray, encoding.Encode(types.LocalState{}))

	// Get server websocket url
	logic.url, err = getWebsocketURL(port)
	if err != nil {
		log.Printf("error websocketURL: %v", err)
		cancel()
		return
	}

	// Connect to the server
	err = logic.connect(ctx)
	if err != nil {
		log.Printf("error connect: %v", err)
		cancel()
		return
	}

	// Client (self) → Server
	var lastMsg []byte
	notify := func(input schema.Input) {
		defer func() {
			if r := recover(); r != nil {
				err := errors.New("Tried to write closed connection")
				fmt.Printf("Recovered. %s: %v\n", err.Error(), r)
				cancel()
			}
		}()
		msg := encoding.EncodeWithCompression(input)
		if bytes.Equal(lastMsg, msg) {
			return
		}
		err := logic.conn.Write(ctx, websocket.MessageBinary, msg)
		if err != nil {
			log.Printf("error send: %v", err)
			cancel()
			return
		}
		lastMsg = msg
	}

	localInputLock := &sync.Mutex{}
	localInput := &LocalInput{}
	window.Set("handleInput", js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 || args[0].IsUndefined() {
			return nil
		}
		arg := args[0]
		local := LocalInput{}
		if value := arg.Get("mouseleft"); !value.IsUndefined() {
			local.mouseleft = value.Bool()
		}
		if value := arg.Get("mousex"); !value.IsUndefined() {
			local.mousex = value.Int()
		}
		if value := arg.Get("mousey"); !value.IsUndefined() {
			local.mousey = value.Int()
		}
		if value := arg.Get("width"); !value.IsUndefined() {
			local.width = value.Int()
		}
		if value := arg.Get("height"); !value.IsUndefined() {
			local.height = value.Int()
		}
		if value := arg.Get("x"); !value.IsUndefined() {
			local.x = value.Int()
		}
		if value := arg.Get("y"); !value.IsUndefined() {
			local.y = value.Int()
		}
		// row 1
		if value := arg.Get("clear"); !value.IsUndefined() {
			local.clear = value.Bool()
		}
		if value := arg.Get("bracket"); !value.IsUndefined() {
			local.bracket = value.Bool()
		}
		if value := arg.Get("percentage"); !value.IsUndefined() {
			local.percentage = value.Bool()
		}
		if value := arg.Get("divide"); !value.IsUndefined() {
			local.divide = value.Bool()
		}
		// row 2
		if value := arg.Get("seven"); !value.IsUndefined() {
			local.seven = value.Bool()
		}
		if value := arg.Get("eight"); !value.IsUndefined() {
			local.eight = value.Bool()
		}
		if value := arg.Get("nine"); !value.IsUndefined() {
			local.nine = value.Bool()
		}
		if value := arg.Get("times"); !value.IsUndefined() {
			local.times = value.Bool()
		}
		// row 3
		if value := arg.Get("four"); !value.IsUndefined() {
			local.four = value.Bool()
		}
		if value := arg.Get("five"); !value.IsUndefined() {
			local.five = value.Bool()
		}
		if value := arg.Get("six"); !value.IsUndefined() {
			local.six = value.Bool()
		}
		if value := arg.Get("minus"); !value.IsUndefined() {
			local.minus = value.Bool()
		}
		// row 4
		if value := arg.Get("one"); !value.IsUndefined() {
			local.one = value.Bool()
		}
		if value := arg.Get("two"); !value.IsUndefined() {
			local.two = value.Bool()
		}
		if value := arg.Get("three"); !value.IsUndefined() {
			local.three = value.Bool()
		}
		if value := arg.Get("plus"); !value.IsUndefined() {
			local.plus = value.Bool()
		}
		// row 5
		if value := arg.Get("negate"); !value.IsUndefined() {
			local.negate = value.Bool()
		}
		if value := arg.Get("zero"); !value.IsUndefined() {
			local.zero = value.Bool()
		}
		if value := arg.Get("decimal"); !value.IsUndefined() {
			local.decimal = value.Bool()
		}
		if value := arg.Get("equals"); !value.IsUndefined() {
			local.equals = value.Bool()
		}
		localInputLock.Lock()
		localInput = &local
		localInputLock.Unlock()
		notify(local.Translate())
		return nil
	}))

	// Server → Client (self)
	go func() {
		var lastStateData []byte
		var pref float32
		var available bool = true
		for available {
			func() {
				defer func() {
					if r := recover(); r != nil {
						err := errors.New("Tried to read closed connection")
						fmt.Printf("Recovered. %s: %v\n", err.Error(), r)
						available = false
						cancel()
					}
				}()
				t, msg, err := logic.conn.Read(ctx)
				if err != nil {
					log.Printf("error read: %v", err)
					cancel()
					return
				}
				switch t {
				case websocket.MessageBinary:
					// continue
				default:
					log.Printf("error unrecognised message type: %v", t)
					cancel()
					return
				}
				prefStart := time.Now()
				data, err := encoding.DecodeWithCompression(msg)
				if err != nil {
					log.Printf("error message decode: %v", t)
					cancel()
					return
				}
				var global schema.PersonalizedState
				err = encoding.Engrain(data.(map[string]any), &global)
				if err != nil {
					log.Printf("error message engrain: %v", t)
					cancel()
					return
				}
				state := types.LocalState{
					State:   global.State,
					Id:      global.Id,
					CpuLoad: pref,
				}
				localInputLock.Lock()
				state = localInput.Meged(state)
				localInputLock.Unlock()
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
	if l.conn == nil {
		return
	}
	l.conn.Close(websocket.StatusGoingAway, "normal closure")
	l.conn = nil
}
