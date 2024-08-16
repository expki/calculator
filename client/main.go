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
	// Listen to user input
	input := userinput.New()
	defer input.Close()

	// Get sharedArray
	sharedArray := js.Global().Get("sharedArray")

	// Connect to server
	comms.New()

	// Run game tick
	logic.LogicLoop(sharedArray)

	// Wait forever
	select {}
}
