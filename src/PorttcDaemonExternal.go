package porttcOSSocket

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
)

func (s *PorttcHTTPServer) startExternalListener(port int) error {
	portString := strconv.Itoa(port)

	if port > 65535 || port < 0 {
		return fmt.Errorf("PorttcDaemonExternal:startExternalListener() port requested was invalid")
	}

	socketBinding := "0.0.0.0:" + portString
	ln, err := net.Listen("tcp", socketBinding)
	if err != nil {
		return fmt.Errorf("PorttcDaemonExternal:startExternalListener() could not start daemon listener %w", err)
	}

	s.Listener = ln
	s.Server = &http.Server{Handler: s.Mux}
	s.Ready = true
	return nil

}
