package repository

import (
	"errors"

	"github.com/mayye4ka/pinder/internal/errs"
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
		r.logger.Err(res.Error).Msg("can't get message transcription")
		return "", false, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get message transcription",
		}
	}
	return t.Transcription, true, nil
}

func (r *Repository) SaveMessageTranscription(msgID uint64, text string) error {
	res := r.db.Create(MessageTranscription{
		MessageID:     msgID,
		Transcription: text,
	})
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't save message transcription")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't save message transcription",
		}
	}
	return nil
}
