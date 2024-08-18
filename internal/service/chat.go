package service

import (
	"github.com/asb1302/innopolis_go_chat/internal/repository/chatdb"
	"github.com/asb1302/innopolis_go_chat/internal/service/pools"
	"github.com/asb1302/innopolis_go_chat/pkg/chatdata"
	"github.com/google/uuid"
	"time"
)

var chats chatdb.DB

func Init(chatDB chatdb.DB) {
	chats = chatDB
}

func NewMessage(msgReq chatdata.MessageChatRequest, fromID chatdata.ID) error {

	msg := chatdata.Message{
		MsgID:  chatdata.ID(uuid.New().String()),
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

	delivery := chatdata.Delivery{
		Type: chatdata.DeliveryTypeNewMsg,
		Data: chatdata.MessageChatDelivery{
			Message: msg,
			Type:    msgReq.Type,
			ChID:    msgReq.ChID,
		},
	}

	for _, userID := range users {
		if userID != fromID {
			pools.Users.Send(userID, delivery)
		}
	}
	return nil
}

func NewChat(uids []chatdata.ID) chatdata.ID {
	return chats.AddChat(uids)
}
