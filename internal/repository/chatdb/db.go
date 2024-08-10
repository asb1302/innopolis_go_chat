package chatdb

import "chat/internal/domain"

type DB interface {
	AddMessage(chatID domain.ID, message domain.Message) error
	DeleteMessage(chatID domain.ID, messageID domain.ID) error
	UpdateMessage(chatID domain.ID, message domain.Message) error
	GetChatUsers(chatID domain.ID) ([]domain.ID, error)
	AddChat(uids []domain.ID) domain.ID
}
