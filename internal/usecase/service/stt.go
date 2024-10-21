package service

import (
	"context"

	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) GetTextFromVoice(ctx context.Context, messageId uint64) (string, bool, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return "", false, errUnauthenticated
	}
	msg, err := s.repository.GetMessage(ctx, messageId)
	if err != nil {
		return "", false, errors.Wrap(err, "can't get message")
	}
	chat, err := s.repository.GetChat(ctx, msg.ChatID)
	if err != nil {
		return "", false, errors.Wrap(err, "can't get chat")
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

	text, found, err := s.repository.GetMessageTranscription(ctx, messageId)
	if err != nil {
		return "", false, errors.Wrap(err, "can't get message transcription")
	}
	if found {
		return text, false, nil
	}

	speech, err := s.filestorage.GetChatVoice(ctx, msg.Payload)
	if err != nil {
		return "", false, errors.Wrap(err, "can't get chat voice")
	}

	err = s.stt.PutTask(ctx, models.SttTask{
		UserID:    userId,
		MessageID: msg.ID,
		Speech:    speech,
	})
	if err != nil {
		return "", false, errors.Wrap(err, "can't put task")
	}
	return "", true, nil
}

func (s *Service) HandleSttResult(ctx context.Context, res models.SttResult) error {
	err := s.repository.SaveMessageTranscription(ctx, res.MessageID, res.Text)
	if err != nil {
		return errors.Wrap(err, "can't save message transcription")
	}
	return s.userNotifier.SendTranscribedMessage(ctx, res.UserID, models.MessageTranscibed{
		MessageID: res.MessageID,
		Text:      res.Text,
	})
}
