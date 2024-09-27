package repository

import (
	"pinder/models"
	"time"
)

func (r *Repository) CreateChat(user1, user2 uint64) error {
	chat := models.Chat{
		User1: user1,
		User2: user2,
	}
	res := r.db.Create(&chat)
	return res.Error
}

func (r *Repository) GetChats(userID uint64) ([]models.Chat, error) {
	var chats []models.Chat
	res := r.db.Model(&models.Chat{}).Where("user_1 = ? or user_2 = ?", userID, userID).Find(&chats)
	if res.Error != nil {
		return nil, res.Error
	}
	return chats, nil
}

func (r *Repository) GetChat(id uint64) (models.Chat, error) {
	var chat models.Chat
	res := r.db.Model(&models.Chat{}).Where("id = ?", id).Find(&chat)
	if res.Error != nil {
		return models.Chat{}, res.Error
	}
	return chat, nil
}

func (r *Repository) SendMessage(chatID, sender uint64, contentType models.MsgContentType, payload string) error {
	message := models.Message{
		ChatID:      chatID,
		SenderID:    sender,
		ContentType: contentType,
		Payload:     payload,
		CreatedAt:   time.Now(),
	}
	res := r.db.Create(&message)
	return res.Error
}

func (r *Repository) GetMessages(chatID uint64) ([]models.Message, error) {
	var messages []models.Message
	res := r.db.Model(&models.Message{}).Where("chat_id = ?", chatID).Order("created_at").Find(&messages)
	if res.Error != nil {
		return nil, res.Error
	}
	return messages, nil
}
