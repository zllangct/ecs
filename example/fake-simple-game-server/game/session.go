package game

import (
	"github.com/zllangct/ecs"
	"test_ecs_fake_server/network"
)

type Session struct {
	SessionID int
	Conn      *network.TcpConn
	Entity    ecs.Entity
}
