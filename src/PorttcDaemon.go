package porttcOSSocket

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

type PorttcHTTPServer struct {
	Server   *http.Server
	Listener net.Listener
	Mux      *http.ServeMux
	Platform string
	Ready    bool
}

type PorttcHTTPServerConfig struct {
	//Pipe (windows)
	//Unix socket (unix)
	//Loopback (all OS)
	//External (all OS)
	//InternalPipe (all OS just uses pipe or unix depending on OS)
	Protocol string

	//external listener info
	ExternalPort int

	//loopback adapter info
	LoopbackPort int

	//windows pipe info
	PipeName string

	//unix socket info
	UnixSockName string
}

func (s *PorttcHTTPServer) Start(config PorttcHTTPServerConfig) error {
	s.Mux = http.NewServeMux()

	switch config.Protocol {

	//windows
	case "Pipe":
		return s.startPipeListener(config.PipeName)

	//unix + win 10 supported by unix sockets
	case "Unix":
		return s.startUnixListener(config.UnixSockName, 0600)

	//platform independant options
	case "Loopback":
		return s.startLoopbackListener(config.LoopbackPort)
	case "External":
		return s.startExternalListener(config.ExternalPort)
	case "InternalPipe":
		var handleName string
		if config.PipeName != "" {
			handleName = config.PipeName
		} else {
			handleName = config.UnixSockName
		}
		return s.startInternalPipeListener(handleName)

	default:
		return fmt.Errorf("PorttcDaemon.go Start() invalid server listener type")
	}
}

func (s *PorttcHTTPServer) Stop() error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	s.Ready = false
	err := s.Server.Shutdown(ctx)
	if err != nil {
		s.Server.Close()
		return fmt.Errorf("porttcDaemon.stop timeout exceeded forcing shutdown")
	}
	return nil
}
