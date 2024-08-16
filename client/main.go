package main

import (
	"calculator/src/logic"
	"calculator/src/logic/comms"
	"calculator/src/userinput"
	"syscall/js"

	"github.com/coder/websocket"
)

type Websocket struct {
	*websocket.Conn
}

func main() {
	// Get sharedArray
	sharedArray := js.Global().Get("sharedArray")

	// Run game tick
	logic.LogicLoop(sharedArray)

	// Listen to user input
	input := userinput.New()
	defer input.Close()

	// Connect to server
	comms.New()

	// Wait forever
	select {}
}
