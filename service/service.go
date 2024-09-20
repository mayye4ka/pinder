package service

import (
	"context"
	"pinder/models"
)

type Service struct {
	repository  Repository
	filestorage FileStorage
}

type Repository interface {
	CreateUser(phoneNumber, passHash string) (*models.User, error)
	GetUserByCreds(phoneNumber, passHash string) (*models.User, error)
	GetUser(id uint64) (*models.User, error)
	GetProfile(uint64) (*models.Profile, error)
	PutProfileData(models.Profile) error
	PutProfilePhoto(userId uint64, photoKey string) error
	GetPreferences(uint64) (*models.Preferences, error)
	PutPreferences(models.Preferences) error

	GetHangingPartner(userID uint64) (*models.Profile, error)
	GetWhoLikedMe(userID uint64) (*models.Profile, error)
	ChooseCandidateAndCreatePairAttempt(userID uint64) (*models.Profile, error)
	CreateEvent(PAID uint64, eventType models.PEType) error
	GetLatestPairAttempt(user1, user2 uint64) (models.PairAttempt, error)
	GetPendingPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error)
	FinishPairAttempt(PAID uint64, PAState models.PAState) error
}

type FileStorage interface {
	SavePhoto(ctx context.Context, photo []byte) (string, error)
	DelPhoto(ctx context.Context, photoKey string) error
	MakeLink(ctx context.Context, photoKey string) (string, error)
}

func New(repo Repository, filestorage FileStorage) *Service {
	return &Service{
		repository:  repo,
		filestorage: filestorage,
	}
}
