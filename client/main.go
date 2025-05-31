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
	ctx, cancel := context.WithCancel(context.Background())

	// Get sharedArray
	sharedArray := js.Global().Get("sharedArray")

	// Run game tick
	l := logic.New(ctx, cancel, port, sharedArray)

	// Wait forever
	<-ctx.Done()
	l.Close()
	cancel()
}
