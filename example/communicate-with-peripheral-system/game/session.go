package game

import "test_ecs/network"

type Session struct {
	SessionID int
	Conn *network.TcpConn
}
