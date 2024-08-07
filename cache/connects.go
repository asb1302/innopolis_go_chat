package cache

import (
	"github.com/gorilla/websocket"

	"sync"
)

var Conns = ConnPool{
	Pool: make(map[string]*websocket.Conn),
}

type ConnPool struct {
	sync.Mutex
	// key - user id
	Pool map[string]*websocket.Conn
}

func (p *ConnPool) SetConn(uid string, conn *websocket.Conn) {
	p.Lock()
	p.Pool[uid] = conn
	p.Unlock()
}

func (p *ConnPool) DeleteConn(uid string) {
	p.Lock()
	p.Pool[uid] = nil
	p.Unlock()
}

func (p *ConnPool) GetConn(uid string) *websocket.Conn {
	p.Lock()
	conn := p.Pool[uid]
	p.Unlock()
	return conn
}
