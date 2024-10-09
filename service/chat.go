package service

import (
	"github.com/mayye4ka/pinder/models"
	"github.com/pkg/errors"

	"golang.org/x/net/context"
)

func getWhoIsNotMe(id1, id2, userId uint64) uint64 {
	if id1 == userId {
		return id2
	}
	return id1
}

func (s *Service) enrichMessageWithLinks(ctx context.Context, message *models.Message) error {
	if message.ContentType == models.ContentPhoto {
		link, err := s.filestorage.MakeChatPhotoLink(ctx, message.Payload)
		if err != nil {
			return errors.Wrap(err, "can't make chat photo link")
		}
		message.Payload = link
	}
	if message.ContentType == models.ContentVoice {
		link, err := s.filestorage.MakeChatVoiceLink(ctx, message.Payload)
		if err != nil {
			return errors.Wrap(err, "can't make chat voice link")
		}
		message.Payload = link
	}
	return nil
}

func (s *Service) ListChats(ctx context.Context) ([]models.ChatShowcase, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return nil, errUnauthenticated
	}
	chats, err := s.repository.GetChats(userId)
	if err != nil {
		return nil, errors.Wrap(err, "can't get chats by user id")
	}
	res := []models.ChatShowcase{}
	for _, chat := range chats {
		user2 := getWhoIsNotMe(chat.User1, chat.User2, userId)
		prof, err := s.repository.GetProfile(user2)
		if err != nil {
			return nil, errors.Wrap(err, "can't get profile")
		}
		photos, err := s.repository.GetUserPhotos(user2)
		if err != nil {
			return nil, errors.Wrap(err, "can't get user photos")
		}
		link, err := s.filestorage.MakeProfilePhotoLink(ctx, photos[0])
		if err != nil {
			return nil, errors.Wrap(err, "can't make profile photo link")
		}
		res = append(res, models.ChatShowcase{
			ID:    chat.ID,
			Name:  prof.Name,
			Photo: link,
		})
	}

	return res, nil
}

func (s *Service) ListMessages(ctx context.Context, chatId uint64) ([]models.MessageShowcase, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return nil, errUnauthenticated
	}
	chat, err := s.repository.GetChat(chatId)
	if err != nil {
		return nil, errors.Wrap(err, "can't get chat by id")
	}
	if chat.User1 != userId && chat.User2 != userId {
		return nil, errPermissionDenied
	}
	messages, err := s.repository.GetMessages(chatId)
	if err != nil {
		return nil, errors.Wrap(err, "can't get messages")
	}
	res := []models.MessageShowcase{}
	for _, msg := range messages {
		sentByMe := true
		if msg.SenderID != userId {
			sentByMe = false
		}
		err = s.enrichMessageWithLinks(ctx, &msg)
		if err != nil {
			return nil, errors.Wrap(err, "can't enrich message with links")
		}
		res = append(res, models.MessageShowcase{
			ID:          msg.ID,
			SentByMe:    sentByMe,
			ContentType: msg.ContentType,
			Payload:     msg.Payload,
			CreatedAt:   msg.CreatedAt,
		})
	}
	return res, nil
}

func (s *Service) SendMessage(ctx context.Context, chatId uint64, contentType models.MsgContentType, payload string) error {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return errUnauthenticated
	}
	chat, err := s.repository.GetChat(chatId)
	if err != nil {
		return errors.Wrap(err, "can't get chat")
	}
	if chat.User1 != userId && chat.User2 != userId {
		return errPermissionDenied
	}
	if contentType == models.ContentVoice {
		key, err := s.filestorage.SaveChatVoice(ctx, []byte(payload))
		if err != nil {
			return errors.Wrap(err, "can't save chat voice")
		}
		payload = key
	} else if contentType == models.ContentPhoto {
		key, err := s.filestorage.SaveChatPhoto(ctx, []byte(payload))
		if err != nil {
			return errors.Wrap(err, "can't save chat photo")
		}
		payload = key
	}
	msg, err := s.repository.SendMessage(chatId, userId, contentType, payload)
	if err != nil {
		return errors.Wrap(err, "can't send message")
	}
	err = s.enrichMessageWithLinks(ctx, &msg)
	if err != nil {
		return errors.Wrap(err, "can't enrich message with links")
	}
	for _, recv := range []uint64{chat.User1, chat.User2} {
		sentByMe := true
		if recv != userId {
			sentByMe = false
		}
		err = s.userNotifier.SendMessage(recv, models.MessageSend{
			ChatID:      msg.ChatID,
			MessageID:   msg.ID,
			SentByMe:    sentByMe,
			ContentType: msg.ContentType,
			Payload:     msg.Payload,
		})
		if err != nil {
			return errors.Wrap(err, "can't send message")
		}
	}
	return nil
}
