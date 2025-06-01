package main

import (
	"calculator/src/logic"
	"context"
	"syscall/js"

	"github.com/coder/websocket"
)

var port string

type Websocket struct {
	*websocket.Conn
}

func main() {
	appCtx, stop := context.WithCancel(context.Background())

	// stop on close
	js.Global().Set("onbeforeunload", js.FuncOf(func(this js.Value, args []js.Value) any {
		stop()
		return nil
	}))

	// Get sharedArray
	sharedArray := js.Global().Get("sharedArray")

	// Reconnect on disconnect
	for {
		ctx, cancel := context.WithCancel(appCtx)

		// Start game tick
		l := logic.New(ctx, cancel, port, sharedArray)

		// Wait for close
		select {
		case <-appCtx.Done(): // app close
			l.Close()
			return
		case <-ctx.Done(): // connection close
			l.Close()
			if appCtx.Err() != nil {
				return
			}
		}
	}
}
