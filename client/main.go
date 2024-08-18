package main

import (
	"calculator/src/logic"
	"calculator/src/logic/comms"
	"calculator/src/userinput"
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
	lgc := logic.New()
	go lgc.LogicLoop(sharedArray)

	// Listen to user input
	input := userinput.New()
	defer input.Close()

	// Connect to server
	cms := comms.New(port, lgc)
	defer cms.Close()

	// Wait forever
	select {}
}
