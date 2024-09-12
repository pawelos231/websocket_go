package protocol

import (
	"bufio"
	"net"
)

type WSConn struct {
	conn        net.Conn
	reader      *bufio.Reader
	writeBuffer []byte
	mutex       chan struct{} // protect write to the connection
}

func CreateConn(conn net.Conn, readBufferSize uint16, writeBufferSize uint16, reader *bufio.Reader, writeBuf []byte) *WSConn {
	return &WSConn{}
}

func (c *Conn) NetConn() net.Conn {
	return c.conn
}
