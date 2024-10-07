package service

import (
	"context"

	"github.com/mayye4ka/pinder/errs"
	"github.com/mayye4ka/pinder/models"
)

func (s *Service) GetTextFromVoice(ctx context.Context, messageId uint64) (string, bool, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return "", false, errUnauthenticated
	}
	msg, err := s.repository.GetMessage(messageId)
	if err != nil {
		return "", false, err
	}
	chat, err := s.repository.GetChat(msg.ChatID)
	if err != nil {
		return "", false, err
	}
	if chat.User1 != userId && chat.User2 != userId {
		return "", false, errPermissionDenied
	}
	if msg.ContentType != models.ContentVoice {
		return "", false, &errs.CodableError{
			Code:    errs.CodeInvalidInput,
			Message: "bad message content to transcribe",
		}
	}

	text, found, err := s.repository.GetMessageTranscription(messageId)
	if err != nil {
		return "", false, err
	}
	if found {
		return text, false, nil
	}

	speech, err := s.filestorage.GetChatVoice(ctx, msg.Payload)
	if err != nil {
		return "", false, err
	}

	err = s.stt.PutTask(models.SttTask{
		UserID:    userId,
		MessageID: msg.ID,
		Speech:    speech,
	})
	if err != nil {
		return "", false, err
	}
	return "", true, nil
}

func (s *Service) Start() error {
	c := s.stt.ResultsChan()
	for res := range c {
		err := s.handleSttResult(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) handleSttResult(res models.SttResult) error {
	err := s.repository.SaveMessageTranscription(res.MessageID, res.Text)
	if err != nil {
		return err
	}
	return s.userNotifier.SendTranscribedMessage(res.UserID, models.MessageTranscibed{
		MessageID: res.MessageID,
		Text:      res.Text,
	})
}
