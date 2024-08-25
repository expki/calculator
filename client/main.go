package main

import (
	"calculator/src/logic"
	"syscall/js"

	"github.com/coder/websocket"
)

var port string

type Websocket struct {
	*websocket.Conn
}

func main() {
	// Get sharedArray
	sharedArray := js.Global().Get("sharedArray")

	// Run game tick
	l := logic.New(port, sharedArray)
	defer l.Close()
	// Wait forever
	select {}
}
