//go:build !windows

package porttcOSSocket

import "fmt"

func (c *PorttcHTTPClient) connectPipe(pipeName string) error {
	return fmt.Errorf("windows named pipes are not supported on this platform")
}
