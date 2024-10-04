package repository

import (
	"time"

	"github.com/mayye4ka/pinder/models"
)

type MsgContentType string

const (
	ContentText  MsgContentType = "text"
	ContentPhoto MsgContentType = "photo"
	ContentVoice MsgContentType = "voice"
)

type Message struct {
	ID          uint64
	ChatID      uint64
	SenderID    uint64
	ContentType MsgContentType
	Payload     string
	CreatedAt   time.Time
}

func (Message) TableName() string {
	return "messages"
}

func (r *Repository) SendMessage(chatID, sender uint64, contentType models.MsgContentType, payload string) (models.Message, error) {
	message := Message{
		ChatID:      chatID,
		SenderID:    sender,
		ContentType: unmapContentType(contentType),
		Payload:     payload,
		CreatedAt:   time.Now(),
	}
	res := r.db.Create(&message)
	if res.Error != nil {
		return models.Message{}, res.Error
	}
	return mapMessage(message), nil
}

func (r *Repository) GetMessages(chatID uint64) ([]models.Message, error) {
	var messages []Message
	res := r.db.Model(&Message{}).Where("chat_id = ?", chatID).Order("created_at").Find(&messages)
	if res.Error != nil {
		return nil, res.Error
	}
	return mapMessages(messages), nil
}

func (r *Repository) GetMessage(msgID uint64) (models.Message, error) {
	var message Message
	res := r.db.Model(&Message{}).Where("id = ?", msgID).First(&message)
	if res.Error != nil {
		return models.Message{}, res.Error
	}
	return mapMessage(message), nil
}

func mapMessages(msgs []Message) []models.Message {
	res := make([]models.Message, len(msgs))
	for i, msg := range msgs {
		res[i] = mapMessage(msg)
	}
	return res
}

func mapMessage(msg Message) models.Message {
	return models.Message{
		ID:          msg.ID,
		ChatID:      msg.ChatID,
		SenderID:    msg.SenderID,
		ContentType: mapContentType(msg.ContentType),
		Payload:     msg.Payload,
		CreatedAt:   msg.CreatedAt,
	}
}

func mapContentType(ct MsgContentType) models.MsgContentType {
	switch ct {
	case ContentPhoto:
		return models.ContentPhoto
	case ContentVoice:
		return models.ContentVoice
	default:
		return models.ContentText
	}
}

func unmapContentType(ct models.MsgContentType) MsgContentType {
	switch ct {
	case models.ContentPhoto:
		return ContentPhoto
	case models.ContentVoice:
		return ContentVoice
	default:
		return ContentText
	}
}
