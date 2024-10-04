package repository

import (
	"time"

	"github.com/mayye4ka/pinder/models"
)

type PairAttempt struct {
	ID        uint64
	User1     uint64
	User2     uint64
	State     PAState
	CreatedAt time.Time
}

func (PairAttempt) TableName() string {
	return "pair_attempts"
}

type PAState string

const (
	PAStatePending  PAState = "pending"
	PAStateMatch    PAState = "match"
	PAStateMismatch PAState = "mismatch"
)

func (r *Repository) GetLatestPairAttempt(user1, user2 uint64) (models.PairAttempt, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("user1 = ? and user2 = ?", user1, user2).
		Order("created_at desc").First(&pair)
	if res.Error != nil {
		return models.PairAttempt{}, res.Error
	}
	return mapPairAttempt(pair), nil
}

func (r *Repository) GetPendingPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("((user1 = ? and user2 = ?) or (user2 = ? and user1 = ?)) and state = ?",
			user1, user2, user1, user2, models.PAStatePending).
		First(&pair)
	if res.Error != nil {
		return models.PairAttempt{}, res.Error
	}
	return mapPairAttempt(pair), nil
}

func (r *Repository) GetLatestPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("((user1 = ? and user2 = ?) or (user2 = ? and user1 = ?))",
			user1, user2, user1, user2).
		Order("created_at desc").
		First(&pair)
	if res.Error != nil {
		return models.PairAttempt{}, res.Error
	}
	return mapPairAttempt(pair), nil
}

func (r *Repository) FinishPairAttempt(PAID uint64, state models.PAState) error {
	res := r.db.Model(&PairAttempt{}).Where("id = ?", PAID).Update("state = ?", unmapPaState(state))
	return res.Error
}

func (r *Repository) CreatePairAttempt(user1, user2 uint64) error {
	res := r.db.Create(&PairAttempt{
		User1:     user1,
		User2:     user2,
		State:     PAStatePending,
		CreatedAt: time.Now(),
	})
	return res.Error
}

func mapPairAttempt(pa PairAttempt) models.PairAttempt {
	return models.PairAttempt{
		ID:        pa.ID,
		User1:     pa.User1,
		User2:     pa.User2,
		State:     mapPaState(pa.State),
		CreatedAt: pa.CreatedAt,
	}
}

func mapPaState(s PAState) models.PAState {
	switch s {
	case PAStateMatch:
		return models.PAStateMatch
	case PAStateMismatch:
		return models.PAStateMismatch
	default:
		return models.PAStatePending
	}
}

func unmapPaState(s models.PAState) PAState {
	switch s {
	case models.PAStateMatch:
		return PAStateMatch
	case models.PAStateMismatch:
		return PAStateMismatch
	default:
		return PAStatePending
	}
}
