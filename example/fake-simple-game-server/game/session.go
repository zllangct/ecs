package game

import (
	"github.com/zllangct/ecs"
	"test_ecs/network"
)

type Session struct {
	SessionID int
	Conn      *network.TcpConn
	Entity  ecs.Entity
}
