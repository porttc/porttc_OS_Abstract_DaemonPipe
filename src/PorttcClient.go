package porttcOSSocket


import (
	"net"
	"net/http"
)

type PorttcHTTPClient struct {
    HTTPClient http.Client
    conn       net.Conn
    platform   string
    ready      bool
}
