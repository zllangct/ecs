package game

import (
	"context"
	"errors"
	"github.com/zllangct/ecs"
	"strconv"
	"strings"
	"sync"
	"test_ecs_fake_server/network"
)

var send chan Msg2Client = make(chan Msg2Client, 10)

type Msg2Client struct {
	SessionID int
	Content   interface{}
}

type FakeGame struct {
	clients  sync.Map
	world    *ecs.AsyncWorld
	chatRoom *ChatRoom
}

func NewGame() *FakeGame {
	return &FakeGame{}
}

func (f *FakeGame) Run(ctx context.Context) {
	f.InitEcs()
	f.InitChat()
	f.InitNetwork()
}

func (f *FakeGame) InitEcs() {
	//create config
	config := ecs.NewDefaultWorldConfig()

	//create a world and startup
	f.world = ecs.NewAsyncWorld(config)
	f.world.Startup()

	//register your system
	ecs.RegisterSystem[MoveSystem](f.world)
	ecs.RegisterSystem[SyncSystem](f.world)
	ecs.RegisterSystem[EmptySystem](f.world)
}

func (f *FakeGame) EnterGame(sess *Session) {
	f.world.Wait(func(gaw ecs.SyncWrapper) error {
		e := gaw.NewEntity()
		gaw.Add(e, &PlayerComponent{
			SessionID: sess.SessionID,
		})
		gaw.Add(e, &Position{
			X: 100,
			Y: 100,
			Z: 100,
		})
		gaw.Add(e, &Movement{
			V:   2000,
			Dir: [3]int{1, 0, 0},
		})
		sess.Entity = e
		return nil
	})
}

func (f *FakeGame) InitNetwork() {
	lis, err := network.Listen()
	if err != nil {
		return
	}

	go func() {
		for {
			select {
			case m := <-send:
				obj, ok := f.clients.Load(m.SessionID)
				if ok {
					sess := obj.(*Session)
					sess.Conn.Write(m.Content)
				}
			}
		}
	}()

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

func SendToClient(sessionId int, content interface{}) {
	send <- Msg2Client{
		SessionID: sessionId,
		Content:   content,
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
		if len(d) != 3 {
			return
		}
		var dir [3]int
		for i := 0; i < 3; i++ {
			value, _ := strconv.Atoi(d[i])
			dir[i] = value
		}

		v, _ := strconv.Atoi(split[2])
		f.Move(sess.Entity, v, dir)
	}
}

func (f *FakeGame) Move(entity ecs.Entity, v int, dir [3]int) {
	if f.world == nil {
		return
	}
	f.world.Sync(func(gaw ecs.SyncWrapper) error {
		u, ok := ecs.GetUtility[MoveSystemUtility](gaw)
		if !ok {
			return errors.New("can not find MoveSystemUtility")
		}
		return u.Move(entity, v, dir)
	})
}

func (f *FakeGame) ChangeMovementTimeScale(timeScale float64) {
	if f.world == nil {
		return
	}
	f.world.Sync(func(gaw ecs.SyncWrapper) error {
		u, ok := ecs.GetUtility[MoveSystemUtility](gaw)
		if !ok {
			return errors.New("can not find MoveSystemUtility")
		}
		return u.UpdateTimeScale(timeScale)
	})
}

func (f *FakeGame) InitChat() {
	f.chatRoom = NewChatRoom(&f.clients)
}
