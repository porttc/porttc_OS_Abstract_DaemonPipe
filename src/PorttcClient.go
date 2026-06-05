package porttcOSSocket


import (
    "context"
    "fmt"
    "net"
    "net/http"
    "runtime"
    "time"
)

type PorttcHTTPClient struct {
    HTTPClient http.Client
    conn       net.Conn
    platform   string
    ready      bool
}
