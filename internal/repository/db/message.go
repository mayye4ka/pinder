package repository

import (
	"context"
	"errors"
	"time"

	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
	"gorm.io/gorm"
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

func (r *Repository) SendMessage(ctx context.Context, chatID, sender uint64, contentType models.MsgContentType, payload string) (models.Message, error) {
	message := Message{
		ChatID:      chatID,
		SenderID:    sender,
		ContentType: unmapContentType(contentType),
		Payload:     payload,
		CreatedAt:   time.Now(),
	}
	res := r.db.WithContext(ctx).Create(&message)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't send message")
		return models.Message{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't send message",
		}
	}
	return mapMessage(message), nil
}

func (r *Repository) GetMessages(ctx context.Context, chatID uint64) ([]models.Message, error) {
	var messages []Message
	res := r.db.WithContext(ctx).Model(&Message{}).Where("chat_id = ?", chatID).Order("created_at").Find(&messages)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get messages in this chat")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get messages in this chat",
		}
	}
	return mapMessages(messages), nil
}

func (r *Repository) GetMessage(ctx context.Context, msgID uint64) (models.Message, error) {
	var message Message
	res := r.db.WithContext(ctx).Model(&Message{}).Where("id = ?", msgID).First(&message)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return models.Message{}, &errs.CodableError{
				Code:    errs.CodeNotFound,
				Message: "not found this message",
			}
		}
		r.logger.Err(res.Error).Msg("can't get message")
		return models.Message{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get message",
		}
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
