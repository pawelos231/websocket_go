package helpers

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
)

func WriteToConnection(conn net.Conn, res http.Response, size uint16) error {
	var resBuf bytes.Buffer

	if err := res.Write(&resBuf); err != nil {
		return fmt.Errorf("failed to write to connection buffer: %w", err)
	}

	if _, err := conn.Write(resBuf.Bytes()); err != nil {
		conn.Close()
		return fmt.Errorf("failed to write to connection: %w", err)
	}

	return nil
}
