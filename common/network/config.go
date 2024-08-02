package network

import (
	"encoding/binary"
	"errors"
	"net"
	"time"
)

type BrokerConf struct {
	ConnectionID      string
	KeepaliveInterval int
	Version           int
	WriteTimeout      time.Duration
	ReadTimeout       time.Duration
	ByteOrder         binary.ByteOrder
	OnCloseCallback   OnConnectionCloseCallback
}

func IsNetTimeout(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}

//
//type ConnectionLayerConfig struct {

//}
//
//type ConnectionConfig struct {
//	IP        string
//	Port      int
//	IPVersion string
//}

//func AcceptHook() {
//		// read connect packet
//		if gs.readTimeout > 0 {
//			_ = conn.SetReadDeadline(time.Now().Add(gs.readTimeout))
//		}
//		msg, err := ReadPacket(conn, gs.byteOrder)
//		if err != nil {
//			return err
//		}
//		if gs.readTimeout > 0 {
//			_ = conn.SetReadDeadline(time.Time{})
//		}
//		packet, ok := msg.(*ConnectPacket)
//		if !ok || packet.Validate() != Accepted {
//			return err
//		}
//		if gs.ConnectionID() != packet.ClientIdentifier {
//			return ErrNotSameConnectionID
//		}
//		gs.lock.Lock()
//		gs.version = packet.ProtocolVersion
//		gs.keepalive = time.Duration(packet.Keepalive)
//		gs.lock.Unlock()
//
//		// send connection ack packet
//		ackPacket := NewControlPacket(ConnectAck).(*ConnectAckPacket)
//		if gs.writeTimeout > 0 {
//			_ = conn.SetWriteDeadline(time.Now().Add(gs.writeTimeout))
//		}
//		_, err = ackPacket.WriteTo(conn, gs.byteOrder)
//		if err != nil {
//			return err
//		}
//		if gs.writeTimeout > 0 {
//			_ = conn.SetWriteDeadline(time.Time{})
//		}
//
//		return nil
//}
//
//func ConnectHook() {
//
//}
