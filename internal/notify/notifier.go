package notify

import (
	"fmt"
	"sync"
)

type InChannel interface {
	In(interface{})
	Close()
}

type Notifier interface {
	Subscribe() chan interface{}
	Unsubscribe()
	Status() int
}

type Connection struct {
	mu    sync.Mutex
	users []*user
	data  chan interface{}
}

type user struct {
	data interface{}
	ch   chan interface{}
}

const (
	Buffer     = 10
	DataBuffer = 10
)

func NewNotifier() *Connection {
	return &Connection{
		data: make(chan interface{}, DataBuffer),
	}
}

func (c *Connection) Subscribe() chan interface{} {
	c.mu.Lock()
	defer c.mu.Unlock()

	ch := make(chan interface{}, Buffer)
	u := &user{ch: ch}

	c.users = append(c.users, u)
	return ch
}

func (c *Connection) Unsubscribe() {
	// TODO(Trung): Implement logic
}

func (c *Connection) Status() int {
	fmt.Printf("Num of users: %d, num of data: %d\n", len(c.users), len(c.data))
	return len(c.users)
}

func (c *Connection) In(packet interface{}) {
	c.data <- packet
}

func (c *Connection) Close() {
	// TODO(Trung): Implement logic
}

func (c *Connection) Notify() {
	users := make(chan user, 10)

	for i := 0; i < 20; i++ {
		go notify(users)
	}

	for d := range c.data {
		c.mu.Lock()
		for _, u := range c.users {
			u.data = d
			users <- *u
		}
		c.mu.Unlock()
	}
}

func notify(user chan user) {
	for u := range user {
		u.ch <- u.data
	}
}
