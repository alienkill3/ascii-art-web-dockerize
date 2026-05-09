package tests

import (
	"ascii-art-web-dockerize/src"
	"net"
	"net/http"
	"testing"
	"time"
)

// TestServerReadTimeout verifies that the server closes connections
// that take too long to send headers/body.
func TestServerReadTimeout(t *testing.T) {
	// Setup a test server using the app's handler
	// We use a very short timeout (100ms) to make the test fast
	server := &http.Server{
		Addr:        "127.0.0.1:0", // 0 lets the OS pick a random free port
		Handler:     src.NewMux(),
		ReadTimeout: 100 * time.Millisecond,
	}

	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// Start server in a background goroutine
	go server.Serve(listener)
	defer server.Close()

	// Connect to the server as a client
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// 1. Send partial headers to keep the connection active
	_, err = conn.Write([]byte("POST /ascii-art HTTP/1.1\r\nHost: localhost\r\n"))
	if err != nil {
		t.Fatal(err)
	}

	// 2. Sleep longer than the ReadTimeout (200ms > 100ms)
	time.Sleep(200 * time.Millisecond)

	// 3. Try to read from the connection.
	// If the timeout worked, the server should have closed the connection,
	// resulting in an EOF or error immediately.
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)

	// We expect an error (like io.EOF) or 0 bytes read indicating closure.
	// If we successfully read data, the server didn't time out.
	if err == nil && n > 0 {
		t.Errorf("Expected server to close connection due to timeout, but got response: %q", string(buf[:n]))
	}
}
