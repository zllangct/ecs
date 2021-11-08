package game

import (
	"sync"
)

type ChatRoom struct {
	clients *sync.Map
}

func NewChatRoom(c *sync.Map) *ChatRoom {
	return &ChatRoom{
		clients: &sync.Map{},
	}
}

func (c *ChatRoom) Talk(content string) {
	c.clients.Range(func(k, v interface{}) bool {
		sess := v.(*Session)
		sess.Conn.Write(content)
		return true
	})
}
