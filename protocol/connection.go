package protocol

import "net"

type WSConn struct {
	Conn net.Conn
}

func NewWebsocketConn() *WSConn {
	return &WSConn{}
}
