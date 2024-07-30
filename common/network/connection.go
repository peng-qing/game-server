package network

//func AcceptTcpConn(cfg *ConnectionConfig) (TcpConnFactory, error) {
//	addr, err := net.ResolveTCPAddr(cfg.IPVersion, fmt.Sprintf("%s:%d", cfg.IP, cfg.Port))
//	if err != nil {
//		return nil, err
//	}
//	listener, err := net.ListenTCP(cfg.IPVersion, addr)
//	if err != nil {
//		gslog.Error("[AcceptTcpConn] listen tcp failed", "ipVersion", cfg.IPVersion, "ip", cfg.IP, "port", cfg.Port)
//		return nil, err
//	}
//
//	return func(ctx context.Context, hook ConnectionHook) *net.TCPConn {
//		conn, err := listener.AcceptTCP()
//		if err != nil {
//			gslog.Critical("[TcpConnFactory] accept tcp conn failed", "addr", listener.Addr().String(), "ipVersion", cfg.IPVersion, "err", err)
//			return nil
//		}
//		if hook != nil {
//			err = hook(conn)
//			if err != nil {
//				return nil
//			}
//		}
//		return conn
//	}, nil
//}
//
//func ConnectTcpConn(cfg *ConnectionConfig, hook ConnectionHook) (TcpConnFactory, error) {
//	addr, err := net.ResolveTCPAddr(cfg.IPVersion, fmt.Sprintf("%s:%d", cfg.IP, cfg.Port))
//	if err != nil {
//		return nil, err
//	}
//	return func(ctx context.Context, hook ConnectionHook) *net.TCPConn {
//		conn, err := net.DialTCP(cfg.IPVersion, nil, addr)
//		if err != nil {
//			gslog.Critical("[TcpConnFactory] connect tcp failed", "addr", addr.String(), "ipVersion", cfg.IPVersion, "err", err)
//			return nil
//		}
//		if hook != nil {
//			err = hook(conn)
//			if err != nil {
//				return nil
//			}
//		}
//		return conn
//	}, nil
//}
