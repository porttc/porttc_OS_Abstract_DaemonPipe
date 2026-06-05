package porttcOSSocket

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
)

func (s *PorttcHTTPServer) startLoopbackListener(port int) error {
	portString := strconv.Itoa(port)

	if port > 65535 || port < 0 {
		return fmt.Errorf("PorttcDaemonLoopback:startLoopbackListener() port requested was invalid")
	}

	socketBinding := "127.0.0.1:" + portString
	ln, err := net.Listen("tcp", socketBinding)
	if err != nil {
		return fmt.Errorf("PorttcDaemonLoopback:startloopbackListener() could not start daemon listener %w", err)
	}

	s.Listener = ln
	s.Server = &http.Server{Handler: s.Mux}
	s.Ready = true
	return nil
}
