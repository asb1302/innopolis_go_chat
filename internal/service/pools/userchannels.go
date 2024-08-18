package pools

import (
	"github.com/asb1302/innopolis_go_chat/pkg/chatdata"
	"sync"
)

var Users = userPool{
	pool: make(map[chatdata.ID]chan interface{}),
}

type userPool struct {
	sync.Mutex
	// key - user id
	pool map[chatdata.ID]chan interface{}
}

func (p *userPool) Send(uid chatdata.ID, msg interface{}) {
	p.Lock()
	defer p.Unlock()
	ch, ok := p.pool[uid]
	if !ok {
		return
	}
	// на подумать: буферизированный или горутину
	ch <- msg
}

func (p *userPool) New(uid chatdata.ID) <-chan interface{} {
	p.Lock()
	ch := make(chan interface{})
	p.pool[uid] = ch
	p.Unlock()

	return ch
}

func (p *userPool) Delete(uid chatdata.ID) bool {
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
