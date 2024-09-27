package service

import (
	"errors"
	"pinder/models"
	"pinder/server"

	"golang.org/x/net/context"
)

func getWhoIsNotMe(chat models.Chat, userId uint64) uint64 {
	if chat.User1 == userId {
		return chat.User2
	}
	return chat.User1
}

func (s *Service) ListChats(ctx context.Context, req *server.RequestWithToken) (*server.ListChatsResponse, error) {
	userID, err := verifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	chats, err := s.repository.GetChats(userID)
	if err != nil {
		return nil, err
	}
	mappedChats := []server.Chat{}
	for _, chat := range chats {
		user2 := getWhoIsNotMe(chat, userID)
		prof, err := s.repository.GetProfile(user2)
		if err != nil {
			return nil, err
		}
		photos, err := s.repository.GetUserPhotos(user2)
		if err != nil {
			return nil, err
		}
		link, err := s.filestorage.MakeLink(ctx, photos[0])
		if err != nil {
			return nil, err
		}
		mappedChats = append(mappedChats, server.Chat{
			ChatID: chat.ID,
			Name:   prof.Name,
			Photo:  link,
		})
	}

	return &server.ListChatsResponse{
		Chats: mappedChats,
	}, nil
}

func (s *Service) ListMessages(ctx context.Context, req *server.ListMessagesRequest) (*server.ListMessagesResponse, error) {
	userID, err := verifyToken(req.Token)
	if err != nil {
		return nil, err
	}
	chat, err := s.repository.GetChat(req.ChatId)
	if err != nil {
		return nil, err
	}
	if chat.User1 != userID && chat.User2 != userID {
		return nil, errors.New("access denied")
	}
	messages, err := s.repository.GetMessages(req.ChatId)
	if err != nil {
		return nil, err
	}
	mappedMessages := []server.Message{}
	for _, msg := range messages {
		sentByMe := true
		if msg.SenderID != userID {
			sentByMe = false
		}
		mappedMessages = append(mappedMessages, server.Message{
			ID:          msg.ID,
			SentByMe:    sentByMe,
			ContentType: unmapContentType(msg.ContentType),
			Payload:     msg.Payload,
			CreatedAt:   msg.CreatedAt,
		})
	}
	return &server.ListMessagesResponse{
		Messages: mappedMessages,
	}, nil
}
