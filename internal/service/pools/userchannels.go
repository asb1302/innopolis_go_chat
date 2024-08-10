package pools

import (
	"chat/internal/domain"
	"sync"
)

var Users = userPool{
	pool: make(map[domain.ID]chan interface{}),
}

type userPool struct {
	sync.Mutex
	// key - user id
	pool map[domain.ID]chan interface{}
}

func (p *userPool) Send(uid domain.ID, msg interface{}) {
	p.Lock()
	defer p.Unlock()
	ch, ok := p.pool[uid]
	if !ok {
		return
	}
	// на подумать: буферизированный или горутину
	ch <- msg
}

func (p *userPool) New(uid domain.ID) <-chan interface{} {
	p.Lock()
	ch := make(chan interface{})
	p.pool[uid] = ch
	p.Unlock()

	return ch
}

func (p *userPool) Delete(uid domain.ID) bool {
	p.Lock()
	defer p.Unlock()
	ch, ok := p.pool[uid]
	if !ok {
		return ok
	}
	delete(p.pool, uid)
	close(ch)
	return ok
}
