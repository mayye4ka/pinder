package service

import (
	"context"

	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
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
	GetProfile(ctx context.Context, userID uint64) (models.Profile, error)
	PutProfile(ctx context.Context, newProfile models.Profile) error
	AddPhoto(ctx context.Context, userID uint64, photoKey string) error
	GetUserPhotos(ctx context.Context, userID uint64) ([]string, error)
	DeleteUserPhoto(ctx context.Context, userID uint64, photoKey string) error
	ReorderPhotos(ctx context.Context, newOrder []string) error
	GetPreferences(ctx context.Context, userID uint64) (models.Preferences, error)
	PutPreferences(ctx context.Context, newPreferences models.Preferences) error
	GetAllValidUsers(ctx context.Context) ([]uint64, error)

	GetPendingPairAttempts(ctx context.Context, user1ID uint64) ([]models.PairAttempt, error)
	GetWhoLikedMe(ctx context.Context, userID uint64) (uint64, error)
	CreateEvent(ctx context.Context, PAID uint64, eventType models.PEType) error
	GetLatestPairAttempt(ctx context.Context, user1, user2 uint64) (models.PairAttempt, error)
	GetLatestPairAttemptByUserPair(ctx context.Context, user1, user2 uint64) (models.PairAttempt, error)
	GetPendingPairAttemptByUserPair(ctx context.Context, user1, user2 uint64) (models.PairAttempt, error)
	CreatePairAttempt(ctx context.Context, user1, user2 uint64) (models.PairAttempt, error)
	FinishPairAttempt(ctx context.Context, PAID uint64, PAState models.PAState) error
	GetLastEvent(ctx context.Context, PAID uint64) (models.PairEvent, error)

	CreateChat(ctx context.Context, user1, user2 uint64) error
	GetChats(ctx context.Context, userID uint64) ([]models.Chat, error)
	GetChat(ctx context.Context, id uint64) (models.Chat, error)
	SendMessage(ctx context.Context, chatID, sender uint64, contentType models.MsgContentType, payload string) (models.Message, error)
	GetMessages(ctx context.Context, chatID uint64) ([]models.Message, error)
	GetMessage(ctx context.Context, msgID uint64) (models.Message, error)

	GetMessageTranscription(ctx context.Context, id uint64) (string, bool, error)
	SaveMessageTranscription(ctx context.Context, id uint64, text string) error
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
	NotifyMatch(ctx context.Context, userId uint64, notification models.MatchNotification) error
	NotifyLiked(ctx context.Context, userId uint64, notification models.LikeNotification) error
	SendMessage(ctx context.Context, userId uint64, notification models.MessageSend) error
	SendTranscribedMessage(ctx context.Context, userId uint64, notification models.MessageTranscibed) error
}

type Stt interface {
	PutTask(ctx context.Context, task models.SttTask) error
}

func New(repo Repository, filestorage FileStorage, userNotifier UserNotifier, stt Stt) *Service {
	return &Service{
		repository:   repo,
		filestorage:  filestorage,
		userNotifier: userNotifier,
		stt:          stt,
	}
}
