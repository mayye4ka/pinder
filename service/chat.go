package service

import (
	"errors"

	"github.com/mayye4ka/pinder/models"

	"golang.org/x/net/context"
)

func getWhoIsNotMe(chat models.Chat, userId uint64) uint64 {
	if chat.User1 == userId {
		return chat.User2
	}
	return chat.User1
}

func (s *Service) enrichMessageWithLinks(ctx context.Context, message *models.Message) error {
	if message.ContentType == models.ContentPhoto {
		link, err := s.filestorage.MakeChatPhotoLink(ctx, message.Payload)
		if err != nil {
			return err
		}
		message.Payload = link
	}
	if message.ContentType == models.ContentVoice {
		link, err := s.filestorage.MakeChatVoiceLink(ctx, message.Payload)
		if err != nil {
			return err
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
		return nil, err
	}
	res := []models.ChatShowcase{}
	for _, chat := range chats {
		user2 := getWhoIsNotMe(chat, userId)
		prof, err := s.repository.GetProfile(user2)
		if err != nil {
			return nil, err
		}
		photos, err := s.repository.GetUserPhotos(user2)
		if err != nil {
			return nil, err
		}
		link, err := s.filestorage.MakeProfilePhotoLink(ctx, photos[0])
		if err != nil {
			return nil, err
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
		return nil, err
	}
	if chat.User1 != userId && chat.User2 != userId {
		return nil, errors.New("access denied")
	}
	messages, err := s.repository.GetMessages(chatId)
	if err != nil {
		return nil, err
	}
	res := []models.MessageShowcase{}
	for _, msg := range messages {
		sentByMe := true
		if msg.SenderID != userId {
			sentByMe = false
		}
		err = s.enrichMessageWithLinks(ctx, &msg)
		if err != nil {
			return nil, err
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
		return err
	}
	if chat.User1 != userId && chat.User2 != userId {
		return errors.New("access denied")
	}
	if contentType == models.ContentVoice {
		key, err := s.filestorage.SaveChatVoice(ctx, []byte(payload))
		if err != nil {
			return nil
		}
		payload = key
	} else if contentType == models.ContentPhoto {
		key, err := s.filestorage.SaveChatPhoto(ctx, []byte(payload))
		if err != nil {
			return nil
		}
		payload = key
	}
	msg, err := s.repository.SendMessage(chatId, userId, contentType, payload)
	if err != nil {
		return err
	}
	err = s.enrichMessageWithLinks(ctx, &msg)
	if err != nil {
		return err
	}
	for _, recv := range []uint64{chat.User1, chat.User2} {
		sentByMe := true
		if recv != userId {
			sentByMe = false
		}
		err = s.userNotifier.SendMessage(recv, models.MessageSend{
			ChatID:      chatId,
			MessageID:   msg.ID,
			SentByMe:    sentByMe,
			ContentType: contentType,
			Payload:     payload,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
