package network

import (
	"encoding/binary"
	"time"
)

type ConnectionLayerConfig struct {
	ConnectionID      string
	KeepaliveInterval int
	Version           int
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ByteOrder         binary.ByteOrder
}

type ConnectionConfig struct {
	IP        string
	Port      int
	IPVersion string
}
