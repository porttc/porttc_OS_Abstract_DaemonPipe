package porttcOSSocket

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"strconv"
)

type PorttcHTTPClientConfig struct {
	Protocol     string // "Pipe", "Unix", "Loopback", "External", "InternalPipe"
	PipeName     string
	UnixSockName string
	LoopbackPort int
	ExternalHost string
	ExternalPort int
}

type PorttcHTTPClient struct {
	HTTPClient http.Client
	baseURL    string
	ready      bool
}

func (c *PorttcHTTPClient) Connect(config PorttcHTTPClientConfig) error {
	switch config.Protocol {
	case "Pipe":
		return c.connectPipe(config.PipeName)
	case "Unix":
		return c.connectUnix(config.UnixSockName)
	case "Loopback":
		return c.connectTCP("127.0.0.1", config.LoopbackPort)
	case "External":
		host := config.ExternalHost
		if host == "" {
			host = "127.0.0.1"
		}
		return c.connectTCP(host, config.ExternalPort)
	case "InternalPipe":
		if runtime.GOOS == "windows" {
			return c.connectPipe(config.PipeName)
		}
		return c.connectUnix(config.UnixSockName)
	default:
		return fmt.Errorf("PorttcHTTPClient.Connect() invalid protocol %q", config.Protocol)
	}
}

func (c *PorttcHTTPClient) connectUnix(socketPath string) error {
	if socketPath == "" {
		return fmt.Errorf("PorttcHTTPClient.connectUnix() socket path cannot be empty")
	}
	c.HTTPClient = http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
	c.baseURL = "http://localhost"
	c.ready = true
	return nil
}

func (c *PorttcHTTPClient) connectTCP(host string, port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("PorttcHTTPClient.connectTCP() invalid port %d", port)
	}
	c.HTTPClient = http.Client{}
	c.baseURL = "http://" + host + ":" + strconv.Itoa(port)
	c.ready = true
	return nil
}

func (c *PorttcHTTPClient) Get(path string) (*http.Response, error) {
	if !c.ready {
		return nil, fmt.Errorf("PorttcHTTPClient.Get() client not connected")
	}
	return c.HTTPClient.Get(c.baseURL + path)
}

func (c *PorttcHTTPClient) Post(path, contentType string, body io.Reader) (*http.Response, error) {
	if !c.ready {
		return nil, fmt.Errorf("PorttcHTTPClient.Post() client not connected")
	}
	return c.HTTPClient.Post(c.baseURL+path, contentType, body)
}

// NewRequest builds a request with the base URL prepended, for use with Do().
func (c *PorttcHTTPClient) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, c.baseURL+path, body)
}

func (c *PorttcHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if !c.ready {
		return nil, fmt.Errorf("PorttcHTTPClient.Do() client not connected")
	}
	return c.HTTPClient.Do(req)
}

func (c *PorttcHTTPClient) Close() {
	c.HTTPClient.CloseIdleConnections()
	c.ready = false
}
