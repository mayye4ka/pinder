package repository

import (
	"errors"

	"gorm.io/gorm"
)

type MessageTranscription struct {
	MessageID     uint64
	Transcription string
}

func (MessageTranscription) TableName() string {
	return "transcriptions"
}

func (r *Repository) GetMessageTranscription(msgID uint64) (string, bool, error) {
	var t MessageTranscription
	res := r.db.Model(&MessageTranscription{}).Where("message_id = ?", msgID).First(&t)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return "", false, nil
	} else if res.Error != nil {
		return "", false, res.Error
	}
	return t.Transcription, true, nil
}

func (r *Repository) SaveMessageTranscription(msgID uint64, text string) error {
	res := r.db.Create(MessageTranscription{
		MessageID:     msgID,
		Transcription: text,
	})
	return res.Error
}
