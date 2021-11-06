package game

import (
	"context"
	"github.com/zllangct/ecs"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"test_ecs/network"
)

type FakeGame struct {
	clients  sync.Map
	world    *ecs.World
	chatRoom *ChatRoom
}

func NewGame() *FakeGame {
	return &FakeGame{}
}

func (f *FakeGame) Run(ctx context.Context) {
	f.InitEcs()
	f.InitNetwork()
}

func (f *FakeGame) InitEcs() {
	//run your ecs runtime
	ecs.RuntimeConfigure(ecs.NewDefaultRuntimeConfig())
	ecs.Run()

	//create a world
	f.world = ecs.CreateWorld(ecs.NewDefaultWorldConfig())
	//world.Run() or ecs.WorldRun(world)
	ecs.WorldRun(f.world)

	//register your system
	// ecs.RegisterSystem[System](world) or world.Register(system)
	ecs.RegisterSystem[MoveSystem](f.world)
	ecs.RegisterSystem[SyncSystem](f.world)
}

func (f *FakeGame) EnterGame(sess *Session) {
	e := f.world.NewEntity()
	e.Add(&PlayerComponent{})
	e.Add(&Position{
		X: 100,
		Y: 100,
		Z: 100,
	})

	sess.EntityId = e.hashKey()
}

func (f *FakeGame) InitNetwork() {
	lis, err := network.Listen()
	if err != nil {
		return
	}

	seq := 0
	for {
		conn := lis.Accept()
		seq++
		sess := &Session{
			SessionID: seq,
			Conn:      conn,
		}

		go func(conn *network.TcpConn, sess *Session) {
			f.OnClientEnter(sess)
			for {
				pkg := conn.Read()
				f.Dispatch(pkg, sess)
			}
		}(conn, sess)
	}
}

func (f *FakeGame) OnClientEnter(sess *Session) {
	f.clients.Store(sess.SessionID, sess)
	f.EnterGame(sess)
}

func (f *FakeGame) Dispatch(pkg interface{}, sess *Session) {
	content, ok := pkg.(string)
	if !ok {
		return
	}

	split := strings.Split(content, ":")
	op := split[0]

	switch op {
	case "chat":
		// not handle by ecs
		f.chatRoom.Talk(split[1])
	case "move":
		// handle by ecs
		if len(split) != 3 {
			return
		}
		d := strings.Split(split[1], ",")
		var dir []int
		for _, s := range d {
			value, _ := strconv.Atoi(s)
			dir = append(dir, value)
		}

		v, _ := strconv.Atoi(split[2])
		e := ecs.GetEntity(f.world, sess.EntityId)
		e.Add()
	}
}

func (f *FakeGame) ChangeMovementTimeScale(timeScale float64) {
	sys, ok := f.world.GetSystem(reflect.TypeOf(&MoveSystem{}))
	if !ok {
		return
	}
	sys.Emit("UpdateTimeScale", timeScale)
}

func (f *FakeGame) InitChat() {
	f.chatRoom = NewChatRoom(&f.clients)
}
