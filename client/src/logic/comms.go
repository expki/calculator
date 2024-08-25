package logic

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"syscall/js"
	"time"

	"github.com/coder/websocket"
)

func (c *Logic) connect() {
	for {
		if success := func() bool {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var err error
			c.conn, _, err = websocket.Dial(ctx, c.url, nil)
			if err != nil {
				log.Printf("websocket.Dial exception: %v", err)
				return false
			}
			return true
		}(); success {
			return
		}
		log.Println("failed to connect, retrying every 3 seconds")
	}
}

func getWebsocketURL(port string) (string, error) {
	uri, err := url.Parse(js.Global().Get("location").Get("href").String())
	if err != nil {
		return "", fmt.Errorf("url.Parse window href: %v", err)
	}
	if port != "" {
		wsPort, err := strconv.Atoi(port)
		if err != nil {
			return "", fmt.Errorf("strconv.Atoi alt port: %v", err)
		}
		uri.Host = fmt.Sprintf("%s:%d", uri.Hostname(), wsPort)
	}
	switch uri.Scheme {
	case "http":
		uri.Scheme = "ws"
	case "https":
		uri.Scheme = "wss"
	default:
		return "", fmt.Errorf("unsupported scheme: %q", uri.Scheme)
	}
	uri.Path = "session"
	return uri.String(), nil
}
