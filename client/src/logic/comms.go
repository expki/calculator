package logic

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"

	"github.com/coder/websocket"
)

func (c *Logic) connect(ctx context.Context) (err error) {
	c.conn, _, err = websocket.Dial(ctx, c.url, nil)
	return err
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
