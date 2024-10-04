package repository

import (
	"github.com/mayye4ka/pinder/models"
)

type Chat struct {
	ID    uint64
	User1 uint64
	User2 uint64
}

func (Chat) TableName() string {
	return "chats"
}

func (r *Repository) CreateChat(user1, user2 uint64) error {
	chat := Chat{
		User1: user1,
		User2: user2,
	}
	res := r.db.Create(&chat)
	return res.Error
}

func (r *Repository) GetChats(userID uint64) ([]models.Chat, error) {
	var chats []Chat
	res := r.db.Model(&Chat{}).Where("user_1 = ? or user_2 = ?", userID, userID).Find(&chats)
	if res.Error != nil {
		return nil, res.Error
	}
	return mapChats(chats), nil
}

func (r *Repository) GetChat(id uint64) (models.Chat, error) {
	var chat Chat
	res := r.db.Model(&Chat{}).Where("id = ?", id).Find(&chat)
	if res.Error != nil {
		return models.Chat{}, res.Error
	}
	return mapChat(chat), nil
}

func mapChats(chats []Chat) []models.Chat {
	res := make([]models.Chat, len(chats))
	for i, chat := range chats {
		res[i] = mapChat(chat)
	}
	return res
}

func mapChat(chat Chat) models.Chat {
	return models.Chat{
		ID:    chat.ID,
		User1: chat.User1,
		User2: chat.User2,
	}
}
