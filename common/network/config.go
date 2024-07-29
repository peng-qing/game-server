package network

import (
	"encoding/binary"
	"time"
)

type Config struct {
	ConnectionID      string
	KeepaliveInterval int
	Version           int
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ByteOrder         binary.ByteOrder
}
