package porttcOSSocket

import (
	"fmt"
	"net/http"

	winio "github.com/microsoft/go-winio"
)

func (s *PorttcHTTPServer) startPipeListener(pipePath string) error {
	if pipePath == "" {
		return fmt.Errorf("PorttcDaemonPipe:startPipeListener() got invalid path")
	}

	ln, err := winio.ListenPipe(`\\.\pipe\`+pipePath, nil)
	if err != nil {
		return err
	}

	s.Listener = ln
	s.Server = &http.Server{Handler: s.Mux}
	s.Ready = true
	return nil
}
