//go:build !windows

package porttcOSSocket

import "fmt"

func (s *PorttcHTTPServer) startPipeListener(pipePath string) error {
	return fmt.Errorf("windows named pipes are not supported on this platform")
}
