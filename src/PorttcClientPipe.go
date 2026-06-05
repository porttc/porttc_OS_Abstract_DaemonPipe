//go:build windows

package porttcOSSocket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	winio "github.com/Microsoft/go-winio"
)

func (c *PorttcHTTPClient) connectPipe(pipeName string) error {
	if pipeName == "" {
		return fmt.Errorf("PorttcHTTPClient.connectPipe() pipe name cannot be empty")
	}
	pipePath := `\\.\pipe\` + pipeName
	c.HTTPClient = http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				timeout := 30 * time.Second
				return winio.DialPipe(pipePath, &timeout)
			},
		},
	}
	c.baseURL = "http://localhost"
	c.ready = true
	return nil
}
