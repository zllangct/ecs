package client

import (
	"context"
	"fmt"
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

	for{
		time.Sleep(time.Second * 5)
		for i, conn := range c {
			conn.Write(fmt.Sprintf("chat:hi, i am %d", i))
			conn.Write(fmt.Sprintf("move:0,0,1:10"))
		}
	}
}
