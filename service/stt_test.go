package service

import (
	"github.com/mayye4ka/pinder/models"
)

var (
	msgID              = uint64(3)
	voiceTranscription = "voice_transcription"
	speech             = "speech"
)

func (s *ServiceTestSuite) TestGetTextFromVoice_VoiceCached() {
	s.repoMock.EXPECT().GetMessage(uint64(msgID)).Return(msgVoice, nil)
	s.repoMock.EXPECT().GetChat(msgVoice.ChatID).Return(chat, nil)
	s.repoMock.EXPECT().GetMessageTranscription(uint64(msgID)).Return(voiceTranscription, true, nil)

	text, shouldWait, err := s.service.GetTextFromVoice(user1Ctx, msgID)
	s.Nil(err)
	s.False(shouldWait)
	s.Equal(text, voiceTranscription)
}

func (s *ServiceTestSuite) TestGetTextFromVoice_CreateTask() {
	s.repoMock.EXPECT().GetMessage(uint64(msgID)).Return(msgVoice, nil)
	s.repoMock.EXPECT().GetChat(msgVoice.ChatID).Return(chat, nil)
	s.repoMock.EXPECT().GetMessageTranscription(uint64(msgID)).Return("", false, nil)
	s.fsMock.EXPECT().GetChatVoice(user1Ctx, msgVoice.Payload).Return(speech, nil)
	s.sttMock.EXPECT().PutTask(models.SttTask{
		UserID:    userId,
		MessageID: msgID,
		Speech:    speech,
	}).Return(nil)
	text, shouldWait, err := s.service.GetTextFromVoice(user1Ctx, msgID)
	s.Nil(err)
	s.True(shouldWait)
	s.Equal("", text)
}

func (s *ServiceTestSuite) TestHandleSttResults() {
	s.repoMock.EXPECT().SaveMessageTranscription(uint64(msgID), voiceTranscription).Return(nil)
	s.userNotifierMock.EXPECT().SendTranscribedMessage(userId, models.MessageTranscibed{
		MessageID: msgID,
		Text:      voiceTranscription,
	}).Return(nil)

	err := s.service.HandleSttResult(models.SttResult{
		UserID:    userId,
		MessageID: uint64(msgID),
		Text:      voiceTranscription,
	})

	s.Nil(err)
}
