package chatdb

import "github.com/asb1302/innopolis_go_chat/pkg/chatdata"

type DB interface {
	AddMessage(chatID chatdata.ID, message chatdata.Message) error
	DeleteMessage(chatID chatdata.ID, messageID chatdata.ID) error
	UpdateMessage(chatID chatdata.ID, message chatdata.Message) error
	GetChatUsers(chatID chatdata.ID) ([]chatdata.ID, error)
	AddChat(uids []chatdata.ID) chatdata.ID
}
