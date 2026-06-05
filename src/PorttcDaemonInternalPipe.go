package porttcOSSocket

import (
	"fmt"
	"runtime"
	"strings"
)

// startInternalPipeListener accepts a Unix-style path (e.g. /tmp/myapp.sock)
// and converts it to a Windows named pipe path if needed.
func (s *PorttcHTTPServer) startInternalPipeListener(pathName string) error {
	if runtime.GOOS == "windows" {
		windowsPipePath := toWindowsPipePath(pathName)
		return s.startPipeListener(windowsPipePath)
	}
	return s.startUnixListener(pathName, 0600)
}

// toWindowsPipePath converts a Unix-style socket path to a Windows named pipe path.
// e.g. "/tmp/myapp.sock" -> "\\.\pipe\myapp.sock"
func toWindowsPipePath(unixPath string) string {
	// Strip leading slashes and replace remaining separators
	trimmed := strings.TrimLeft(unixPath, "/")
	normalized := strings.ReplaceAll(trimmed, "/", "-")
	return fmt.Sprintf(`\\.\pipe\%s`, normalized)
}
