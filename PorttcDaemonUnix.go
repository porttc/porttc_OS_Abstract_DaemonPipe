package porttcOSSocket

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func (s *PorttcHTTPServer) startUnixListener(socketPath string, perms int) error {
	if socketPath == "" {
		return fmt.Errorf("PorttcDaemonUnix:startUnixListener() socket path cannot be empty")
	}

	os.Remove(socketPath)

	ln, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}

	s.Listener = ln
	s.Server = &http.Server{Handler: s.Mux}
	s.Ready = true
	return nil
}
