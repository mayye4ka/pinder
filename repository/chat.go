package repository

import (
	"errors"

	"github.com/mayye4ka/pinder/errs"
	"github.com/mayye4ka/pinder/models"
	"gorm.io/gorm"
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
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't create chat")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't create chat",
		}
	}
	return nil
}

func (r *Repository) GetChats(userID uint64) ([]models.Chat, error) {
	var chats []Chat
	res := r.db.Model(&Chat{}).Where("user_1 = ? or user_2 = ?", userID, userID).Find(&chats)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get chats")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get chats",
		}
	}
	return mapChats(chats), nil
}

func (r *Repository) GetChat(id uint64) (models.Chat, error) {
	var chat Chat
	res := r.db.Model(&Chat{}).Where("id = ?", id).Find(&chat)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return models.Chat{}, &errs.CodableError{
				Code:    errs.CodeNotFound,
				Message: "no such chat",
			}
		}
		r.logger.Err(res.Error).Msg("can't find chat")
		return models.Chat{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't find chat",
		}
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
