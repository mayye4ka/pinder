package service

import (
	"github.com/mayye4ka/pinder/internal/models"
)

var (
	chat = models.Chat{
		ID:    1,
		User1: userId,
		User2: user2Id,
	}
	chatShowcase = models.ChatShowcase{
		ID:    1,
		Name:  userName,
		Photo: photo1Link,
	}

	msgText = models.Message{
		ID:          1,
		ChatID:      1,
		SenderID:    userId,
		ContentType: models.ContentText,
		Payload:     "text",
	}
	msgPhoto = models.Message{
		ID:          2,
		ChatID:      1,
		SenderID:    userId,
		ContentType: models.ContentPhoto,
		Payload:     chatPhoto,
	}
	msgVoice = models.Message{
		ID:          3,
		ChatID:      1,
		SenderID:    userId,
		ContentType: models.ContentVoice,
		Payload:     voice,
	}
	msgTextShowcase = models.MessageShowcase{
		ID:          1,
		SentByMe:    true,
		ContentType: models.ContentText,
		Payload:     "text",
	}
	msgPhotoShowcase = models.MessageShowcase{
		ID:          2,
		SentByMe:    true,
		ContentType: models.ContentPhoto,
		Payload:     chatPhotoLink,
	}
	msgVoiceShowcase = models.MessageShowcase{
		ID:          3,
		SentByMe:    true,
		ContentType: models.ContentVoice,
		Payload:     voiceLink,
	}

	voice         = "voice"
	chatPhoto     = "chat_photo"
	voiceLink     = "voice_link"
	chatPhotoLink = "chat_photo_link"
	voiceBytes    = []byte("some voice bytes")
	photoBytes    = []byte("some photo bytes")
)

func (s *ServiceTestSuite) TestListChats() {
	s.repoMock.EXPECT().GetChats(userId).Return([]models.Chat{chat}, nil)
	s.repoMock.EXPECT().GetProfile(user2Id).Return(models.Profile{UserID: user2Id, Name: userName}, nil)
	s.repoMock.EXPECT().GetUserPhotos(user2Id).Return([]string{photo1, photo2}, nil)
	s.fsMock.EXPECT().MakeProfilePhotoLink(user1Ctx, photo1).Return(photo1Link, nil)

	chats, err := s.service.ListChats(user1Ctx)

	s.Nil(err)
	s.Equal([]models.ChatShowcase{chatShowcase}, chats)
}

func (s *ServiceTestSuite) TestListMessages() {
	s.repoMock.EXPECT().GetChat(chat.ID).Return(chat, nil)
	s.repoMock.EXPECT().GetMessages(chat.ID).Return([]models.Message{msgText, msgPhoto, msgVoice}, nil)
	s.fsMock.EXPECT().MakeChatPhotoLink(user1Ctx, msgPhoto.Payload).Return(chatPhotoLink, nil)
	s.fsMock.EXPECT().MakeChatVoiceLink(user1Ctx, msgVoice.Payload).Return(voiceLink, nil)

	messages, err := s.service.ListMessages(user1Ctx, chat.ID)

	s.Nil(err)
	s.Equal(messages, []models.MessageShowcase{
		msgTextShowcase, msgPhotoShowcase, msgVoiceShowcase,
	})
}

func (s *ServiceTestSuite) TestSendMessage_ContentTypeText() {
	s.repoMock.EXPECT().GetChat(chat.ID).Return(chat, nil)
	s.repoMock.EXPECT().SendMessage(chat.ID, userId, models.ContentText, "text").Return(msgText, nil)
	s.userNotifierMock.EXPECT().SendMessage(chat.User1, models.MessageSend{
		ChatID:      chat.ID,
		MessageID:   msgText.ID,
		SentByMe:    true,
		ContentType: models.ContentText,
		Payload:     msgText.Payload,
	}).Return(nil)
	s.userNotifierMock.EXPECT().SendMessage(chat.User2, models.MessageSend{
		ChatID:      chat.ID,
		MessageID:   msgText.ID,
		SentByMe:    false,
		ContentType: models.ContentText,
		Payload:     msgText.Payload,
	}).Return(nil)

	err := s.service.SendMessage(user1Ctx, chat.ID, models.ContentText, "text")

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSendMessage_ContentTypeVoice() {
	s.repoMock.EXPECT().GetChat(chat.ID).Return(chat, nil)
	s.fsMock.EXPECT().SaveChatVoice(user1Ctx, voiceBytes).Return(voice, nil)
	s.repoMock.EXPECT().SendMessage(chat.ID, userId, models.ContentVoice, voice).Return(msgVoice, nil)
	s.fsMock.EXPECT().MakeChatVoiceLink(user1Ctx, voice).Return(voiceLink, nil)

	s.userNotifierMock.EXPECT().SendMessage(chat.User1, models.MessageSend{
		ChatID:      chat.ID,
		MessageID:   msgVoice.ID,
		SentByMe:    true,
		ContentType: models.ContentVoice,
		Payload:     voiceLink,
	}).Return(nil)
	s.userNotifierMock.EXPECT().SendMessage(chat.User2, models.MessageSend{
		ChatID:      chat.ID,
		MessageID:   msgVoice.ID,
		SentByMe:    false,
		ContentType: models.ContentVoice,
		Payload:     voiceLink,
	}).Return(nil)

	err := s.service.SendMessage(user1Ctx, chat.ID, models.ContentVoice, string(voiceBytes))

	s.Nil(err)
}

func (s *ServiceTestSuite) TestSendMessage_ContentTypePhoto() {
	s.repoMock.EXPECT().GetChat(chat.ID).Return(chat, nil)
	s.fsMock.EXPECT().SaveChatPhoto(user1Ctx, photoBytes).Return(chatPhoto, nil)
	s.repoMock.EXPECT().SendMessage(chat.ID, userId, models.ContentPhoto, chatPhoto).Return(msgPhoto, nil)
	s.fsMock.EXPECT().MakeChatPhotoLink(user1Ctx, chatPhoto).Return(chatPhotoLink, nil)

	s.userNotifierMock.EXPECT().SendMessage(chat.User1, models.MessageSend{
		ChatID:      chat.ID,
		MessageID:   msgPhoto.ID,
		SentByMe:    true,
		ContentType: models.ContentPhoto,
		Payload:     chatPhotoLink,
	}).Return(nil)
	s.userNotifierMock.EXPECT().SendMessage(chat.User2, models.MessageSend{
		ChatID:      chat.ID,
		MessageID:   msgPhoto.ID,
		SentByMe:    false,
		ContentType: models.ContentPhoto,
		Payload:     chatPhotoLink,
	}).Return(nil)

	err := s.service.SendMessage(user1Ctx, chat.ID, models.ContentPhoto, string(photoBytes))

	s.Nil(err)
}
