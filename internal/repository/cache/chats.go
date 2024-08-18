package cache

import (
	"context"
	"errors"
	"github.com/asb1302/innopolis_go_chat/pkg/chatdata"
	"github.com/google/uuid"
	"sync"
)

const chatDumpFileName = "chats.json"

type ChatsPool struct {
	sync.Mutex
	// key - ChatId
	pool map[chatdata.ID]*chatdata.Chat
}

func ChatCacheInit(ctx context.Context, wg *sync.WaitGroup) (*ChatsPool, error) {
	var chats = ChatsPool{
		pool: make(map[chatdata.ID]*chatdata.Chat),
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		makeDump(chatDumpFileName, chats.pool)
	}()

	if err := loadFromDump(chatDumpFileName, &chats.pool); err != nil {
		return nil, err
	}

	return &chats, nil
}

func (p *ChatsPool) AddChat(uids []chatdata.ID) chatdata.ID {
	chid := chatdata.ID(uuid.New().String())
	nc := chatdata.Chat{
		UIDs: uids,
		ChID: chid,
	}

	p.Lock()
	p.pool[chid] = &nc
	p.Unlock()

	return chid
}

func (p *ChatsPool) AddMessage(chatID chatdata.ID, message chatdata.Message) error {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return errors.New("chat not found")
	}
	chat.Messages = append(chat.Messages, message)
	p.pool[chatID] = chat
	return nil
}

func (p *ChatsPool) DeleteMessage(chatID chatdata.ID, messageID chatdata.ID) error {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return errors.New("chat not found")
	}

	// TODO think
	for i, message := range chat.Messages {
		if message.MsgID == messageID {
			chat.Messages = append(chat.Messages[:i], chat.Messages[i+1:]...)
			break
		}
	}
	p.pool[chatID] = chat
	return nil
}

func (p *ChatsPool) UpdateMessage(chatID chatdata.ID, message chatdata.Message) error {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return errors.New("chat not found")
	}

	// TODO think
	for i, m := range chat.Messages {
		if m.MsgID == message.MsgID {
			// TODO надо проверить, тот же автор ли?
			chat.Messages[i] = message
			break
		}
	}
	p.pool[chatID] = chat
	return nil
}

func (p *ChatsPool) GetChatUsers(chatID chatdata.ID) ([]chatdata.ID, error) {
	p.Lock()
	defer p.Unlock()

	chat, ok := p.pool[chatID]
	if !ok {
		return nil, errors.New("chat not found")
	}
	return chat.UIDs, nil
}
