package porttcOSSocket_test

import (
	"io"
	"net"
	"net/http"
	"path/filepath"
	"runtime"
	"testing"

	porttc "github.com/porttc/porttc_OS_Abstract_DaemonPipe"
)

// startServer starts a server with the given config, registers a /ping handler,
// and schedules cleanup via t.Cleanup.
func startServer(t *testing.T, config porttc.PorttcHTTPServerConfig) *porttc.PorttcHTTPServer {
	t.Helper()
	var srv porttc.PorttcHTTPServer
	if err := srv.Start(config); err != nil {
		t.Fatalf("Start: %v", err)
	}
	srv.Mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	go srv.Server.Serve(srv.Listener)
	t.Cleanup(func() { srv.Stop() })
	return &srv
}

func connectClient(t *testing.T, config porttc.PorttcHTTPClientConfig) *porttc.PorttcHTTPClient {
	t.Helper()
	var c porttc.PorttcHTTPClient
	if err := c.Connect(config); err != nil {
		t.Fatalf("Connect: %v", err)
	}
	t.Cleanup(func() { c.Close() })
	return &c
}

func assertPing(t *testing.T, c *porttc.PorttcHTTPClient) {
	t.Helper()
	resp, err := c.Get("/ping")
	if err != nil {
		t.Fatalf("GET /ping: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "pong" {
		t.Fatalf("expected body %q, got %q", "pong", string(body))
	}
}

func boundTCPPort(t *testing.T, ln net.Listener) int {
	t.Helper()
	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("listener is not TCP")
	}
	return addr.Port
}

func TestLoopback(t *testing.T) {
	srv := startServer(t, porttc.PorttcHTTPServerConfig{
		Protocol:     "Loopback",
		LoopbackPort: 0,
	})
	port := boundTCPPort(t, srv.Listener)

	c := connectClient(t, porttc.PorttcHTTPClientConfig{
		Protocol:     "Loopback",
		LoopbackPort: port,
	})
	assertPing(t, c)
}

func TestExternal(t *testing.T) {
	srv := startServer(t, porttc.PorttcHTTPServerConfig{
		Protocol:     "External",
		ExternalPort: 0,
	})
	port := boundTCPPort(t, srv.Listener)

	c := connectClient(t, porttc.PorttcHTTPClientConfig{
		Protocol:     "External",
		ExternalHost: "127.0.0.1",
		ExternalPort: port,
	})
	assertPing(t, c)
}

func TestUnix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix domain sockets not supported on Windows")
	}
	sockPath := filepath.Join(t.TempDir(), "test.sock")

	startServer(t, porttc.PorttcHTTPServerConfig{
		Protocol:     "Unix",
		UnixSockName: sockPath,
	})
	c := connectClient(t, porttc.PorttcHTTPClientConfig{
		Protocol:     "Unix",
		UnixSockName: sockPath,
	})
	assertPing(t, c)
}

func TestPipe(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows named pipes only available on Windows")
	}
	const pipeName = "porttc-test-pipe"

	startServer(t, porttc.PorttcHTTPServerConfig{
		Protocol: "Pipe",
		PipeName: pipeName,
	})
	c := connectClient(t, porttc.PorttcHTTPClientConfig{
		Protocol: "Pipe",
		PipeName: pipeName,
	})
	assertPing(t, c)
}

func TestInternalPipe(t *testing.T) {
	srvCfg := porttc.PorttcHTTPServerConfig{Protocol: "InternalPipe"}
	cliCfg := porttc.PorttcHTTPClientConfig{Protocol: "InternalPipe"}

	if runtime.GOOS == "windows" {
		srvCfg.PipeName = "porttc-test-internal"
		cliCfg.PipeName = "porttc-test-internal"
	} else {
		sockPath := filepath.Join(t.TempDir(), "internal.sock")
		srvCfg.UnixSockName = sockPath
		cliCfg.UnixSockName = sockPath
	}

	startServer(t, srvCfg)
	c := connectClient(t, cliCfg)
	assertPing(t, c)
}
