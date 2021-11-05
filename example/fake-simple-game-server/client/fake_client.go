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
			c[idx].Write(fmt.Sprintf("move:0,0,1:10"))
		}
	}()

	// simulation to accept pkg
	for i, conn := range c {
		go func(idx int) {
			for {
				pkg := conn.Read()
				ecs.Log.Infof("client[%d] recv: %+v", idx, pkg)
			}
		}(i)
	}

	<-ctx.Done()
}
