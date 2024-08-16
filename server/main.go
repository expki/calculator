package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/expki/calculator/lib/schema"
	"github.com/gorilla/websocket"
)

func main() {

	var state schema.State
	fmt.Println(state)

	http.HandleFunc("/", echoHandler) // Set the route handler

	// Start the server on localhost port 8080 and log any errors
	log.Println("WebSocket server starting on :8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		// Read message from browser
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break // Exit the loop on error
		}

		// Print the message to the console
		fmt.Printf("Received: %s\n", message)

		// Write message back to browser
		if err := conn.WriteMessage(mt, message); err != nil {
			log.Println("Write error:", err)
			break // Exit the loop on error
		}
	}
}
