package comms

import (
	"fmt"
	"net/url"
	"strconv"
	"syscall/js"
)

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
