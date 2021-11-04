package game

import (
	"context"
	"github.com/zllangct/ecs"
	"reflect"
	"test_ecs/network"
)
type FakeGame struct {
	clients map[int]*Session
	world *ecs.World
}

func NewGame() *FakeGame {
	return &FakeGame{
		clients: map[int]*Session{},
	}
}

func (f *FakeGame) Run(ctx context.Context) {
	f.InitEcs()
	f.InitNetwork()
}

func (f *FakeGame) InitEcs() {
	rt := ecs.Runtime
	rt.Run()

	f.world = rt.NewWorld(ecs.NewDefaultWorldConfig())
	f.world.Run()

	f.world.Register(&MoveSystem{})
	f.world.Register(&SyncSystem{})
}

func (f *FakeGame) EnterGame(sess *Session) {
	e := f.world.NewEntity()
	e.Add(&PlayerComponent{})
	e.Add(&Position{
		X: 100,
		Y: 100,
		Z: 100,
	})
}

func (f *FakeGame) InitNetwork() {
	lis, err := network.Listen()
	if err != nil {
		return
	}

	seq := 0
	for{
		conn := lis.Accept()
		seq++
		sess := &Session{
			SessionID: seq,
			Conn: conn,
		}

		go func(conn *network.TcpConn, sess *Session) {
			for{
				pkg := conn.Read()
				f.Dispatch(pkg)
			}
		}(conn, sess)
	}
}

func (f *FakeGame) OnClientEnter(sess *Session) {
	f.clients[sess.SessionID] = sess
	f.EnterGame(sess)
}

func (f *FakeGame) Dispatch(pkg interface{}) {
	//TODO for controller system
}

func (f *FakeGame) ChangeMovementTimeScale() {
	sys, ok := f.world.GetSystem(reflect.TypeOf(&MoveSystem{}))
	if !ok {
		return
	}
	sys.Emit("UpdateTimeScale", float64(1.2))
}
