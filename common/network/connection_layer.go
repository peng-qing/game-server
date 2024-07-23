package network

import (
	"net"
	"time"
)

type ConnectionKeeper struct {
	net.Conn
	cfg *ConnectionConfig
}

func NewConnectionLayer(conn net.Conn, cfg *ConnectionConfig) ConnectionLayer {

	return &ConnectionKeeper{Conn: conn, cfg: cfg}
}

func (gs *ConnectionKeeper) ConnectionID() string {
	//TODO implement me
	panic("implement me")
}

func (gs *ConnectionKeeper) Heartbeat() time.Duration {
	//TODO implement me
	panic("implement me")
}
