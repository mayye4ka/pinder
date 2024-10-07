package service

import (
	"context"

	"github.com/mayye4ka/pinder/errs"
	"github.com/mayye4ka/pinder/models"
)

const userIdContextKey = "user_id"

var (
	errUnauthenticated = &errs.CodableError{
		Code:    errs.CodePermissionDenied,
		Message: "unauthenticated for this endpoint",
	}
	errPermissionDenied = &errs.CodableError{
		Code:    errs.CodePermissionDenied,
		Message: "no permissions for this action",
	}
)

type Service struct {
	repository   Repository
	filestorage  FileStorage
	userNotifier UserNotifier
	stt          Stt
}

type Repository interface {
	GetProfile(uint64) (models.Profile, error)
	PutProfile(models.Profile) error
	AddPhoto(userID uint64, photoKey string) error
	GetUserPhotos(userID uint64) ([]string, error)
	DeleteUserPhoto(userID uint64, photoKey string) error
	GetPreferences(uint64) (models.Preferences, error)
	PutPreferences(models.Preferences) error
	GetAllValidUsers() ([]uint64, error)

	GetPendingPairAttempts(user1ID uint64) ([]models.PairAttempt, error)
	GetWhoLikedMe(userID uint64) (uint64, error)
	CreateEvent(PAID uint64, eventType models.PEType) error
	GetLatestPairAttempt(user1, user2 uint64) (models.PairAttempt, error)
	GetLatestPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error)
	GetPendingPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error)
	CreatePairAttempt(user1, user2 uint64) (models.PairAttempt, error)
	FinishPairAttempt(PAID uint64, PAState models.PAState) error
	GetLastEvent(PAID uint64) (models.PairEvent, error)

	CreateChat(user1, user2 uint64) error
	GetChats(userID uint64) ([]models.Chat, error)
	GetChat(id uint64) (models.Chat, error)
	SendMessage(chatID, sender uint64, contentType models.MsgContentType, payload string) (models.Message, error)
	GetMessages(chatID uint64) ([]models.Message, error)
	GetMessage(msgID uint64) (models.Message, error)

	GetMessageTranscription(id uint64) (string, bool, error)
	SaveMessageTranscription(id uint64, text string) error
}

type FileStorage interface {
	SaveProfilePhoto(ctx context.Context, photo []byte) (string, error)
	DelProfilePhoto(ctx context.Context, photoKey string) error
	MakeProfilePhotoLink(ctx context.Context, photoKey string) (string, error)

	MakeChatPhotoLink(ctx context.Context, key string) (string, error)
	SaveChatPhoto(ctx context.Context, paylaod []byte) (string, error)

	MakeChatVoiceLink(ctx context.Context, key string) (string, error)
	SaveChatVoice(ctx context.Context, payload []byte) (string, error)
	GetChatVoice(ctx context.Context, key string) (string, error)
}

type UserNotifier interface {
	NotifyMatch(userId uint64, notification models.MatchNotification) error
	NotifyLiked(userId uint64, notification models.LikeNotification) error
	SendMessage(userId uint64, notification models.MessageSend) error
	SendTranscribedMessage(userId uint64, notification models.MessageTranscibed) error
}

type Stt interface {
	PutTask(task models.SttTask) error
	ResultsChan() <-chan models.SttResult
}

func New(repo Repository, filestorage FileStorage, userNotifier UserNotifier, stt Stt) *Service {
	return &Service{
		repository:   repo,
		filestorage:  filestorage,
		userNotifier: userNotifier,
		stt:          stt,
	}
}
