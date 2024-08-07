package cache

import (
	"play/websocket/domain"
	"sync"
)

var Dials = DialPool{
	Pool: make(map[string]*domain.UserDialog),
}

type DialPool struct {
	sync.Mutex
	// key - ChatId
	Pool map[string]*domain.UserDialog
}

func (p *DialPool) Set(chid string, dial *domain.UserDialog) {
	p.Lock()
	p.Pool[chid] = dial
	p.Unlock()
}

func (p *DialPool) Get(chid string) *domain.UserDialog {
	p.Lock()
	dial := p.Pool[chid]
	p.Unlock()
	return dial
}
