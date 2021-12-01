package client

import (
	"context"
	"fmt"
	"github.com/zllangct/ecs"
	"math/rand"
	"test_ecs/network"
	"time"
)

type FakeClient struct {
}

func NewClient() *FakeClient {
	return &FakeClient{}
}

func (f *FakeClient) Run(ctx context.Context) {
	var c []*network.TcpConn
	for i := 0; i < 5; i++ {
		conn := network.Dial("127.0.0.1:3333")
		c = append(c, conn)
	}

	// simulation to send pkg
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			idx := rand.Intn(len(c))
			c[idx].Write(fmt.Sprintf("chat:hi, i am %d", idx))
		}
	}()

	// simulation to control player
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(rand.Intn(5)))
			idx := rand.Intn(len(c))
			v := rand.Intn(1000)
			dir := [3]int{0, 0, 0}
			dir[rand.Intn(3)] = 1
			c[idx].Write(fmt.Sprintf("move:%d,%d,%d:%d", dir[0], dir[1], dir[2], v))
		}
	}()

	// simulation to accept pkg
	for i, conn := range c {
		go func(idx int, conn *network.TcpConn) {
			for {
				pkg := conn.Read()
				ecs.Log.Infof("client[%d] recv: %+v", idx, pkg)
			}
		}(i, conn)
	}

	<-ctx.Done()
}
