package service

import (
	"context"
	"pinder/models"
)

type Service struct {
	repository     Repository
	filestorage    FileStorage
	userInteractor UserInteractor
}

type Repository interface {
	CreateUser(phoneNumber, passHash string) (*models.User, error)
	GetUserByCreds(phoneNumber, passHash string) (*models.User, error)
	GetUser(id uint64) (*models.User, error)
	GetProfile(uint64) (*models.Profile, error)
	PutProfile(models.Profile) error
	AddPhoto(userID uint64, photoKey string) error
	GetUserPhotos(userID uint64) ([]string, error)
	DeleteUserPhoto(userID uint64, photoKey string) error
	GetPreferences(uint64) (*models.Preferences, error)
	PutPreferences(models.Preferences) error

	GetHangingPartner(userID uint64) (*models.Profile, error)
	GetWhoLikedMe(userID uint64) (*models.Profile, error)
	ChooseCandidateAndCreatePairAttempt(userID uint64) (*models.Profile, error)
	CreateEvent(PAID uint64, eventType models.PEType) error
	GetLatestPairAttempt(user1, user2 uint64) (models.PairAttempt, error)
	GetPendingPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error)
	FinishPairAttempt(PAID uint64, PAState models.PAState) error

	CreateChat(user1, user2 uint64) error
	GetChats(userID uint64) ([]models.Chat, error)
	GetChat(id uint64) (models.Chat, error)
	SendMessage(chatID, sender uint64, contentType models.MsgContentType, payload string) (models.Message, error)
	GetMessages(chatID uint64) ([]models.Message, error)
}

type FileStorage interface {
	SaveProfilePhoto(ctx context.Context, photo []byte) (string, error)
	DelProfilePhoto(ctx context.Context, photoKey string) error
	MakeProfilePhotoLink(ctx context.Context, photoKey string) (string, error)

	MakeChatPhotoLink(ctx context.Context, key string) (string, error)
	SaveChatPhoto(ctx context.Context, paylaod []byte) (string, error)

	MakeChatVoiceLink(ctx context.Context, key string) (string, error)
	SaveChatVoice(ctx context.Context, payload []byte) (string, error)
}

type UserInteractor interface {
	SendMessage(chat models.Chat, message models.Message)
}

func New(repo Repository, filestorage FileStorage, userInteractor UserInteractor) *Service {
	return &Service{
		repository:     repo,
		filestorage:    filestorage,
		userInteractor: userInteractor,
	}
}
