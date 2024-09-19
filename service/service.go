package service

import (
	"pinder/models"
)

type Service struct {
	repository Repository
}

type Repository interface {
	CreateUser(phoneNumber, passHash string) (*models.User, error)
	GetUserByCreds(phoneNumber, passHash string) (*models.User, error)
	GetUser(id uint64) (*models.User, error)
	GetProfile(uint64) (*models.Profile, error)
	PutProfile(models.Profile) error
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

func New(repo Repository) *Service {
	return &Service{
		repository: repo,
	}
}
