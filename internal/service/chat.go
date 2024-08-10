package service

import (
	"chat/internal/domain"
	"chat/internal/repository/chatdb"
	"chat/internal/service/pools"
	"github.com/google/uuid"
	"time"
)

var chats chatdb.DB

func Init(chatDB chatdb.DB) {
	chats = chatDB
}

func NewMessage(msgReq domain.MessageChatRequest, fromID domain.ID) error {

	msg := domain.Message{
		MsgID:  domain.ID(uuid.New().String()),
		Body:   msgReq.Msg,
		TDate:  time.Now(),
		FromID: fromID,
	}

	if err := chats.AddMessage(msgReq.ChID, msg); err != nil {
		return err
	}

	users, err := chats.GetChatUsers(msgReq.ChID)
	if err != nil {
		return err
	}
	toDelivery := domain.MessageChatDelivery{
		Message: msg,
		Type:    msgReq.Type,
		ChID:    msgReq.ChID,
	}
	for _, userID := range users {
		if userID != fromID {
			pools.Users.Send(userID, toDelivery)
		}
	}
	return nil
}

func NewChat(uids []domain.ID) domain.ID {
	return chats.AddChat(uids)
}
